package cachedcarthage

// ProjectState represents a snapshot of a cached Carthage project.
type ProjectState struct {
	buildDirNotEmpty bool

	cacheFileExists  bool
	cacheFileContent string

	resolvedFileExists  bool
	resolvedFileContent string

	carthageDirExists bool
}

func (state ProjectState) isCacheIntact() bool {
	return state.buildDirNotEmpty &&
		state.cacheFileExists &&
		state.resolvedFileExists
}
