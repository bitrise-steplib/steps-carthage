package cachedcarthage

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-steplib/steps-carthage/carthage"
	"github.com/hashicorp/go-version"
)

const (
	bootstrapCommand = "bootstrap"
	projectDirArg    = "--project-directory"
	carthageDirName  = "Carthage"
	buildDirName     = "Build"
	cacheFileName    = "Cachefile"
	resolvedFileName = "Cartfile.resolved"
)

// Runner ...
type Runner struct {
	command           string
	args              []string
	sourceDir         string
	githubAccessToken stepconf.Secret
}

// NewRunner ...
func NewRunner(
	command string,
	args []string,
	sourceDir string,
	githubAccessToken stepconf.Secret,
) Runner {
	return Runner{
		command:           command,
		args:              args,
		sourceDir:         sourceDir,
		githubAccessToken: githubAccessToken,
	}
}

// Run ...
func (runner Runner) Run() error {
	// Environment
	fmt.Println()
	log.Infof("Environment:")

	carthageVersion, err := runner.getCarthageVersion()
	if err != nil {
		// TODO fail("Failed to get carthage version, error: %s", err)
		return err
	}
	log.Printf("- CarthageVersion: %s", carthageVersion.String())

	swiftVersion, err := runner.getSwiftVersion()
	if err != nil {
		// TODO fail("Failed to get swift version, error: %s", err)
		return err
	}
	log.Printf("- SwiftVersion: %s", strings.Replace(swiftVersion, "\n", "- ", -1))
	// --

	projectDir := runner.sourceDir

	isNextOptionProjectDir := false
	for _, option := range runner.args {
		if option == projectDirArg {
			isNextOptionProjectDir = true
			continue
		}

		if isNextOptionProjectDir {
			projectDir = option

			fmt.Println()
			log.Infof("--project-directory flag found with value: %s", projectDir)
			log.Printf("using %s as working directory", projectDir)

			break
		}
	}
	// ---

	// Exit if bootstrap is cached
	if runner.command == bootstrapCommand {
		fmt.Println()
		log.Infof("Check if cache is available")

		cacheAvailable, err := runner.isCacheAvailable(projectDir, swiftVersion)
		if err != nil {
			log.Warnf("Failed to check if cached is available, error: %s", err)
		}

		if cacheAvailable {
			log.Successf("Cache available")
		} else {
			log.Errorf("Cache not available")
		}

		if cacheAvailable {
			if err := runner.collectCarthageCache(projectDir); err != nil {
				log.Warnf("Cache collection skipped: %s", err)
			}
			log.Donef("Using cached dependencies for bootstrap command. If you would like to force update your dependencies, select `update` as CarthageCommand and re-run your build.")
			os.Exit(0)
		}
	}
	// ---

	fmt.Println()
	log.Infof("Running Carthage command")

	cmd := carthage.New().AddGitHubToken(runner.githubAccessToken).Command(runner.args...)

	// TODO
	// log.Printf("Appending GITHUB_ACCESS_TOKEN to process environments")

	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	log.Donef("$ %s", cmd.PrintableCommandArgs())
	fmt.Println()

	if err := cmd.Run(); err != nil {
		return err
		// TODO: fail("Carthage command failed, error: %s", err)
	}
	// ---

	// Create cache
	if runner.command == bootstrapCommand {
		fmt.Println()
		log.Infof("Creating cache")

		cacheFilePth := filepath.Join(projectDir, carthageDirName, cacheFileName)

		resolvedFilePath := filepath.Join(projectDir, resolvedFileName)
		resolvedFileContent, err := runner.contentOfFile(resolvedFilePath)
		if err != nil {
			// TODO fail("Failed to read Cartfile.resolved, error: %s", err)
			return err
		}
		if resolvedFileContent == "" {
			log.Warnf("Cartfile.resolved is empty or not exists at: %s", resolvedFilePath)
			log.Warnf("Skipping caching")
			os.Exit(1)
		}

		cacheContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s", swiftVersion, resolvedFileName, resolvedFileContent, resolvedFileName)

		carthageDir := filepath.Join(projectDir, carthageDirName)
		if exist, err := pathutil.IsPathExists(carthageDir); err != nil {
			// TODO fail("Failed to check if dir exists at (%s), error: %s", carthageDir, err)
			return err
		} else if !exist {
			if err := os.Mkdir(carthageDir, 0777); err != nil {
				// TODO fail("Failed to create dir (%s), error: %s", carthageDir, err)
				return err
			}
		}

		if err := fileutil.WriteStringToFile(cacheFilePth, cacheContent); err != nil {
			// TODO fail("Failed to write cahe file, error: %s", err)
			return err
		}

		log.Donef("Cachefile created: %s", cacheFilePth)

		if err := runner.collectCarthageCache(projectDir); err != nil {
			log.Warnf("Cache collection skipped: %s", err)
		}
	}
	// ---

	return nil
}

func (runner Runner) getCarthageVersion() (*version.Version, error) {
	cmd := carthage.New().Version()
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return nil, err
	}

	// if output is multi-line, get the last line of string
	// parse Version from cmd output
	for _, outLine := range strings.Split(out, "\n") {
		if currentVersion, err := version.NewVersion(outLine); err == nil {
			return currentVersion, nil
		}
	}

	return nil, fmt.Errorf("failed to parse `$ carthage version` output: %s", out)
}

