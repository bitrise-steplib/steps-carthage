package cachedcarthage

import "path/filepath"

const (
	carthageDirName  = "Carthage"
	buildDirName     = "Build"
	resolvedFileName = "Cartfile.resolved"
	cacheFileName    = "Cachefile"
)

// Project ...
type Project struct {
	projectDir string
}

// NewProject ...
func NewProject(projectDir string) Project {
	return Project{projectDir: projectDir}
}

func (project Project) carthageDir() string {
	return filepath.Join(project.projectDir, carthageDirName)
}

func (project Project) cacheFilePath() string {
	return filepath.Join(project.carthageDir(), cacheFileName)
}

func (project Project) buildDir() string {
	return filepath.Join(project.carthageDir(), buildDirName)
}

func (project Project) resolvedFilePath() string {
	return filepath.Join(project.projectDir, resolvedFileName)
}
