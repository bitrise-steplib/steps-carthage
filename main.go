package main

// ConfigsModel ...
import (
	"errors"
	"fmt"
	"os"

	"path/filepath"

	"strings"

	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	version "github.com/hashicorp/go-version"
	"github.com/kballard/go-shellquote"
)

const (
	carthageDirName               = "Carthage"
	buildDirName                  = "Build"
	cacheFileName                 = "Cachefile"
	resolvedFileName              = "Cartfile.resolved"
	buildCacheSupportSinceVersion = "0.20.0"
)

// ConfigsModel ...
type ConfigsModel struct {
	GithubAccessToken string
	CarthageCommand   string
	CarthageOptions   string
	SourceDir         string
}

func createConfigsModelFromEnvs() ConfigsModel {
	return ConfigsModel{
		CarthageCommand:   os.Getenv("carthage_command"),
		CarthageOptions:   os.Getenv("carthage_options"),
		GithubAccessToken: os.Getenv("github_access_token"),
		SourceDir:         os.Getenv("BITRISE_SOURCE_DIR"),
	}
}

func (configs ConfigsModel) print() {
	log.Infof("Configs:")

	log.Printf("- CarthageCommand: %s", configs.CarthageCommand)
	log.Printf("- CarthageOptions: %s", configs.CarthageOptions)
	log.Printf("- GithubAccessToken: %s", configs.GithubAccessToken)

	fmt.Println()
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func (configs ConfigsModel) validate() error {
	if configs.CarthageCommand == "" {
		return errors.New("no CarthageCommand parameter specified")
	}

	return nil
}

func contentsOfCartfileResolved(pth string) (string, error) {
	content, err := fileutil.ReadStringFromFile(pth)
	if err != nil {
		return "", err
	}
	return content, nil
}

func swiftVersion() (string, error) {
	cmd := command.New("swift", "-version")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", err
	}
	return out, nil
}

func indexInStringSlice(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func isCarthageBuildCacheSupported() (bool, *version.Version, error) {
	// get carthage version cmd
	cmd := command.New("carthage", "version")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return false, nil, err
	}

	// get Version which supports build cache
	buildCacheSupportVersion, err := version.NewVersion(buildCacheSupportSinceVersion)
	if err != nil {
		return false, nil, err
	}

	// if output is multi-line, get the last line of string
	// parse Version from cmd output
	for _, outLine := range strings.Split(out, "\n") {
		if currentVersion, err := version.NewVersion(outLine); err == nil {
			return !currentVersion.LessThan(buildCacheSupportVersion), currentVersion, nil
		}
	}

	return false, nil, errors.New("failed to parse `$ carthage version` output")
}

func isCacheAvailable(srcDir string) (bool, error) {
	carthageDir := filepath.Join(srcDir, carthageDirName)
	if exist, err := pathutil.IsPathExists(carthageDir); err != nil {
		return false, err
	} else if !exist {
		return false, nil
	}

	buildDir := filepath.Join(carthageDir, buildDirName)
	if exist, err := pathutil.IsPathExists(buildDir); err != nil {
		return false, err
	} else if exist {
		pattern := filepath.Join(buildDir, "*")
		files, err := filepath.Glob(pattern)
		if err != nil {
			return false, err
		}
		if len(files) == 0 {
			return false, nil
		}
	} else {
		return false, nil
	}

	// read cache
	cacheContent := ""

	cacheFilePth := filepath.Join(srcDir, carthageDirName, cacheFileName)
	if exist, err := pathutil.IsPathExists(cacheFilePth); err != nil {
		return false, err
	} else if exist {
		cacheContent, err = fileutil.ReadStringFromFile(cacheFilePth)
		if err != nil {
			return false, err
		}
	} else {
		return false, nil
	}

	swiftVersion, err := swiftVersion()
	if err != nil {
		return false, err
	}

	resolvedFilePath := filepath.Join(srcDir, resolvedFileName)
	resolved, err := contentsOfCartfileResolved(resolvedFilePath)
	if err != nil {
		return false, err
	}

	desiredCacheContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s", swiftVersion, resolvedFileName, resolved, resolvedFileName)

	return cacheContent == desiredCacheContent, nil
}

