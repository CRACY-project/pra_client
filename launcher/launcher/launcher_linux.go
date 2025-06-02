package launcher

import (
	"os/exec"
	"time"

	"github.com/jimbersoftware/pra_client/logging"
	"github.com/jimbersoftware/pra_client/utils"
)

var pendingUpdatePath = utils.GetCurrentDir() + `/jimberpra-next`
var jimberBinaryName = `/jimberpra`
var executablePath = utils.GetCurrentDir() + jimberBinaryName

func makeExecutable() {
	utils.RunCommandWithOutput("chmod +x " + executablePath)
}

func startService() int {
	logging.Log(logging.INFO, "Starting service", executablePath)

	if _desktopApp {
		_cmd = exec.Command(executablePath, "service")
	} else {
		_cmd = exec.Command(executablePath)
	}

	stderr, err := _cmd.StderrPipe()
	if err != nil {
		logging.Log(logging.ERROR, "Error creating stderr pipe", err)
		return 1
	}
	if err := _cmd.Start(); err != nil {
		logging.Log(logging.ERROR, "Can't start service", executablePath)
		panic(err)
	}

	logging.LogReader(stderr, logging.ERROR, "Service crashed with error: ")

	if err := _cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			logging.Log(logging.ERROR, "!Service crashed with error", exiterr)
		} else {
			logging.Log(logging.ERROR, "!Service stopped", err)
		}

		if !_serviceIsUpdating {
			time.Sleep(1 * time.Second)
			startService()
		}
	}

	logging.Log(logging.INFO, "Service stopped, no longer waiting", _serviceIsUpdating)
	return 0
}

func killRemainingProcesses() {
}
