package carthage

import (
	"fmt"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/bitrise-io/go-utils/command"
)

// Carthage ...
type Carthage struct {
	githubToken stepconf.Secret
}

// New ...
func New() Carthage {
	return Carthage{
		githubToken: "",
	}
}

// AddGitHubToken ...
func (carthage Carthage) AddGitHubToken(githubToken stepconf.Secret) Carthage {
	carthage.githubToken = githubToken
	return carthage
}

// Command ...
func (carthage Carthage) Command(args ...string) *command.Model {
	cmd := command.New("carthage", args...)
	if carthage.githubToken != "" {
		cmd.AppendEnvs(fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", string(carthage.githubToken)))
	}
	return cmd
}

// Version ...
func (carthage Carthage) Version() *command.Model {
	cmd := carthage.Command("version")
	return cmd
}
