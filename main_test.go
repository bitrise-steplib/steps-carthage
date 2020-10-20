package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// parseProjectDir
func Test_GivenCustomDirProvided_WhenParseProjectDirCalled_ThenExpectCustomDir(t *testing.T) {
	// Given
	expectedDir := "/customDir"
	customOptions := []string{"--project-directory", expectedDir}

	// When
	acutalProjectDir := parseProjectDir("/originalDir", customOptions)

	//Then
	assert.Equal(t, expectedDir, acutalProjectDir)
}

func Test_GivenCustomDirNotProvided_WhenParseProjectDirCalled_ThenExpectOriginalDir(t *testing.T) {
	// Given
	expectedDir := "/originalDir"
	customOptions := []string{"--someparam", `"some value"`}

	// When
	acutalProjectDir := parseProjectDir(expectedDir, customOptions)

	//Then
	assert.Equal(t, expectedDir, acutalProjectDir)
}

// parseCarthageOptions
func Test_WhenParseCarthageOptionsCalled_ThenExpectCorrectValue(t *testing.T) {
	// Given
	expectedOpts := []string{"--parameter", "value"}
	options := Config{
		CarthageOptions: "--parameter value",
	}

	// When
	actualOpts := parseCarthageOptions(options)

	// Then
	assert.Equal(t, expectedOpts, actualOpts)
}

// parseXCConfigPath
func Test_GivenXCConfigAsInputAndFileProviderSucceeds_WhenParseXCConfigPathCalled_ThenExpectPath(t *testing.T) {
	// Given
	expectedPath := "/path/from/input.xcconfig"
	mockFileProvider := givenMockFileProvider().
		GivenLocalPathSucceeds(expectedPath)

	// When
	actualPath, err := parseXCConfigPath(expectedPath, "", mockFileProvider)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)
}

func Test_GivenXCConfigAsInputAndFileProviderFails_WhenParseXCConfigPathCalled_ThenExpectError(t *testing.T) {
	// Given
	expectedError := errors.New("sad error")
	mockFileProvider := givenMockFileProvider().
		GivenLocalPathFails(expectedError)

	// When
	actualPath, actualErr := parseXCConfigPath("whatever", "", mockFileProvider)

	// Then
	assert.EqualError(t, expectedError, actualErr.Error())
	assert.Empty(t, actualPath)
}

func Test_GivenBothXCConfigAsInputAndEnvPassed_WhenParseXCConfigPathCalled_ThenExpectPathFromInput(t *testing.T) {
	// Given
	expectedPath := "/path/from/input.xcconfig"
	envPath := "/path/from/env.xcconfig"
	mockFileProvider := givenMockFileProvider().
		GivenLocalPathSucceeds(expectedPath)

	// When
	actualPath, err := parseXCConfigPath(expectedPath, envPath, mockFileProvider)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)
}

func Test_GivenXCConfigAsEnvPassed_WhenParseXCConfigPathCalled_ThenExpectPath(t *testing.T) {
	// Given
	expectedPath := "/path/from/env.xcconfig"

	// When
	actualPath, err := parseXCConfigPath("", expectedPath, nil)

	// Then
	assert.NoError(t, err)
	assert.Equal(t, expectedPath, actualPath)
}

func givenMockFileProvider() *MockFileProvider {
	return new(MockFileProvider)
}
