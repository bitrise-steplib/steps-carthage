package cachedcarthage

import (
	"errors"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
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
	mockCarthageCache.AssertNotCalled(t, "Collect")
	mockCarthageCache.AssertNotCalled(t, "Create")
}

func Test_GivenBootstrapCommandAndCacheNotAvailableAndCacheCreateFails_WhenRunCalled_ThenExpectError(t *testing.T) {
	// Given
	expectedError := errors.New("sad error")
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(false).
		GivenCreateFails(expectedError)
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
		GivenCreateSucceeds().
		GivenCollectSucceeds()
	runner := Runner{
		carthageCommand: "bootstrap",
		cache:           mockCarthageCache,
		commandBuilder:  givenStubbedCommandBuilder(),
	}

	// When
	error := runner.Run()

	// Then
	assert.NoError(t, error)
	mockCarthageCache.AssertCalled(t, "Create")
	mockCarthageCache.AssertCalled(t, "Collect")
}

func Test_GivenBootstrapCommandAndCacheAvailableAndCollectFails_WhenRunCalled_ThenExpectCommandExecutedCacheCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(true).
		GivenCollectFails(errors.New("sad error")).
		GivenCreateSucceeds()

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
	mockCarthageCache.AssertCalled(t, "Create")
	mockCarthageCache.AssertCalled(t, "Collect")
}

func Test_GivenBootstrapCommandAndCacheAvailableAndCollectSucceeds_WhenRunCalled_ThenExpectCommandNotExecutedAndCacheCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(true).
		GivenCollectSucceeds()

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
	mockCarthageCache.AssertNotCalled(t, "Create")
	mockCarthageCache.AssertNumberOfCalls(t, "Collect", 1)
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
	command := command.New("echo", "hello")
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(command)
	return mockCommandBuilder
}

func givenStubbedCommandBuilderReturnFailingCommand() *MockCommandBuilder {
	command := command.New("fail")
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(command)
	return mockCommandBuilder
}
