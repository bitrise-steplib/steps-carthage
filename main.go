package main

// ConfigsModel ...
import (
	"errors"
	"fmt"
	"os"

	"path/filepath"

	"github.com/bitrise-io/go-utils/cmdex"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/kballard/go-shellquote"
)

const (
	carthageDirName  = "Carthage"
	buildDirName     = "Build"
	cacheFileName    = "Cachefile"
	resolvedFileName = "Cartfile.resolved"
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
	log.Info("Configs:")

	log.Detail("- CarthageCommand: %s", configs.CarthageCommand)
	log.Detail("- CarthageOptions: %s", configs.CarthageOptions)
	log.Detail("- GithubAccessToken: %s", configs.GithubAccessToken)

	fmt.Println()
}

func fail(format string, v ...interface{}) {
	log.Error(format, v...)
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
	cmd := cmdex.NewCommand("swift", "-version")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil {
		return "", err
	}
	return out, nil
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

	// read cahce
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

	//
	// Exit if bootstrap is cached
	log.Info("Check if cache is available")

	hasCachedItems, err := isCacheAvailable(configs.SourceDir)
	if err != nil {
		fail("Failed to check cached files, error: %s", err)
	}

	log.Detail("has cahched items: %v", hasCachedItems)

	if configs.CarthageCommand == "bootstrap" && hasCachedItems {
		log.Done("Using cached dependencies for bootstrap command. If you would like to force update your dependencies, select `update` as CarthageCommand and re-run your build.")
		os.Exit(0)
	}

	fmt.Println()
	// ---

	//
	// Run carthage command
	log.Info("Running Carthage command")

	args := append([]string{configs.CarthageCommand}, customOptions...)
	cmd := cmdex.NewCommand("carthage", args...)

	if configs.GithubAccessToken != "" {
		log.Detail("Appending GITHUB_ACCESS_TOKEN to process environments")

		cmd.AppendEnvs([]string{fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", configs.GithubAccessToken)})
	}

	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	log.Done("$ %s", cmdex.PrintableCommandArgs(false, cmd.GetCmd().Args))
	fmt.Println()

	if err := cmd.Run(); err != nil {
		fail("Carthage command failed, error: %s", err)
	}
	// ---

	//
	// Create cache
	if configs.CarthageCommand == "bootstrap" {
		fmt.Println()
		log.Info("Creating cache")

		cacheFilePth := filepath.Join(configs.SourceDir, carthageDirName, cacheFileName)

		swiftVersion, err := swiftVersion()
		if err != nil {
			fail("Failed to get swift version, error: %s", err)
		}

		resolvedFilePath := filepath.Join(configs.SourceDir, resolvedFileName)
		resolved, err := contentsOfCartfileResolved(resolvedFilePath)
		if err != nil {
			fail("Failed to get resolved file content, error: %s", err)
		}

		cacheContent := fmt.Sprintf("--Swift version: %s --Swift version \n --%s: %s --%s", swiftVersion, resolvedFileName, resolved, resolvedFileName)

		carthageDir := filepath.Join(configs.SourceDir, carthageDirName)
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

		log.Done("Cachefile: %s", cacheFilePth)
	}
	// ---
}
