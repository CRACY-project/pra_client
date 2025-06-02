package launcher

import (
	"os/exec"

	"github.com/jimbersoftware/pra_client/logging"
	"github.com/jimbersoftware/pra_client/utils"
)

var pendingUpdatePath = utils.GetCurrentDir() + `\jimberpra-next`
var executablePath = utils.GetCurrentDir() + `\jimberpra.exe`

func makeExecutable() {
}

func startService() int {

	logging.Log(logging.INFO, "Starting service with adjust", executablePath)

	g, err := NewProcessExitGroup()
	if err != nil {
		logging.Log(logging.ERROR, "Can't create process exit group", executablePath)
		panic(err)
	}
	defer g.Dispose()
	logging.Log(logging.INFO, "Starting with process group for shutdown", executablePath)

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
		logging.Log(logging.ERROR, "Can't start service", executablePath, err)
		panic(err)
	}

	if err := g.AddProcess(_cmd.Process); err != nil {
		logging.Log(logging.ERROR, "Can't add process to process exit group", executablePath, err)
		panic(err)
	}

	logging.LogReader(stderr, logging.ERROR, "Service crashed with error: ")

	if err := _cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			logging.Log(logging.ERROR, "!Service crashed with error", err, exiterr)
		} else {
			logging.Log(logging.ERROR, "!Service stopped", err)
		}
		if !_serviceIsUpdating {
			startService()
		}
	}
	logging.Log(logging.INFO, "Service stopped, no longer waiting")

	return 0
}

// killRemainingProcesses forcefully terminates any running instances of "JimberNetworkIsolation.exe"
func killRemainingProcesses() {
	killExecutable := exec.Command("taskkill", "/IM", "JimberNetworkIsolation.exe", "/F")
	if err := killExecutable.Run(); err != nil {
		logging.Log(logging.ERROR, "Failed to execute `task kill` command: ", err)
	}

	logging.Log(logging.INFO, "Forcefully terminated 'JimberNetworkIsolation.exe' successfully")
}
