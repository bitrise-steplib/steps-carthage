package cachedcarthage

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/retry"
)

const (
	bootstrapCommand = "bootstrap"
	updateCommand    = "update"
)

// CarthageCache ...
type CarthageCache interface {
	Commit() error
	CreateIndicator() error
	IsAvailable() (bool, error)
}

// CommandBuilder ...
type CommandBuilder interface {
	AddGitHubToken(githubToken stepconf.Secret) CommandBuilder
	AddXCConfigFile(path string) CommandBuilder
	Append(args ...string) CommandBuilder
	Command() *command.Model
}

// Runner can be used to execute Carthage command and cache the results.
type Runner struct {
	carthageCommand   string
	args              []string
	githubAccessToken stepconf.Secret
	xcconfigPath      string
	cache             CarthageCache
	commandBuilder    CommandBuilder
}

// NewRunner ...
func NewRunner(
	carthageCommand string,
	args []string,
	githubAccessToken stepconf.Secret,
	xcconfigPath string,
	cache CarthageCache,
	commandBuilder CommandBuilder,
) Runner {
	return Runner{
		carthageCommand:   carthageCommand,
		args:              args,
		githubAccessToken: githubAccessToken,
		xcconfigPath:      xcconfigPath,
		cache:             cache,
		commandBuilder:    commandBuilder,
	}
}

// Run ...
func (runner Runner) Run() error {

	if runner.carthageCommand == bootstrapCommand {
		if runner.isCacheAvailable() {
			log.Donef("Cache available")

			log.Infof("Committing Cachefile...")
			err := runner.cache.Commit()
			if err == nil {
				log.Donef("Using cached dependencies for bootstrap command. If you would like to force update your dependencies, select `update` as CarthageCommand and re-run your build.")
				return nil
			}

			log.Warnf("Cache collection skipped: %s", err)
		} else {
			log.Warnf("Cache not available")
		}
	}

	if err := runner.perform(); err != nil {
		if runnerErr, ok := err.(*RunnerError); ok {
			runnerErr.Err = fmt.Errorf("Carthage command failed, error: %s", runnerErr.Err)
		}

		return err
	}

	if runner.carthageCommand == bootstrapCommand {
		log.Infof("Creating cache indicator")
		if err := runner.cache.CreateIndicator(); err != nil {
			return err
		}

		if err := runner.cache.Commit(); err != nil {
			log.Warnf("Cache committing skipped: %s", err)
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

func (runner Runner) perform() error {
	var function = runner.executeCommand

	if contains(getRetryableCommands(), runner.carthageCommand) {
		function = func() error {
			return retry.Times(1).Wait(3 * time.Second).TryWithAbort(func(attempt uint) (error, bool) {
				if attempt > 0 {
					log.Warnf("Carthage %s (possible) network failure, retrying ...", runner.carthageCommand)
				}

				err := runner.executeCommand()

				return err, !hasRetryableFailure(err)
			})
		}
	}

	return function()
}

func (runner Runner) executeCommand() error {
	log.Infof("Running Carthage command")

	builder := runner.commandBuilder.
		AddGitHubToken(runner.githubAccessToken).
		AddXCConfigFile(runner.xcconfigPath).
		Append(runner.carthageCommand).
		Append(runner.args...)
	var stderrBuf bytes.Buffer

	cmd := builder.Command()
	cmd.SetStdout(os.Stdout)
	cmd.SetStderr(io.MultiWriter(os.Stderr, &stderrBuf))

	log.Donef("$ %s", cmd.PrintableCommandArgs())

	err := cmd.Run()

	if err == nil {
		return nil
	}

	return &RunnerError{stderrBuf.String(), err}
}

func contains(slice []string, value string) bool {
	for _, item := range slice {
		if value == item {
			return true
		}
	}

	return false
}
