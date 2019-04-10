package main

// ConfigsModel ...
import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-tools/go-steputils/cache"
	"github.com/bitrise-tools/go-steputils/stepconf"
	version "github.com/hashicorp/go-version"
	"github.com/kballard/go-shellquote"
)

const (
	carthageDirName  = "Carthage"
	buildDirName     = "Build"
	cacheFileName    = "Cachefile"
	resolvedFileName = "Cartfile.resolved"

	bootstrapCommand = "bootstrap"
)

// Config ...
type Config struct {
	GithubAccessToken stepconf.Secret `env:"github_access_token"`
	CarthageCommand   string          `env:"carthage_command,required"`
	CarthageOptions   string          `env:"carthage_options"`
	SourceDir         string          `env:"BITRISE_SOURCE_DIR"`
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func getSwiftVersion() (string, error) {
	cmd := command.New("swift", "-version")
	return cmd.RunAndReturnTrimmedCombinedOutput()
}

func getCarthageVersion() (*version.Version, error) {
	cmd := command.New("carthage", "version")
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

func contentOfFile(pth string) (string, error) {
	if exist, err := pathutil.IsPathExists(pth); err != nil {
		return "", err
	} else if !exist {
		return "", fmt.Errorf("File not exists: %s", pth)
	}

	return fileutil.ReadStringFromFile(pth)
}

func isDirEmpty(dirPth string) (bool, error) {
	if exist, err := pathutil.IsDirExists(dirPth); err != nil {
		return false, err
	} else if !exist {
		return false, nil
	}

	pattern := filepath.Join(dirPth, "*")
	files, err := filepath.Glob(pattern)
	if err != nil {
		return false, err
	}

	return (len(files) == 0), nil
}

func isCacheAvailable(srcDir string, swiftVersion string) (bool, error) {
	// check for built dependencies (Carthage/Build/*)
	carthageDir := filepath.Join(srcDir, carthageDirName)
	carthageBuildDir := filepath.Join(carthageDir, buildDirName)
	if empty, err := isDirEmpty(carthageBuildDir); err != nil {
		return false, err
	} else if empty {
		return false, nil
	}
	// ---

	// read cache indicator file (Carthage/Cachefile)
	cacheFilePth := filepath.Join(carthageDir, cacheFileName)
	cacheFileContent, err := contentOfFile(cacheFilePth)
	if err != nil {
		return false, err
	}
	log.Debugf(cacheFileName + " exists in " + cacheFilePth + "\n")

	if cacheFileContent == "" {
		return false, nil
	}
	log.Debugf(cacheFileName + " content: " + cacheFileContent + "\n")

	// --

	// read Cartfile.resolved
	resolvedFilePath := filepath.Join(srcDir, resolvedFileName)
	log.Debugf(resolvedFileName + " exists in " + resolvedFilePath + "\n")

	resolvedFileContent, err := contentOfFile(resolvedFilePath)
	if err != nil {
		return false, err

	}
	log.Debugf(resolvedFileName + " content: " + resolvedFileContent + "\n")

	if resolvedFileContent == "" {
		return false, nil
	}

	// ---

	desiredCacheContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s", swiftVersion, resolvedFileName, resolvedFileContent, resolvedFileName)
	return cacheFileContent == desiredCacheContent, nil
}

func collectCarthageCache(projectDir string) error {
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

func main() {
	var configs Config
	if err := stepconf.Parse(&configs); err != nil {
		fail("Could not create config: %s", err)
	}
	stepconf.Print(configs)

	log.SetEnableDebugLog(true)

	// Environment
	fmt.Println()
	log.Infof("Environment:")

	carthageVersion, err := getCarthageVersion()
	if err != nil {
		fail("Failed to get carthage version, error: %s", err)
	}
	log.Printf("- CarthageVersion: %s", carthageVersion.String())

	swiftVersion, err := getSwiftVersion()
	if err != nil {
		fail("Failed to get swift version, error: %s", err)
	}
	log.Printf("- SwiftVersion: %s", strings.Replace(swiftVersion, "\n", " - ", -1))
	// --

	// Parse options
	customOptions := []string{}
	if configs.CarthageOptions != "" {
		options, err := shellquote.Split(configs.CarthageOptions)
		if err != nil {
			fail("Failed to shell split CarthageOptions (%s), error: %s", configs.CarthageOptions)
		}
		customOptions = options
	}

	projectDir := configs.SourceDir

	isNextOptionProjectDir := false
	for _, option := range customOptions {
		if option == "--project-directory" {
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
	if configs.CarthageCommand == bootstrapCommand {
		fmt.Println()
		log.Infof("Check if cache is available")

		cacheAvailable, err := isCacheAvailable(projectDir, swiftVersion)
		if err != nil {
			fail("Failed to check if cached is available, error: %s", err)
		}

		log.Printf("cache available: %v", cacheAvailable)

		if cacheAvailable {
			if err := collectCarthageCache(projectDir); err != nil {
				log.Warnf("Cache collection skipped: %s", err)
			}
			log.Donef("Using cached dependencies for bootstrap command. If you would like to force update your dependencies, select `update` as CarthageCommand and re-run your build.")
			os.Exit(0)
		}
	}
	// ---

	// Run carthage command
	fmt.Println()
	log.Infof("Running Carthage command")

	args := append([]string{configs.CarthageCommand}, customOptions...)
	cmd := command.New("carthage", args...)

	if configs.GithubAccessToken != "" {
		log.Printf("Appending GITHUB_ACCESS_TOKEN to process environments")

		cmd.AppendEnvs(fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", string(configs.GithubAccessToken)))
	}

	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	log.Donef("$ %s", cmd.PrintableCommandArgs())
	fmt.Println()

	if err := cmd.Run(); err != nil {
		fail("Carthage command failed, error: %s", err)
	}
	// ---

	// Create cache
	if configs.CarthageCommand == bootstrapCommand {
		fmt.Println()
		log.Infof("Creating cache")

		cacheFilePth := filepath.Join(projectDir, carthageDirName, cacheFileName)

		resolvedFilePath := filepath.Join(projectDir, resolvedFileName)
		resolvedFileContent, err := contentOfFile(resolvedFilePath)
		if err != nil {
			fail("Failed to read Cartfile.resolved, error: %s", err)
		}
		if resolvedFileContent == "" {
			log.Warnf("Cartfile.resolved is empty or not exists at: %s", resolvedFilePath)
			log.Warnf("Skipping caching")
			os.Exit(1)
		}

		cacheContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s", swiftVersion, resolvedFileName, resolvedFileContent, resolvedFileName)

		carthageDir := filepath.Join(projectDir, carthageDirName)
		if exist, err := pathutil.IsPathExists(carthageDir); err != nil {
			fail("Failed to check if dir exists at (%s), error: %s", carthageDir, err)
		} else if !exist {
			if err := os.Mkdir(carthageDir, 0777); err != nil {
				fail("Failed to create dir (%s), error: %s", carthageDir, err)
			}
		}

		if err := fileutil.WriteStringToFile(cacheFilePth, cacheContent); err != nil {
			fail("Failed to write cahe file, error: %s", err)
		}

		log.Donef("Cachefile created: %s", cacheFilePth)

		if err := collectCarthageCache(projectDir); err != nil {
			log.Warnf("Cache collection skipped: %s", err)
		}
	}
	// ---
}
