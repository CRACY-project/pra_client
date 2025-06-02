package launcher

import (
	"encoding/json"
	"os"
	"strconv"

	events "github.com/jimbersoftware/pra_client/events"
	"github.com/jimbersoftware/pra_client/extendedrpc"
	"github.com/jimbersoftware/pra_client/launcher/launcher/types"
	"github.com/jimbersoftware/pra_client/launcher/launcher/updater"

	"github.com/jimbersoftware/pra_client/logging"
)

var customPortEnv = os.Getenv("JIMBER_PORT")
var customPort, _ = strconv.Atoi(customPortEnv)

var _launcherPort = "127.0.0.1:11003"

func GetLauncherPort() string {
	if customPortEnv != "" {
		return "127.0.0.1:" + strconv.Itoa(customPort+2)
	}
	return _launcherPort
}

func StartLauncherRpc() {
	go startListening()
}

func startListening() {
	server := extendedrpc.NewServer()

	server.Register(server)
	server.Register(&extendedrpc.Protocol{})

	logging.Log(logging.INFO, "Launcher RPC listening on ", GetLauncherPort())
	err := server.Listen("tcp", GetLauncherPort())

	if err != nil {
		logging.Log(logging.ERROR, "Can't start service daemon on port ", GetLauncherPort(), "with error: ", err)
	}
	server.Accept()
}

func handleMessage(name string, message string) string {
	logging.Log(logging.INFO, "Received message: ", name, " ", message)
	switch name {
	case events.DOWNLOAD_UPDATE:
		updateMessage := types.UpdateMessage{}
		json.Unmarshal([]byte(message), &updateMessage)
		logging.Log(logging.INFO, "downloadupdate message received, will download")
		return strconv.FormatBool(updater.DownloadClientUpdate(updateMessage))
	case events.UPDATE:
		go events.Get().Emit(events.UPDATE, nil)
		logging.Log(logging.INFO, "Update message received, will update")
		return ""
	case events.CLEANUP_COMPLETED:
		events.Get().Emit(events.CLEANUP_COMPLETED, nil)
	}
	return ""
}
