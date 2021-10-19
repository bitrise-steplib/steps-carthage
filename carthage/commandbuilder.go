package carthage

import (
	"fmt"
	"github.com/bitrise-io/go-utils/env"
	"io"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-steplib/steps-carthage/cachedcarthage"
)

// CLIBuilder can be used to build cli Carthage commands.
type CLIBuilder struct {
	args []string
	envs []string
	commandFactory command.Factory
}

// NewCLIBuilder ...
func NewCLIBuilder() CLIBuilder {
	return CLIBuilder{
		args: []string{},
		envs: []string{},
		commandFactory: command.NewFactory(env.NewRepository()),
	}
}

// AddGitHubToken appends the provided GitHub token to the builder.
func (builder CLIBuilder) AddGitHubToken(githubToken stepconf.Secret) cachedcarthage.CommandBuilder {
	if githubToken != "" {
		builder.envs = append(builder.envs, fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", string(githubToken)))
	}
	return builder
}

// AddXCConfigFile appends the provided .xcconfig file path to the builder.
func (builder CLIBuilder) AddXCConfigFile(path string) cachedcarthage.CommandBuilder {
	if path != "" {
		builder.envs = append(builder.envs, fmt.Sprintf("XCODE_XCCONFIG_FILE=%s", path))
	}
	return builder
}

// Append adds the arguments to the builder.
func (builder CLIBuilder) Append(args ...string) cachedcarthage.CommandBuilder {
	builder.args = append(builder.args, args...)
	return builder
}

// Command returns the built command.
func (builder CLIBuilder) Command(stdout io.Writer, stderr io.Writer) command.Command {
	command := builder.commandFactory.Create("carthage", builder.args, &command.Opts{
		Stdout: stdout,
		Stderr: stderr,
		Env:    builder.envs,
	})
	return command
}
