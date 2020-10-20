package carthage

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-steplib/steps-carthage/cachedcarthage"
)

// CLIBuilder can be used to build cli Carthage commands.
type CLIBuilder struct {
	cmd *command.Model
}

// NewCLIBuilder ...
func NewCLIBuilder() CLIBuilder {
	return CLIBuilder{
		cmd: command.New("carthage"),
	}
}

// AddGitHubToken appends the provided GitHub token to the builder.
func (builder CLIBuilder) AddGitHubToken(githubToken stepconf.Secret) cachedcarthage.CommandBuilder {
	if githubToken != "" {
		builder.cmd.AppendEnvs(fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", string(githubToken)))
	}
	return builder
}

// AddXCConfigFile appends the provided .xcconfig file path to the builder.
func (builder CLIBuilder) AddXCConfigFile(path string) cachedcarthage.CommandBuilder {
	if path != "" {
		builder.cmd.AppendEnvs(fmt.Sprintf("XCODE_XCCONFIG_FILE=%s", path))
	}
	return builder
}

// Append adds the arguments to the builder.
func (builder CLIBuilder) Append(args ...string) cachedcarthage.CommandBuilder {
	builder.cmd.GetCmd().Args = append(builder.cmd.GetCmd().Args, args...)
	return builder
}

// Command returns the built command.
func (builder CLIBuilder) Command() *command.Model {
	return builder.cmd
}
