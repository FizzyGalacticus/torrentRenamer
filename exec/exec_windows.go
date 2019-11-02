package exec

// ExecuteSudoCommand - Executes the given command with attempted elevated priveleges
func ExecuteSudoCommand(cmdStr string, args ...string) (string, error) {
	return ExecuteCommand(cmdStr, args...)
}
