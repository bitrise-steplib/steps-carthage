package cachedcarthage

import (
	"errors"
	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The first part writes the given string to stderr and the second part provides the exit code 1.
const (
	failingCommandWithTimeoutStderr           = "echo timed out 1>&2 && false"
	failingCommandWithFailedToConnectToStderr = "echo failed to connect to 1>&2 && false"
)

// Run
func Test_GivenNotBootstrapCommand_WhenRunCalled_ThenExpectNoErrorAndCacheNotCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache()
	mockCommandBuilder := givenStubbedCommandBuilder()
	runner := Runner{
		carthageCommand: "version",
		cache:           mockCarthageCache,
		commandBuilder:  mockCommandBuilder,
	}

	// When
	err := runner.Run()

	// Then
	assert.NoError(t, err)
	mockCarthageCache.AssertNotCalled(t, "IsAvailable")
	mockCarthageCache.AssertNotCalled(t, "Commit")
	mockCarthageCache.AssertNotCalled(t, "CreateIndicator")
}

func Test_GivenBootstrapCommandAndCacheNotAvailableAndCacheCreateFails_WhenRunCalled_ThenExpectError(t *testing.T) {
	// Given
	expectedError := errors.New("sad error")
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(false).
		GivenCreateIndicatorFails(expectedError)
	runner := Runner{
		carthageCommand: "bootstrap",
		cache:           mockCarthageCache,
		commandBuilder:  givenStubbedCommandBuilder(),
	}

	// When
	err := runner.Run()

	// Then
	assert.EqualError(t, expectedError, err.Error())
}

func Test_GivenBootstrapCommandAndCacheNotAvailableAndCacheCreateSucceeds_WhenRunCalled_ThenExpectNoError(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(false).
		GivenCreateIndicatorSucceeds().
		GivenCommitSucceeds()
	runner := Runner{
		carthageCommand: "bootstrap",
		cache:           mockCarthageCache,
		commandBuilder:  givenStubbedCommandBuilder(),
	}

	// When
	err := runner.Run()

	// Then
	assert.NoError(t, err)
	mockCarthageCache.AssertCalled(t, "CreateIndicator")
	mockCarthageCache.AssertCalled(t, "Commit")
}

func Test_GivenBootstrapCommandAndCacheAvailableAndCollectFails_WhenRunCalled_ThenExpectCommandExecutedCacheCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(true).
		GivenCommitFails(errors.New("sad error")).
		GivenCreateIndicatorSucceeds()

	mockCommandBuilder := givenStubbedCommandBuilder()

	runner := Runner{
		carthageCommand: "bootstrap",
		cache:           mockCarthageCache,
		commandBuilder:  mockCommandBuilder,
	}

	// When
	err := runner.Run()

	// Then
	assert.NoError(t, err)
	mockCarthageCache.AssertCalled(t, "CreateIndicator")
	mockCarthageCache.AssertCalled(t, "Commit")
}

func Test_GivenBootstrapCommandAndCacheAvailableAndCollectSucceeds_WhenRunCalled_ThenExpectCommandNotExecutedAndCacheCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(true).
		GivenCommitSucceeds()

	mockCommandBuilder := givenStubbedCommandBuilderReturnFailingCommand()
	runner := Runner{
		carthageCommand: "bootstrap",
		cache:           mockCarthageCache,
		commandBuilder:  mockCommandBuilder,
	}

	// When
	err := runner.Run()

	// Then
	assert.NoError(t, err)
	mockCarthageCache.AssertNotCalled(t, "CreateIndicator")
	mockCarthageCache.AssertNumberOfCalls(t, "Commit", 1)
}

