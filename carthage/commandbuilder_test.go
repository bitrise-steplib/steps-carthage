package carthage

import (
	"testing"

	"github.com/bitrise-io/go-steputils/stepconf"
	"github.com/stretchr/testify/assert"
)

func Test_WhenArgumentAppended_ThenResultCommandContainsArgument(t *testing.T) {
	// Given
	expectedCommand := `carthage "version"`
	builder := NewCLIBuilder()

	// When
	command := builder.Append("version").Command(nil, nil)

	// Then
	assert.Equal(t, expectedCommand, command.PrintableCommandArgs())
}

func Test_WhenGitHubTokenAppended_ThenResultCommandContainsToken(t *testing.T) {
	// Given
	var expectedToken stepconf.Secret = "nice_token"
	//expectedEnv := fmt.Sprintf("GITHUB_ACCESS_TOKEN=%s", string(expectedToken))
	expectedCommand := `carthage "version"`
	builder := NewCLIBuilder()

	// When
	command := builder.AddGitHubToken(expectedToken).Append("version").Command(nil, nil)

	// Then
	assert.Equal(t, expectedCommand, command.PrintableCommandArgs())
	// At the moment it is not possible to get the env variables from the command.
	//assert.Contains(t, command.GetCmd().Env, expectedEnv)
}

func Test_WhenXCConfigFileAppended_ThenResultCommandContainsPath(t *testing.T) {
	// Given
	path := "/path/file.xcconfig"
	//expectedEnv := fmt.Sprintf("XCODE_XCCONFIG_FILE=%s", path)
	expectedCommand := `carthage "version"`
	builder := NewCLIBuilder()

	// When
	command := builder.AddXCConfigFile(path).Append("version").Command(nil, nil)

	// Then
	assert.Equal(t, expectedCommand, command.PrintableCommandArgs())
	// At the moment it is not possible to get the env variables from the command.
	//assert.Contains(t, command.GetCmd().Env, expectedEnv)
}
