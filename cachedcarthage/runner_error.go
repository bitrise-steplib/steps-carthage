package cachedcarthage

import (
	"errors"
	"strings"
)

// RunnerError ...
type RunnerError struct {
	Output string
	Err error
}

// Error ...
func (e *RunnerError) Error() string {
	return e.Err.Error()
}

func getRetryableCommands() []string {
	return []string{bootstrapCommand, updateCommand}
}

func getErrorSlices() []string {
	return []string{"failed to connect to", "timed out"}
}

func hasRetryableFailure(err error) bool {
	var runnerError *RunnerError

	if errors.As(err, &runnerError) {
		output := strings.ToLower(runnerError.Output)

		for _, str := range getErrorSlices(){
			if strings.Contains(output, str) {
				return true
			}
		}
	}

	return false
}
