package launcher

import (
	"os"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/jimbersoftware/pra_client/logging"
)

var _desktopApp = false

var _environment = ""
var logMaxSize = 10

// DISPLAY variable in linux to restart application
var _physicalDisplay = ""

var _cmd *exec.Cmd
var _serviceIsUpdating = false

func StartLauncherApp(desktopApp bool, logFile string, debugLogFile string, environment string) {
	_desktopApp = desktopApp
	_environment = environment

	logging.InitLog(logFile, debugLogFile, logMaxSize)

	StartLauncherRpc()
	// global_events := events.Get()
	// global_events.On(events.UPDATE, func(payload ...interface{}) {
	// 	UpdateClient()
	// 	_serviceIsUpdating = false
	// })

	startService()

	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func UpdateClient() {
	_serviceIsUpdating = true

	// an update is pending
	logging.Log(logging.ERROR, "Update pending, removing previous", executablePath)

	// Attempt to kill the process associated with the global _cmd variable
	if err := _cmd.Process.Kill(); err != nil {
		logging.Log(logging.ERROR, "Failed to kill process: ", err)
	}

	processState, err := _cmd.Process.Wait()

	// forcefully terminate "JimberNetworkIsolation.exe" on Windows
	killRemainingProcesses()

	logging.Log(logging.INFO, "Process state: ", processState, err)

	for i := 1; i < 10; i++ {
		logging.Log(logging.INFO, "Trying replace, "+strconv.Itoa(i))
		time.Sleep(200 * time.Millisecond)
		errRemove := os.Remove(executablePath)
		if errRemove != nil {
			logging.Log(logging.ERROR, "Could not remove JimberNetworkIsolation executable. Can't update:", errRemove)
			continue
		}

		errRename := os.Rename(pendingUpdatePath, executablePath)
		if errRename != nil {
			logging.Log(logging.ERROR, "Could not update JimberNetworkIsolation executable.", errRename)
			continue
		} else {
			break
		}

	}

	makeExecutable()
	_serviceIsUpdating = false
	go startService()
}
