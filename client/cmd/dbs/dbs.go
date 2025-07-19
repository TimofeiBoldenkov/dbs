package main

import (
	"time"

	"github.com/TimofeiBoldenkov/dbs/client/providers/RAM_info_provider"
	"github.com/TimofeiBoldenkov/dbs/lib/providers_manager"
)

func main() {
	manager := providersmanager.ProvidersManager{}
	manager.Add(raminfoprovider.RAMInfoProvider{}, time.Second * 3)
	manager.Add(raminfoprovider.RAMInfoProvider{}, time.Second * 3)
	manager.Run()
}