// Retry on failure
func Test_GivenBootstrapCommandAndSingleNetworkFailure_WhenRunCalled_ThenExpectCommandToBeRetriedAndSucceed(t *testing.T) {
	// Given
	commands := []*command.Model{
		command.New("bash", "-c", failingCommandWithTimeoutStderr),
		command.New("echo", "hello"),
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("bootstrap", commands)

	// When
	err := runner.Run()

	// Then
	assert.NoError(t, err)
}

func Test_GivenBootstrapCommandAndPermanentNetworkFailure_WhenRunCalled_ThenExpectCommandToFail(t *testing.T) {
	// Given
	blueprints := []*command.Model{
		command.New("bash", "-c", failingCommandWithTimeoutStderr),
		command.New("bash", "-c", failingCommandWithFailedToConnectToStderr),
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("bootstrap", blueprints)

	// When
	err := runner.Run()

	// Then
	assert.Error(t, err)
}

func Test_GivenUpdateCommandAndSingleNetworkFailure_WhenRunCalled_ThenExpectCommandToBeRetriedAndSucceed(t *testing.T) {
	// Given
	blueprints := []*command.Model{
		command.New("bash", "-c", failingCommandWithFailedToConnectToStderr),
		command.New("echo", "hello"),
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("update", blueprints)

	// When
	err := runner.Run()

	// Then
	assert.NoError(t, err)
}

func Test_GivenUpdateCommandAndPermanentNetworkFailure_WhenRunCalled_ThenExpectCommandToFail(t *testing.T) {
	// Given
	blueprints := []*command.Model{
		command.New("bash", "-c", failingCommandWithFailedToConnectToStderr),
		command.New("bash", "-c", failingCommandWithTimeoutStderr),
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("update", blueprints)

	// When
	err := runner.Run()

	// Then
	assert.Error(t, err)
}

// isCacheAvailable
func Test_GivenCarthageCacheAvailableFails_WhenIsCacheAvailableCalled_ThenExpectFalse(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableFails(errors.New("whatever"))

	runner := Runner{
		cache: mockCarthageCache,
	}

	// When
	available := runner.isCacheAvailable()

	// Then
	assert.False(t, available)
}

func Test_GivenCarthageCacheAvailableSucceeds_WhenIsCacheAvailableCalled_ThenExpectResult(t *testing.T) {
	// Given
	const expectedAvailable = true
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(expectedAvailable)

	runner := Runner{
		cache: mockCarthageCache,
	}

	// When
	available := runner.isCacheAvailable()

	// Then
	assert.Equal(t, expectedAvailable, available)
}

// ecxecuteCommand
func Test_GivenCommadSucceeds_WhenExecuteCommandCalled_ThenExpectNoError(t *testing.T) {
	// Given
	mockCommandBuilder := givenStubbedCommandBuilder()
	cmd := "version"
	args := []string{"arg1"}
	runner := Runner{
		carthageCommand: cmd,
		args:            args,
		commandBuilder:  mockCommandBuilder,
	}

	// When
	err := runner.executeCommand()

	// Then
	assert.NoError(t, err)
	mockCommandBuilder.AssertCalled(t, "AddGitHubToken", mock.Anything)
	mockCommandBuilder.AssertCalled(t, "AddXCConfigFile", mock.Anything)
	mockCommandBuilder.AssertCalled(t, "Append", []string{cmd})
	mockCommandBuilder.AssertCalled(t, "Append", args)
}

// helpers
func givenMockCarthageCache() *MockCarthageCache {
	return new(MockCarthageCache)
}

func givenStubbedCommandBuilder() *MockCommandBuilder {
	cmd := command.New("echo", "hello")
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(cmd)
	return mockCommandBuilder
}

func givenStubbedCommandBuilderReturnFailingCommand() *MockCommandBuilder {
	cmd := command.New("fail")
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(cmd)
	return mockCommandBuilder
}

func givenRunnerWithMainAndCommandBuilderCommands(mainCommand string, commands []*command.Model) Runner {
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(false).
		GivenCreateIndicatorSucceeds().
		GivenCommitSucceeds()

	return Runner{
		carthageCommand: mainCommand,
		cache:           mockCarthageCache,
		commandBuilder:  givenStubbedCommandBuilderReturnsCommands(commands),
	}
}

func givenStubbedCommandBuilderReturnsCommands(commandBlueprints []*command.Model) *MockCommandBuilder {
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandsReturned(commandBlueprints)
	return mockCommandBuilder
}