func main() {
	configs := createConfigsModelFromEnvs()

	fmt.Println()
	configs.print()

	if err := configs.validate(); err != nil {
		fail("Issue with input: %s", err)
	}

	customOptions := []string{}
	if configs.CarthageOptions != "" {
		options, err := shellquote.Split(configs.CarthageOptions)
		if err != nil {
			fail("Failed to shell split CarthageOptions (%s), error: %s", configs.CarthageOptions)
		}
		customOptions = options
	}

	// get build cache support and version
	isCarthageBuildCacheSupported, currentCarthageVersion, err := isCarthageBuildCacheSupported()
	if err != nil {
		fail("Failed to get carthage version, error: %s", err)
	}

	cacheBuildFlagInCustomOptions := indexInStringSlice("--cache-builds", customOptions)

	log.Infof("Carthage version: %s", currentCarthageVersion.String())
	if !cacheBuildFlagInCustomOptions && isCarthageBuildCacheSupported {
		log.Warnf("Built in cache is available, adding --cache-builds flag")
	} else if cacheBuildFlagInCustomOptions {
		if !isCarthageBuildCacheSupported {
			log.Warnf("Invalid flag --cache-builds")
			log.Printf("It's supported since carthage version (%s), your carthage version: %s", buildCacheSupportSinceVersion, currentCarthageVersion.String())
			fmt.Println()
		} else {
			log.Printf("--cache-builds flag found")
		}
	}
	log.Printf("To save cache files use Cache Pull and Cache Push steps")
	fmt.Println()

	projectDir := configs.SourceDir
	isNextOptionProjectDir := false
	for _, option := range customOptions {
		if option == "--project-directory" {
			isNextOptionProjectDir = true
			continue
		}

		if isNextOptionProjectDir {
			projectDir = option

			log.Infof("--project-directory flag found with value: %s", projectDir)
			log.Printf("using %s as working directory", projectDir)

			break
		}
	}

	//
	// Exit if bootstrap is cached
	if !isCarthageBuildCacheSupported {
		fmt.Println()
		log.Infof("Check if cache is available")

		hasCachedItems, err := isCacheAvailable(projectDir)
		if err != nil {
			fail("Failed to check cached files, error: %s", err)
		}

		log.Printf("has cached items: %v", hasCachedItems)

		if configs.CarthageCommand == "bootstrap" && hasCachedItems {
			log.Donef("Using cached dependencies for bootstrap command. If you would like to force update your dependencies, select `update` as CarthageCommand and re-run your build.")
			os.Exit(0)
		}
	}
	// ---

	//
	// Run carthage command
	fmt.Println()
	log.Infof("Running Carthage command")

	args := append([]string{configs.CarthageCommand}, customOptions...)

	if isCarthageBuildCacheSupported && !cacheBuildFlagInCustomOptions && configs.CarthageCommand == "bootstrap" {
		args = append(args, "--cache-builds")
	}

	cmd := command.New("carthage", args...)

	if configs.GithubAccessToken != "" {
		log.Printf("Appending GITHUB_ACCESS_TOKEN to process environments")

		cmd.AppendEnvs(fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", configs.GithubAccessToken))
	}

	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	log.Donef("$ %s", command.PrintableCommandArgs(false, cmd.GetCmd().Args))
	fmt.Println()

	if err := cmd.Run(); err != nil {
		fail("Carthage command failed, error: %s", err)
	}
	// ---

	//
	// Create cache
	if configs.CarthageCommand == "bootstrap" && !isCarthageBuildCacheSupported {
		fmt.Println()
		log.Infof("Creating cache")

		cacheFilePth := filepath.Join(projectDir, carthageDirName, cacheFileName)

		swiftVersion, err := swiftVersion()
		if err != nil {
			fail("Failed to get swift version, error: %s", err)
		}

		resolvedFilePath := filepath.Join(projectDir, resolvedFileName)
		resolved, err := contentsOfCartfileResolved(resolvedFilePath)
		if err != nil {
			fail("Failed to get resolved file content, error: %s", err)
		}

		cacheContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s", swiftVersion, resolvedFileName, resolved, resolvedFileName)

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

		log.Donef("Cachefile: %s", cacheFilePth)
	}
	// ---
}
