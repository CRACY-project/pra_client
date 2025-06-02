package updater

import (
	"github.com/jimbersoftware/pra_client/launcher/launcher/types"
	"github.com/jimbersoftware/pra_client/logging"
	"github.com/jimbersoftware/pra_client/utils/updatekit"
)

func DownloadClientUpdate(updateMessage types.UpdateMessage) bool {
	if updateMessage.CompanyName == "" || updateMessage.OsUsername == "" || updateMessage.SignalServer == "" || updateMessage.Tag == "" || updateMessage.Platform == "" {
		logging.Log(logging.ERROR, "Update info not initialized, can't perform update")
		return false
	}
	updateKit := updatekit.MakeUpdateKit(updateMessage.SignalServer, updateMessage.CompanyName, updateMessage.Platform, updateMessage.Tag, "client")
	if updateKit.IsUpdateNeeded() {
		logging.Log(logging.INFO, "Update available, downloading")
		return updateKit.Start()
	}
	return false
}
