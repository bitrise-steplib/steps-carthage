package main

// ConfigsModel ...
import (
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-steplib/steps-carthage/cachedcarthage"
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

	// Parse options
	args := parseCarthageOptions(configs)

	runner := cachedcarthage.NewRunner(
		configs.CarthageCommand,
		args,
		configs.SourceDir,
		configs.GithubAccessToken,
	)
	if err := runner.Run(); err != nil {
		// TODO fail()
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
