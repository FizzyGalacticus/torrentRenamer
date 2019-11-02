package exec

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func parseCommandArgInterfaces(args ...interface{}) ([]string, []string, error) {
	executeArgs := make([]string, 0)
	logArgs := make([]string, 0)
	var err error

	for _, arg := range args {
		if val, ok := arg.(string); ok {
			executeArgs = append(executeArgs, val)
			logArgs = append(logArgs, val)
		} else if arr, ok := arg.([]string); ok {
			executeArgs = append(executeArgs, arr[0])
			logArgs = append(logArgs, arr[1])
		} else {
			err = errors.New("Could not determine type of given command.")
		}
	}

	return executeArgs, logArgs, err
}

func Execute(args []string, inReader io.Reader, outWriter io.Writer, errWriter io.Writer) error {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = inReader
	cmd.Stdout = outWriter
	cmd.Stderr = errWriter

	return cmd.Run()
}

// ExecuteCommandWithSTDOutput - Executes the given command, printing output to
// stdout and stderr.
func ExecuteCommandWithSTDOutput(cmdStr string, args ...string) error {
	return Execute(append([]string{cmdStr}, args...), os.Stdin, os.Stdout, os.Stderr)
}

// ExecuteCommand - Executes the given command, inserting output into a string
// and error.
func ExecuteCommand(cmdStr string, args ...string) (string, error) {
	var out strings.Builder
	var stderr strings.Builder

	err := Execute(append([]string{cmdStr}, args...), os.Stdin, &out, &stderr)
	if err != nil {
		err = errors.New(stderr.String())
	}

	return out.String(), err
}

// ExecuteCommandWithLog - Functions just as exec.ExecuteCommand, but also prints
// what command is being executed.
func ExecuteCommandWithLog(cmdStr string, args ...interface{}) (string, error) {
	executeArgs, logArgs, err := parseCommandArgInterfaces(args...)
	if err != nil {
		return "", err
	}

	fmt.Println("Executing command:", cmdStr, strings.Join(logArgs, " "))
	return ExecuteCommand(cmdStr, executeArgs...)
}

// ExecuteCommandWithSTDOutputAndLog - Combines exec.ExecuteCommandWithSTDOutput and
// exec.ExecuteCommandWithLog
func ExecuteCommandWithSTDOutputAndLog(cmdStr string, args ...interface{}) error {
	executeArgs, logArgs, err := parseCommandArgInterfaces(args...)
	if err != nil {
		return err
	}

	fmt.Println("Executing command:", cmdStr, strings.Join(logArgs, " "))
	return ExecuteCommandWithSTDOutput(cmdStr, executeArgs...)
}

// IsCommandInPath - Returns true if command succeeds, false otherwise
func IsCommandInPath(cmd string, args ...string) bool {
	_, err := ExecuteCommand(cmd, args...)

	return err == nil
}
