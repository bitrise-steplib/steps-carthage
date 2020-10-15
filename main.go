package main

// ConfigsModel ...
import (
	"fmt"
	"os"
	"strings"

	cacheutil "github.com/bitrise-io/go-steputils/cache"
	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-carthage/cachedcarthage"
	"github.com/bitrise-steplib/steps-carthage/carthage"
	"github.com/hashicorp/go-version"
	"github.com/kballard/go-shellquote"
)

const (
	projectDirArg = "--project-directory"
)

// Config ...
type Config struct {
	GithubAccessToken stepconf.Secret `env:"github_access_token"`
	CarthageCommand   string          `env:"carthage_command,required"`
	CarthageOptions   string          `env:"carthage_options"`
	SourceDir         string          `env:"BITRISE_SOURCE_DIR"`

	// Debug
	VerboseLog bool `env:"verbose_log,opt[yes,no]"`
}

func fail(format string, v ...interface{}) {
	log.Errorf(format, v...)
	os.Exit(1)
}

func main() {
	var configs Config
	if err := stepconf.Parse(&configs); err != nil {
		fail("Could not create config: %s", err)
	}
	stepconf.Print(configs)

	log.SetEnableDebugLog(configs.VerboseLog)

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
	log.Printf("- SwiftVersion: %s", strings.Replace(swiftVersion, "\n", "- ", -1))
	// --

	// Parse options
	args := parseCarthageOptions(configs)

	projectDir := parseProjectDir(configs.SourceDir, args)
	project := cachedcarthage.NewProject(projectDir)
	filecache := cacheutil.New()
	stateProvider := cachedcarthage.DefaultStateProvider{}

	runner := cachedcarthage.NewRunner(
		configs.CarthageCommand,
		args,
		configs.GithubAccessToken,
		cachedcarthage.NewCache(project, swiftVersion, &filecache, stateProvider),
		carthage.NewCLIBuilder(),
	)
	if err := runner.Run(); err != nil {
		fail("Failed to execute step: %s", err)
	}
}

func parseCarthageOptions(config Config) []string {
	var customCarthageOptions []string
	if config.CarthageOptions != "" {
		options, err := shellquote.Split(config.CarthageOptions)
		if err != nil {
			fail("Failed to shell split CarthageOptions (%s), error: %s", config.CarthageOptions)
		}
		customCarthageOptions = options
	}
	return customCarthageOptions
}

func getCarthageVersion() (*version.Version, error) {
	cmd := carthage.NewCLIBuilder().Append("version").Command()
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

func getSwiftVersion() (string, error) {
	cmd := command.New("swift", "-version")
	return cmd.RunAndReturnTrimmedCombinedOutput()
}

func parseProjectDir(originalDir string, customCarthageOptions []string) string {
	projectDir := originalDir

	isNextOptionProjectDir := false
	for _, option := range customCarthageOptions {
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

	return projectDir
}
