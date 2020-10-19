package cachedcarthage

import "os"

// ProjectState represents a snapshot of a cached Carthage project.
type ProjectState struct {
	buildDirExists bool
	buildDirFiles  []os.FileInfo

	cacheFileExists  bool
	cacheFileContent string

	resolvedFileExists  bool
	resolvedFileContent string

	carthageDirExists bool
}

func (state ProjectState) isCacheIntact() bool {
	return state.buildDirExists && state.cacheFileExists && state.resolvedFileExists
}
