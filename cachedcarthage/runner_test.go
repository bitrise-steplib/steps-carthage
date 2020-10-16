package cachedcarthage

import (
	"errors"
	"testing"

	"github.com/bitrise-io/go-utils/command"
	"github.com/stretchr/testify/assert"
	mock "github.com/stretchr/testify/mock"
)

// Run
func Test_GivenNotBootstratpedCommand_WhenRunCalled_ThenExpectNoErrorAndCacheNotCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache()
	mockCommandBuilder := givenStubedCommandBuilder()
	runner := Runner{
		command:        "version",
		cache:          mockCarthageCache,
		commandBuilder: mockCommandBuilder,
	}

	// When
	error := runner.Run()

	// Then
	assert.NoError(t, error)
	mockCarthageCache.AssertNotCalled(t, "IsAvailable")
	mockCarthageCache.AssertNotCalled(t, "Collect")
	mockCarthageCache.AssertNotCalled(t, "Create")
}

func Test_GivenBootstratpedCommandAndCacheNotAvailableAndCacheCreateFails_WhenRunCalled_ThenExpectError(t *testing.T) {
	// Given
	expectedError := errors.New("sad error")
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(false).
		GivenCreateFails(expectedError)
	runner := Runner{
		command:        "bootstrap",
		cache:          mockCarthageCache,
		commandBuilder: givenStubedCommandBuilder(),
	}

	// When
	error := runner.Run()

	// Then
	assert.EqualError(t, expectedError, error.Error())
}

func Test_GivenBootstratpedCommandAndCacheNotAvailableAndCacheCreateSucceeds_WhenRunCalled_ThenExpectNoError(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(false).
		GivenCreateSucceeds().
		GivenCollectSucceeds()
	runner := Runner{
		command:        "bootstrap",
		cache:          mockCarthageCache,
		commandBuilder: givenStubedCommandBuilder(),
	}

	// When
	error := runner.Run()

	// Then
	assert.NoError(t, error)
}

func Test_GivenBootstratpedCommandAndCacheAvailableAndCollectFails_WhenRunCalled_ThenExpectCommandExecutedCacheCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(true).
		GivenCollectFails(errors.New("sad error")).
		GivenCreateSucceeds()

	mockCommandBuilder := givenStubedCommandBuilder()

	runner := Runner{
		command:        "bootstrap",
		cache:          mockCarthageCache,
		commandBuilder: mockCommandBuilder,
	}

	// When
	error := runner.Run()

	// Then
	assert.NoError(t, error)
	mockCarthageCache.AssertCalled(t, "Create")
	mockCarthageCache.AssertCalled(t, "Collect")
}

func Test_GivenBootstratpedCommandAndCacheAvailableAndCollectSucceeds_WhenRunCalled_ThenExpectCommandNotExecutedAndCacheCreated(t *testing.T) {
	// Given
	mockCarthageCache := givenMockCarthageCache().
		GivenIsAvailableSucceeds(true).
		GivenCollectSucceeds()

	mockCommandBuilder := givenStubedCommandBuilderReturnFailingCommand()
	runner := Runner{
		command:        "bootstrap",
		cache:          mockCarthageCache,
		commandBuilder: mockCommandBuilder,
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
	mockCommandBuilder := givenStubedCommandBuilder()
	command := "version"
	args := []string{"arg1"}
	runner := Runner{
		command:        command,
		args:           args,
		commandBuilder: mockCommandBuilder,
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

func givenStubedCommandBuilder() *MockCommandBuilder {
	command := command.New("echo", "hello")
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(command)
	return mockCommandBuilder
}

func givenStubedCommandBuilderReturnFailingCommand() *MockCommandBuilder {
	command := command.New("fail")
	mockCommandBuilder := new(MockCommandBuilder).
		GivenAddGitHubTokenSucceeds().
		GivenAddXCConfigFileSucceeds().
		GivenAppendSucceeds().
		GivenCommandReturned(command)
	return mockCommandBuilder
}
