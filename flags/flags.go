package flags

import (
	"os"

	"github.com/jimbersoftware/pra_client/environment"
)

type Config struct {
	NoLauncher bool
}

func ShowVersionInfo() {
	println("-------------------------------------------------------------")
	println("Version: " + environment.GetVersion())
	println("Build date: " + environment.GetBuildDate())
	println("Environment: " + environment.GetEnvironment())
	println("SHA commit hash: " + environment.GetShaCommitHash())
	println("-------------------------------------------------------------")
	println("Signal server endpoint: " + environment.GetSignalServer())
	println("Websocket endpoint: " + environment.GetSignalServerWs())
	println("-------------------------------------------------------------")
	println("Database base path: " + os.Getenv("DATABASE_BASE_PATH"))
}
