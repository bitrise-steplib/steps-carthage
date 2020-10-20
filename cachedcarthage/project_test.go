package cachedcarthage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WhenCarthageDirCalled_ThenExpectCorrectPath(t *testing.T) {
	// Given
	expectedPath := "/base/dir/Carthage"
	project := Project{"/base/dir"}

	// When
	actualPath := project.carthageDir()

	// Then
	assert.Equal(t, expectedPath, actualPath)
}

func Test_WhenCacheFilePathCalled_ThenExpectCorrectPath(t *testing.T) {
	// Given
	expectedPath := "/base/dir/Carthage/Cachefile"
	project := Project{"/base/dir"}

	// When
	actualPath := project.cacheFilePath()

	// Then
	assert.Equal(t, expectedPath, actualPath)
}

func Test_WhenBuildDirCalled_ThenExpectCorrectPath(t *testing.T) {
	// Given
	expectedPath := "/base/dir/Carthage/Build"
	project := Project{"/base/dir"}

	// When
	actualPath := project.buildDir()

	// Then
	assert.Equal(t, expectedPath, actualPath)
}

func Test_WhenResolvedFilePathCalled_ThenExpectCorrectPath(t *testing.T) {
	// Given
	expectedPath := "/base/dir/Cartfile.resolved"
	project := Project{"/base/dir"}

	// When
	actualPath := project.resolvedFilePath()

	// Then
	assert.Equal(t, expectedPath, actualPath)
}
