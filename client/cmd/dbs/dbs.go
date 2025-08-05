package main

import (
	"time"

	"github.com/TimofeiBoldenkov/dbs/client/providers/RAM_info_provider"
	processesinfoprovider "github.com/TimofeiBoldenkov/dbs/client/providers/processes_info_provider"
	"github.com/TimofeiBoldenkov/dbs/client/providers_manager"
)

func main() {
	manager := providersmanager.ProvidersManager{}
	manager.Add(raminfoprovider.RAMInfoProvider{}, "RAM", time.Second * 60)
	manager.Add(processesinfoprovider.ProcessesInfoProvider{}, "processes", time.Second * 60)
	manager.Run()
}
