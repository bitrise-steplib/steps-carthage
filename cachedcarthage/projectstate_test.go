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
		{ProjectState{buildDirExists: true, cacheFileExists: true, resolvedFileExists: true}, true},
		{ProjectState{buildDirExists: false, cacheFileExists: true, resolvedFileExists: true}, false},
		{ProjectState{buildDirExists: true, cacheFileExists: false, resolvedFileExists: true}, false},
		{ProjectState{buildDirExists: true, cacheFileExists: true, resolvedFileExists: false}, false},
		{ProjectState{buildDirExists: false, cacheFileExists: false, resolvedFileExists: true}, false},
		{ProjectState{buildDirExists: true, cacheFileExists: false, resolvedFileExists: false}, false},
		{ProjectState{buildDirExists: false, cacheFileExists: true, resolvedFileExists: false}, false},
		{ProjectState{buildDirExists: false, cacheFileExists: false, resolvedFileExists: false}, false},
	}

	for _, scenario := range testScenarios {
		// When
		actual := scenario.state.isCacheIntact()

		// Then
		assert.Equal(t, scenario.expected, actual)
	}
}
