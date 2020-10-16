package carthage

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-steplib/steps-carthage/cachedcarthage"
)

// CLIBuilder ...
type CLIBuilder struct {
	cmd *command.Model
}

// NewCLIBuilder ...
func NewCLIBuilder() CLIBuilder {
	return CLIBuilder{
		cmd: command.New("carthage"),
	}
}

// AddGitHubToken ...
func (builder CLIBuilder) AddGitHubToken(githubToken stepconf.Secret) cachedcarthage.CommandBuilder {
	if githubToken != "" {
		builder.cmd.AppendEnvs(fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", string(githubToken)))
	}
	return builder
}

// AddXCConfigFile ...
func (builder CLIBuilder) AddXCConfigFile(path string) cachedcarthage.CommandBuilder {
	if path != "" {
		builder.cmd.AppendEnvs(fmt.Sprintf("XCODE_XCCONFIG_FILE=%s", path))
	}
	return builder
}

// Append ...
func (builder CLIBuilder) Append(args ...string) cachedcarthage.CommandBuilder {
	builder.cmd.GetCmd().Args = append(builder.cmd.GetCmd().Args, args...)
	return builder
}

// Command ...
func (builder CLIBuilder) Command() *command.Model {
	return builder.cmd
}
