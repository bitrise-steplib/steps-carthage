package cachedcarthage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_WhenIsCacheIntactCalled_ThenExpectCorrectValue(t *testing.T) {
	testScenarios := []struct {
		state    ProjectState
		expected bool
	}{
		{ProjectState{buildDirNotEmpty: true, cacheFileExists: true, resolvedFileExists: true}, true},
		{ProjectState{buildDirNotEmpty: false, cacheFileExists: true, resolvedFileExists: true}, false},
		{ProjectState{buildDirNotEmpty: true, cacheFileExists: false, resolvedFileExists: true}, false},
		{ProjectState{buildDirNotEmpty: true, cacheFileExists: true, resolvedFileExists: false}, false},
		{ProjectState{buildDirNotEmpty: false, cacheFileExists: false, resolvedFileExists: true}, false},
		{ProjectState{buildDirNotEmpty: true, cacheFileExists: false, resolvedFileExists: false}, false},
		{ProjectState{buildDirNotEmpty: false, cacheFileExists: true, resolvedFileExists: false}, false},
		{ProjectState{buildDirNotEmpty: false, cacheFileExists: false, resolvedFileExists: false}, false},
	}

	for _, scenario := range testScenarios {
		// When
		actual := scenario.state.isCacheIntact()

		// Then
		assert.Equal(t, scenario.expected, actual)
	}
}
