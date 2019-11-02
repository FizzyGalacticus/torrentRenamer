// +build darwin linux

package exec

// ExecuteSudoCommand - Executes the given command with attempted elevated priveleges
func ExecuteSudoCommand(cmdStr string, args ...string) (string, error) {
	return ExecuteCommand("sudo", append([]string{cmdStr}, args...)...)
}
