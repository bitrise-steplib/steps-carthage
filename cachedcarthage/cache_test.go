package cachedcarthage

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Create
func Test_GivenStateCouldNotBeParsed_WhenCreateCalled_ThenExpectError(t *testing.T) {
	// Given
	expectedError := errors.New("sad error")
	mockStateProvider := givenMockProjectStateProvider().GivenParseStateFails(expectedError)
	mockFileCache := givenMockFileCache()

	cache := Cache{
		project:       Project{},
		swiftVersion:  "whatever",
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualError := cache.Create()

	// Then
	assert.EqualError(t, expectedError, actualError.Error())
}

func Test_GivenCarthageDirDoesNotExist_WhenCreateCalled_ThenExpectCarthageDirAndCacheFileToExist(t *testing.T) {
	// Given
	tempDir := givenTempDir(t)
	expectedDir := filepath.Join(tempDir, "Carthage")
	expectedFile := filepath.Join(expectedDir, "Cachefile")
	defer func() {
		err := os.RemoveAll(tempDir)
		require.NoError(t, err)
	}()
	project := Project{tempDir}
	mockStateProvider := givenMockProjectStateProvider().GivenParseStateSucceeds(ProjectState{})
	mockFileCache := givenMockFileCache()

	cache := Cache{
		project:       project,
		swiftVersion:  "whatever",
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualError := cache.Create()

	// Then
	assert.NoError(t, actualError)
	assert.DirExists(t, expectedDir)
	assert.FileExists(t, expectedFile)
}

// cacheFileContent
func Test_WhenCacheFileContentCalled_ThenExpectCorrectValue(t *testing.T) {
	// Given
	resolvedFileName := "Cartfile.resolved"
	content := "nice content"
	swiftVersion := "5.0.2"

	expectedContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s",
		swiftVersion,
		resolvedFileName,
		content,
		resolvedFileName)

	mockStateProvider := givenMockProjectStateProvider()
	mockFileCache := givenMockFileCache()

	cache := Cache{
		project:       Project{},
		swiftVersion:  swiftVersion,
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualContent := cache.cacheFileContent(content)

	// Then
	assert.Equal(t, expectedContent, actualContent)
}

// Collect
func Test_GivenFileCacheCommitFails_WhenCollectCalled_ThenExpectError(t *testing.T) {
	// Given
	expectedError := errors.New("failed to commit cache paths")
	mockStateProvider := givenMockProjectStateProvider()
	mockFileCache := givenMockFileCache().
		GivenIncludeSucceeds().
		GivenCommitFails(expectedError)
	cache := Cache{
		project:       Project{},
		swiftVersion:  "whatever",
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualError := cache.Collect()

	// Then
	assert.EqualError(t, expectedError, actualError.Error())
}

func Test_GivenFileCacheCommitSucceeds_WhenCollectCalled_ThenExpectIncludePathCalledWithCorrectValue(t *testing.T) {
	// Given
	projectDir := "/awesomepath"
	expectedCacheCall := fmt.Sprintf(
		"%s -> %s",
		filepath.Join(projectDir, "Carthage"),
		filepath.Join(projectDir, "Carthage/Cachefile"))
	mockStateProvider := givenMockProjectStateProvider()
	mockFileCache := givenMockFileCache().
		GivenIncludeSucceeds().
		GivenCommitSucceeds()
	cache := Cache{
		project:       Project{projectDir},
		swiftVersion:  "whatever",
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualError := cache.Collect()

	// Then
	assert.NoError(t, actualError)
	mockFileCache.AssertCalled(t, "IncludePath", expectedCacheCall)
}

// IsAvailable
func Test_GivenStateCouldNotBeParsed_WhenIsAvailableCalled_ThenExpectError(t *testing.T) {
	// Given
	expectedError := errors.New("sad error")
	mockStateProvider := givenMockProjectStateProvider().GivenParseStateFails(expectedError)
	mockFileCache := givenMockFileCache()

	cache := Cache{
		project:       Project{},
		swiftVersion:  "whatever",
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualValue, actualError := cache.IsAvailable()

	// Then
	assert.EqualError(t, expectedError, actualError.Error())
	assert.False(t, actualValue)
}

func Test_GivenStateIsNotIntact_WhenIsAvailableCalled_ThenExpectFalse(t *testing.T) {
	// Given
	state := ProjectState{}
	mockStateProvider := givenMockProjectStateProvider().GivenParseStateSucceeds(state)
	mockFileCache := givenMockFileCache()

	cache := Cache{
		project:       Project{},
		swiftVersion:  "whatever",
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualValue, err := cache.IsAvailable()

	// Then
	assert.NoError(t, err)
	assert.False(t, actualValue)
}

func Test_GivenStateIsIntactButCacheFileIsCorrupt_WhenIsAvailableCalled_ThenExpectTrue(t *testing.T) {
	// Given
	state := ProjectState{
		buildDirExists:     true,
		cacheFileExists:    true,
		cacheFileContent:   "corrupt",
		resolvedFileExists: true,
	}
	mockStateProvider := givenMockProjectStateProvider().GivenParseStateSucceeds(state)
	mockFileCache := givenMockFileCache()

	cache := Cache{
		project:       Project{},
		swiftVersion:  "whatever",
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualValue, err := cache.IsAvailable()

	// Then
	assert.NoError(t, err)
	assert.False(t, actualValue)
}

func Test_GivenStateIsIntactAndCacheFileIsCorrect_WhenIsAvailableCalled_ThenExpectTrue(t *testing.T) {
	// Given
	resolvedFileName := "Cartfile.resolved"
	resolvedContent := "nice content"
	swiftVersion := "5.0.2"

	expectedContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s",
		swiftVersion,
		resolvedFileName,
		resolvedContent,
		resolvedFileName)

	state := ProjectState{
		buildDirExists:      true,
		cacheFileExists:     true,
		cacheFileContent:    expectedContent,
		resolvedFileExists:  true,
		resolvedFileContent: resolvedContent,
	}
	mockStateProvider := givenMockProjectStateProvider().GivenParseStateSucceeds(state)
	mockFileCache := givenMockFileCache()

	cache := Cache{
		project:       Project{},
		swiftVersion:  swiftVersion,
		filecache:     mockFileCache,
		stateProvider: mockStateProvider,
	}

	// When
	actualValue, err := cache.IsAvailable()

	// Then
	assert.NoError(t, err)
	assert.True(t, actualValue)
}

// helpers
func givenMockProjectStateProvider() *MockProjectStateProvider {
	return new(MockProjectStateProvider)
}

func givenMockFileCache() *MockFileCache {
	return new(MockFileCache)
}

func givenTempDir(t *testing.T) string {
	path, err := pathutil.NormalizedOSTempDirPath("test")
	require.NoError(t, err)
	return path
}
