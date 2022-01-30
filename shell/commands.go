package shell

import "os/exec"

type RunAndWaitOutput struct {
	ExitCode       int
	CombinedOutput string
}

// RunAndWait
func RunAndWait(cmdRoot string, cmdOpts []string) (result RunAndWaitOutput, err error) {
	cmd := exec.Command(cmdRoot, cmdOpts...)
	b, err := cmd.CombinedOutput()
	result.CombinedOutput = string(b)
	result.ExitCode = cmd.ProcessState.ExitCode()
	return result, err
}
