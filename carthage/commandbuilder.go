package carthage

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-steplib/steps-carthage/cachedcarthage"
)

// CLIBuilder ...
type CLIBuilder struct {
	githubToken stepconf.Secret
	cmd         *command.Model
}

// NewCLIBuilder ...
func NewCLIBuilder() CLIBuilder {
	return CLIBuilder{
		githubToken: "",
		cmd:         command.New("carthage"),
	}
}

// AddGitHubToken ...
func (builder CLIBuilder) AddGitHubToken(githubToken stepconf.Secret) cachedcarthage.CommandBuilder {
	builder.githubToken = githubToken
	return builder
}

// Append ...
func (builder CLIBuilder) Append(args ...string) cachedcarthage.CommandBuilder {
	builder.cmd.GetCmd().Args = append(builder.cmd.GetCmd().Args, args...)
	return builder
}

// Command ...
func (builder CLIBuilder) Command() *command.Model {
	if builder.githubToken != "" {
		builder.cmd.AppendEnvs(fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", string(builder.githubToken)))
	}
	return builder.cmd
}
