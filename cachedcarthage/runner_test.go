package cachedcarthage

import (
	"errors"
	"github.com/stretchr/testify/mock"
	"testing"

	"github.com/stretchr/testify/assert"
)

// The first part writes the given string to stderr and the second part provides the exit code 1.
const(
	failingCommandWithTimeoutStderr = "echo timed out 1>&2 && false"
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
	error := runner.Run()

	// Then
	assert.NoError(t, error)
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
	error := runner.Run()

	// Then
	assert.EqualError(t, expectedError, error.Error())
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
	error := runner.Run()

	// Then
	assert.NoError(t, error)
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
	error := runner.Run()

	// Then
	assert.NoError(t, error)
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
	error := runner.Run()

	// Then
	assert.NoError(t, error)
	mockCarthageCache.AssertNotCalled(t, "CreateIndicator")
	mockCarthageCache.AssertNumberOfCalls(t, "Commit", 1)
}

// Retry on failure
func Test_GivenBootstrapCommandAndSingleNetworkFailure_WhenRunCalled_ThenExpectCommandToBeRetriedAndSucceed(t *testing.T) {
	// Given
	blueprints := []CommandBlueprint{
		{
			Command:   "bash",
			Arguments: []string{"-c", failingCommandWithTimeoutStderr},
		},
		{
			Command:   "echo",
			Arguments: []string{"hello"},
		},
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("bootstrap", blueprints)

	// When
	error := runner.Run()

	// Then
	assert.NoError(t, error)
}

func Test_GivenBootstrapCommandAndPermanentNetworkFailure_WhenRunCalled_ThenExpectCommandToFail(t *testing.T) {
	// Given
	blueprints := []CommandBlueprint{
		{
			Command:   "bash",
			Arguments: []string{"-c", failingCommandWithTimeoutStderr},
		},
		{
			Command:   "bash",
			Arguments: []string{"-c", failingCommandWithFailedToConnectToStderr},
		},
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("bootstrap", blueprints)

	// When
	error := runner.Run()

	// Then
	assert.Error(t, error)
}

func Test_GivenUpdateCommandAndSingleNetworkFailure_WhenRunCalled_ThenExpectCommandToBeRetriedAndSucceed(t *testing.T) {
	// Given
	blueprints := []CommandBlueprint{
		{
			Command:   "bash",
			Arguments: []string{"-c", failingCommandWithFailedToConnectToStderr},
		},
		{
			Command:   "echo",
			Arguments: []string{"hello"},
		},
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("update", blueprints)

	// When
	error := runner.Run()

	// Then
	assert.NoError(t, error)
}

func Test_GivenUpdateCommandAndPermanentNetworkFailure_WhenRunCalled_ThenExpectCommandToFail(t *testing.T) {
	// Given
	blueprints := []CommandBlueprint{
		{
			Command:   "bash",
			Arguments: []string{"-c", failingCommandWithFailedToConnectToStderr},
		},
		{
			Command:   "bash",
			Arguments: []string{"-c", failingCommandWithTimeoutStderr},
		},
	}
	runner := givenRunnerWithMainAndCommandBuilderCommands("update", blueprints)

	// When
	error := runner.Run()

	// Then
	assert.Error(t, error)
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
	expectedAvailable := true
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
	command := "version"
	args := []string{"arg1"}
	runner := Runner{
		carthageCommand: command,
		args:            args,
		commandBuilder:  mockCommandBuilder,
	}

	// When
	error := runner.executeCommand()

	// Then
	assert.NoError(t, error)
	mockCommandBuilder.AssertCalled(t, "AddGitHubToken", mock.Anything)
	mockCommandBuilder.AssertCalled(t, "AddXCConfigFile", mock.Anything)
	mockCommandBuilder.AssertCalled(t, "Append", []string{command})
	mockCommandBuilder.AssertCalled(t, "Append", args)
}

// helpers
func givenMockCarthageCache() *MockCarthageCache {
	return new(MockCarthageCache)
}

func givenStubbedCommandBuilder() *MockCommandBuilder {
	blueprint := CommandBlueprint{
		Command:   "echo",
		Arguments: []string{"hello"},
	}
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(blueprint)
	return mockCommandBuilder
}

func givenStubbedCommandBuilderReturnFailingCommand() *MockCommandBuilder {
	blueprint := CommandBlueprint{
		Command:   "fail",
		Arguments: nil,
	}
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(blueprint)
	return mockCommandBuilder
}

func givenRunnerWithMainAndCommandBuilderCommands(mainCommand string, commandBlueprints []CommandBlueprint) Runner {
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(false).
		GivenCreateIndicatorSucceeds().
		GivenCommitSucceeds()

	return Runner{
		carthageCommand: mainCommand,
		cache:           mockCarthageCache,
		commandBuilder:  givenStubbedCommandBuilderReturnsCommands(commandBlueprints),
	}
}

func givenStubbedCommandBuilderReturnsCommands(commandBlueprints []CommandBlueprint) *MockCommandBuilder {
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandsReturned(commandBlueprints)
	return mockCommandBuilder
}