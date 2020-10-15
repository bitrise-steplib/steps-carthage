package cachedcarthage

import (
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
)

const (
	bootstrapCommand = "bootstrap"
)

// CarthageCache ...
type CarthageCache interface {
	Collect() error
	Create() error
	IsAvailable() (bool, error)
}

// CommandBuilder ...
type CommandBuilder interface {
	AddGitHubToken(githubToken stepconf.Secret) CommandBuilder
	Append(args ...string) CommandBuilder
	Command() *command.Model
}

// Runner ...
type Runner struct {
	command           string
	args              []string
	githubAccessToken stepconf.Secret
	cache             CarthageCache
	commandBuilder    CommandBuilder
}

// NewRunner ...
func NewRunner(
	command string,
	args []string,
	githubAccessToken stepconf.Secret,
	cache CarthageCache,
	commandBuilder CommandBuilder,
) Runner {
	return Runner{
		command:           command,
		args:              args,
		githubAccessToken: githubAccessToken,
		cache:             cache,
		commandBuilder:    commandBuilder,
	}
}

// Run ...
func (runner Runner) Run() error {

	if runner.command == bootstrapCommand {
		if runner.isCacheAvailable() {
			log.Successf("Cache available")

			log.Infof("Collecting carthage caches...")
			err := runner.cache.Collect()
			if err == nil {
				log.Donef("Using cached dependencies for bootstrap command. If you would like to force update your dependencies, select `update` as CarthageCommand and re-run your build.")
				return nil
			}

			log.Warnf("Cache collection skipped: %s", err)
		} else {
			log.Errorf("Cache not available")
		}
	}

	if err := runner.executeCommand(); err != nil {
		return fmt.Errorf("Carthage command failed, error: %s", err)
	}

	if runner.command == bootstrapCommand {
		log.Infof("Creating cache")
		if err := runner.cache.Create(); err != nil {
			return err
		}

		if err := runner.cache.Collect(); err != nil {
			log.Warnf("Cache collection skipped: %s", err)
		}
	}

	return nil
}

func (runner Runner) isCacheAvailable() bool {
	log.Infof("Check if cache is available")

	cacheAvailable, err := runner.cache.IsAvailable()
	if err != nil {
		log.Warnf("Failed to check if cached is available, error: %s", err)
	}

	return cacheAvailable
}

func (runner Runner) executeCommand() error {
	log.Infof("Running Carthage command")

	// log.Printf("Appending GITHUB_ACCESS_TOKEN to process environments")
	builder := runner.commandBuilder.
		AddGitHubToken(runner.githubAccessToken).
		Append(runner.args...)
	cmd := builder.Command()

	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(os.Stderr)

	log.Donef("$ %s", cmd.PrintableCommandArgs())

	return cmd.Run()
}
