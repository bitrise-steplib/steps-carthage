package cachedcarthage

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
)

// DefaultStateProvider reads the current state of a cached Carthage project.
type DefaultStateProvider struct {
}

// ParseState ...
func (provider DefaultStateProvider) ParseState(project Project) (ProjectState, error) {
	buildDirNotEmpty, buildDirFiles := provider.parseBuildDirectoryState(project.buildDir())
	cacheFileExists, cacheFileContent, err := provider.parseCacheFileState(project.cacheFilePath())
	if err != nil {
		return ProjectState{}, err
	}
	resolvedFileExists, resolvedFileContent, err := provider.parseResolvedFileState(project.resolvedFilePath())
	if err != nil {
		return ProjectState{}, err
	}

	carthageDirExists, err := pathutil.IsPathExists(project.carthageDir())
	if err != nil {
		return ProjectState{}, fmt.Errorf("failed to check if dir exists at (%s), error: %s", project.carthageDir(), err)
	}

	return ProjectState{
		buildDirNotEmpty: buildDirNotEmpty && len(buildDirFiles) != 0,

		cacheFileExists:  cacheFileExists,
		cacheFileContent: cacheFileContent,

		resolvedFileExists:  resolvedFileExists,
		resolvedFileContent: resolvedFileContent,

		carthageDirExists: carthageDirExists,
	}, nil
}

func (provider DefaultStateProvider) parseBuildDirectoryState(buildDir string) (bool, []os.FileInfo) {
	files, err := ioutil.ReadDir(buildDir)
	if err != nil {
		return false, nil
	}

	return true, files
}

func (provider DefaultStateProvider) parseCacheFileState(path string) (bool, string, error) {
	if exist, err := pathutil.IsPathExists(path); err != nil {
		return false, "", err
	} else if !exist {
		return false, "", nil
	}

	cacheFileContent, err := provider.contentOfFile(path)
	if err != nil {
		return false, "", err
	}

	return true, cacheFileContent, nil
}

func (provider DefaultStateProvider) parseResolvedFileState(path string) (bool, string, error) {
	if exist, err := pathutil.IsPathExists(path); err != nil {
		return false, "", err
	} else if !exist {
		return false, "", nil
	}

	resolvedFileContent, err := provider.contentOfFile(path)
	if err != nil {
		return false, "", err
	} else if resolvedFileContent == "" {
		return false, "", fmt.Errorf("Catfile.resolved is empty")
	}

	return true, resolvedFileContent, nil
}

func (provider DefaultStateProvider) contentOfFile(pth string) (string, error) {
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("file does not exist: %s", pth)
	}

	return fileutil.ReadStringFromFile(pth)
}
