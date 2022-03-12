package shell

import (
	"os"
	"os/exec"
	"path"
)

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

// MkdirP run mkdir -p for a given path
func MkdirP(dirPath string) (err error) {
	args := []string{"-p", dirPath}
	_, err = RunAndWait("mkdir", args)
	return err
}

// DefaultConfigFilePath return the absolute path to the default config file
// $HOME/.stayback.json
func DefaultConfigFilePath() (configPath string, err error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return path.Join(home, ".stayback.json"), err
}
