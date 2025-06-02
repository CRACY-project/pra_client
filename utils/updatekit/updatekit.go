package updatekit

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/jimbersoftware/pra_client/logging"
	"github.com/jimbersoftware/pra_client/utils"
	"github.com/jimbersoftware/pra_client/utils/filekit"
)

const TIMEOUT = 2 * time.Second

type UpdateKit struct {
	signalServerUrl string
	platform        string
	companyName     string
	myTag           string
	appName         string
}

func MakeUpdateKit(signalServerUrl, companyName, platform, myTag, appName string) *UpdateKit {
	return &UpdateKit{
		signalServerUrl: signalServerUrl,
		platform:        platform,
		companyName:     strings.ToLower(companyName),
		myTag:           myTag,
		appName:         appName,
	}
}

func (u *UpdateKit) Start() bool {
	tag := u.getCurrentTag()
	logging.Log(logging.DEBUG, "[UPDATEKIT] Checking TAG ", tag, " vs ", u.myTag+"_"+u.platform)
	if !u.IsUpdateNeeded() {
		return false
	}

	fileManager := filekit.NewFileManager(u.signalServerUrl)
	newBinary := filekit.File{
		SourceName: tag + "_" + u.platform,
		SourcePath: "/binaries/",
		DestName:   "jimberfw-next",
		DestPath:   utils.GetCurrentDir(),
	}

	err := fileManager.DownloadAndVerifyFile(newBinary)
	if err != nil {
		logging.Log(logging.ERROR, "[UPDATEKIT] Error downloading the new version, can't update! ", tag, err.Error())
		return false
	}

	return true
}

func (u *UpdateKit) IsUpdateNeeded() bool {
	tag := u.getCurrentTag()
	if tag == u.myTag || tag == "" {
		return false
	}
	return true
}

func (u *UpdateKit) getCurrentTag() string {
	client := utils.GetHttpClient(TIMEOUT)

	response, err := client.Get(u.signalServerUrl + "/binaries/updates.json")
	if err != nil {
		fmt.Print(err.Error())
		return ""
	}

	responseData, err := ioutil.ReadAll(response.Body)
	var raw interface{}
	json.Unmarshal(responseData, &raw)
	if raw == nil {
		logging.Log(logging.ERROR, "[UPDATEKIT] Can't decode JSON for ", u.platform, u.appName, u.companyName)
		return ""
	}
	allData := raw.(map[string]interface{})

	if allData[u.platform] == nil {
		logging.Log(logging.ERROR, "[UPDATEKIT] Can't find platform ", u.platform, " in JSON")
		return ""
	}
	platformData := allData[u.platform].(map[string]interface{})

	if platformData[u.appName] == nil {
		logging.Log(logging.ERROR, "[UPDATEKIT] Can't find app ", u.appName, " in JSON")
		return ""
	}

	platformDataForApp := platformData[u.appName].(map[string]interface{})

	if platformDataForApp[u.companyName] != nil {
		logging.Log(logging.DEBUG, "Found tag for ", u.companyName, " in JSON", platformDataForApp[u.companyName].(string))
		return platformDataForApp[u.companyName].(string)
	}
	if platformDataForApp["all"] != nil {
		logging.Log(logging.DEBUG, "Found tag for all in JSON", platformDataForApp["all"].(string))
		return platformDataForApp["all"].(string)
	}
	return ""
}
