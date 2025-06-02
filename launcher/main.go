package main

import (
	"fmt"

	"github.com/jimbersoftware/pra_client/environment"
	"github.com/jimbersoftware/pra_client/launcher/launcher"
	"github.com/kardianos/service"
)

const (
	serviceName        = "Jimber PRA client"
	serviceDescription = "Client for Jimber PRA"
)

type program struct{}

func (p program) Start(s service.Service) error {
	fmt.Println(s.String() + " started")
	go p.run()
	return nil
}

func (p program) Stop(s service.Service) error {
	fmt.Println(s.String() + " stopped")
	return nil
}

func (p program) run() {
	desktopMode := false
	// Move old log files

	launcher.StartLauncherApp(desktopMode, environment.GetLauncherLogPath(), environment.GetDebugLogPath(), environment.GetPlatform())
}

func main() {
	// gracefulshutdown.ListenToRPCEvents()

	serviceConfig := &service.Config{
		Name:        serviceName,
		DisplayName: serviceName,
		Description: serviceDescription,
	}
	prg := &program{}
	s, err := service.New(prg, serviceConfig)
	if err != nil {
		fmt.Println("Cannot create the service: " + err.Error())
	}
	err = s.Run()
	if err != nil {
		fmt.Println("Cannot start the service: " + err.Error())
	}
}
