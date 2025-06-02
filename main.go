//go:build !niac
// +build !niac

package main

import (
	"time"

	"github.com/jimbersoftware/pra_client/api"
	"github.com/jimbersoftware/pra_client/libraryinfo/ubuntu"
	"github.com/jimbersoftware/pra_client/logging"
	"github.com/jimbersoftware/pra_client/settings"
	"github.com/jimbersoftware/pra_client/sysinfo"
)

var client *api.Client

func main() {

	settings, err := settings.LoadSettings()

	if err != nil {
		logging.Log(logging.ERROR, "Error loading settings:", err)
		return
	}
	logging.Log(logging.INFO, "Starting heartbeat with token:", settings.Token)
	client = api.NewClient(settings.Token)

	SendSystemInfo()
	SendLibraryInfo()
	time.Sleep(5 * time.Minute) // Wait for the client to initialize

}

func SendSystemInfo() {
	logging.Log(logging.INFO, "Gathering system information...")
	systemInformation, err := sysinfo.GatherSystemInfo()
	if err != nil {
		logging.Log(logging.ERROR, "Error gathering system information:", err)

	}
	client.SendSystemInfo(*systemInformation)

}

func SendLibraryInfo() {
	logging.Log(logging.INFO, "Gathering library information...")
	packages, err := ubuntu.GetInstalledPackages()

	if err != nil {
		logging.Log(logging.ERROR, "Error gathering library information:", err)
		return
	}

	client.SendLibraryInfo(packages)
	logging.Log(logging.INFO, "Library information sent successfully")
}