func (runner Runner) getSwiftVersion() (string, error) {
	cmd := command.New("swift", "-version")
	return cmd.RunAndReturnTrimmedCombinedOutput()
}

func (runner Runner) isCacheAvailable(srcDir string, swiftVersion string) (bool, error) {
	log.Printf("Check cache specific files")
	buildDirAvailable, cacheFileAvailable, resolvedFileAvailable := true, true, true

	// check for built dependencies (Carthage/Build/*)
	fmt.Print("- Carthage/Build directory: ")
	carthageDir := filepath.Join(srcDir, carthageDirName)
	carthageBuildDir := filepath.Join(carthageDir, buildDirName)

	files, err := ioutil.ReadDir(carthageBuildDir)
	if err != nil {
		buildDirAvailable = false
		log.Errorf("not found")
	} else if len(files) == 0 {
		buildDirAvailable = false
		log.Errorf("empty")
	} else {
		log.Successf("found")
	}

	// ---

	// read cache indicator file (Carthage/Cachefile)
	fmt.Print("- Cachefile: ")
	var cacheFileContent string

	cacheFilePth := filepath.Join(carthageDir, cacheFileName)
	if exist, err := pathutil.IsPathExists(cacheFilePth); err != nil {
		return false, err
	} else if !exist {
		cacheFileAvailable = false
		log.Errorf("not available yet")
	} else {
		log.Successf("found")
		log.Debugf(cacheFileName + " exists in " + cacheFilePth + "\n")

		cacheFileContent, err = runner.contentOfFile(cacheFilePth)
		if err != nil {
			return false, err
		} else if cacheFileContent == "" {
			log.Errorf("Cachefile is empty")
		} else {
			log.Debugf(cacheFileName + " content: " + cacheFileContent + "\n")
		}
	}

	// --

	// read Cartfile.resolved
	fmt.Print("- Cartfile.resolved: ")
	var resolvedFileContent string

	resolvedFilePath := filepath.Join(srcDir, resolvedFileName)
	if exist, err := pathutil.IsPathExists(resolvedFilePath); err != nil {
		return false, err
	} else if !exist {
		log.Errorf("not found")
		resolvedFileAvailable = false
	} else {
		log.Successf("found")

		resolvedFileContent, err = runner.contentOfFile(resolvedFilePath)
		if err != nil {
			return false, err

		} else if resolvedFileContent == "" {
			return false, fmt.Errorf("Catfile.resolved is empty")
		}
		log.Debugf(resolvedFileName + " content: " + resolvedFileContent + "\n")
	}

	// ---

	// Warn messages about the missing files

	// Print the warning about the missing Cachefile only if the other required file (Cartfile.resolved) is available.
	// If the Cartfile.resolved is not found, then we don't want to mislead the user with this warning.
	if !cacheFileAvailable && resolvedFileAvailable {
		fmt.Println()
		log.Warnf("The " + cacheFileName + " is generated by the step. Probably cache not initialised yet (first cache push initialises the cache), nothing to worry about ;)")

	}

	if !resolvedFileAvailable {
		fmt.Println()
		log.Warnf("No "+resolvedFileName+" found at: %s", resolvedFilePath)
		log.Warnf("Make sure it's committed into your repository!")
		log.Warnf(resolvedFileName + " presence ensures that Bitrise will use exactly the same versions of dependencies as you in your local environment. ")
		log.Warnf("The dependencies will not be cached until the " + resolvedFileName + " file presents in the repository.")
	}

	if buildDirAvailable && cacheFileAvailable && resolvedFileAvailable {
		desiredCacheContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s", swiftVersion, resolvedFileName, resolvedFileContent, resolvedFileName)
		if cacheFileContent != desiredCacheContent {
			log.Debugf(
				"Cachefile is not valid.\n" +
					"Desired cache content:\n" +
					desiredCacheContent +
					"CacheFile content:\n" +
					cacheFileContent,
			)

			return false, nil
		}
	}

	return (buildDirAvailable && cacheFileAvailable && resolvedFileAvailable), nil
}

func (runner Runner) contentOfFile(pth string) (string, error) {
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("file not exists: %s", pth)
	}

	return fileutil.ReadStringFromFile(pth)
}

func (runner Runner) collectCarthageCache(projectDir string) error {
	fmt.Println()
	log.Infof("Collecting carthage caches...")

	absCarthageDir, err := filepath.Abs(filepath.Join(projectDir, carthageDirName))
	if err != nil {
		return fmt.Errorf("failed to determine cache paths")
	}
	absCacheFilePth, err := filepath.Abs(filepath.Join(projectDir, carthageDirName, cacheFileName))
	if err != nil {
		return fmt.Errorf("failed to determine cache paths")
	}
	carthageCache := cache.New()
	carthageCache.IncludePath(fmt.Sprintf("%s -> %s", absCarthageDir, absCacheFilePth))
	if err := carthageCache.Commit(); err != nil {
		return fmt.Errorf("failed to commit cache paths")
	}

	return nil
}
