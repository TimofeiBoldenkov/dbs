package providersmanager

import (
	"dbs/lib/info_provider"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type ProviderInfo struct {
	Provider infoprovider.InfoProvider
	SleepBetweenRuns time.Duration
	SleepBeforeFirstRun time.Duration
}

type ProvidersManager struct {
	providerInfos []ProviderInfo
}

func (pm *ProvidersManager) Add(Provider infoprovider.InfoProvider, SleepBetweenRuns time.Duration) {
	sleepBeforeFirstRun := time.Duration(rand.Float64() * float64(SleepBetweenRuns))

	pm.AddSetDelay(Provider, SleepBetweenRuns, sleepBeforeFirstRun)
}

func (pm *ProvidersManager) AddSetDelay(
	Provider infoprovider.InfoProvider, 
	SleepBetweenRuns time.Duration, 
	SleepBeforeFirstRun time.Duration) {
	
	pm.providerInfos = append(
		pm.providerInfos,
		ProviderInfo{Provider, SleepBetweenRuns, SleepBeforeFirstRun})
}

func (pm *ProvidersManager) Run() {
	var wg sync.WaitGroup
	for _, info := range pm.providerInfos {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(info.SleepBeforeFirstRun)
			for {
				fmt.Println(info.Provider.GetInfo())
				time.Sleep(info.SleepBetweenRuns)
			}
		}()
	}

	wg.Wait()
}
