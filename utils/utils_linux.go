package utils

import (
	"os/exec"

	"github.com/jimbersoftware/pra_client/logging"
)

func RunCommandWithOutput(command string) string {
	logging.Log(logging.DEBUG, "RUNNING COMMAND: ", command)
	cmd, err := exec.Command("/bin/sh", "-c", command).Output()

	logging.Log(logging.DEBUG, "OUTPUT: ", string(cmd))

	if err != nil {
		logging.Log(logging.WARNING, "CMD FAILED", command, err.Error())
	}

	return string(cmd)
}
