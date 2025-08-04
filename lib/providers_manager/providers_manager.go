package providersmanager

import (
	"bytes"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/TimofeiBoldenkov/dbs/lib/info_provider"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type ProviderInfo struct {
	Provider            infoprovider.InfoProvider
	Name				string
	SleepBetweenRuns    time.Duration
	SleepBeforeFirstRun time.Duration
}

type ProvidersManager struct {
	providerInfos []ProviderInfo
}

func (pm *ProvidersManager) Add(
	Provider infoprovider.InfoProvider, 
	Name string, 
	SleepBetweenRuns time.Duration) {
	sleepBeforeFirstRun := time.Duration(rand.Float64() * float64(SleepBetweenRuns))

	pm.AddSetDelay(Provider, Name, SleepBetweenRuns, sleepBeforeFirstRun)
}

func (pm *ProvidersManager) AddSetDelay(
	Provider infoprovider.InfoProvider,
	Name string,
	SleepBetweenRuns time.Duration,
	SleepBeforeFirstRun time.Duration) {

	pm.providerInfos = append(
		pm.providerInfos,
		ProviderInfo{Provider, Name, SleepBetweenRuns, SleepBeforeFirstRun})
}

func (pm *ProvidersManager) Run() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("can't load .env: %v", err)
		return
	}
	var API_URL = os.Getenv("API_URL")

	var wg sync.WaitGroup

	for _, info := range pm.providerInfos {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(info.SleepBeforeFirstRun)
			for {
				data, err := info.Provider.GetInfo()
				if err != nil {
					log.Error(err)
					continue
				}
				body, err := json.Marshal(data)
				if err != nil {
					log.Error(err)
					continue
				}


				req, err := http.NewRequest("POST", API_URL + info.Name, bytes.NewBuffer(body))
				if err != nil {
					log.Error(err)
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Error(err)
					continue
				}
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Error(err)
					continue
				}

				log.Debugf("status: %v", resp.Status)
				log.Debugf("response: %v", string(respBody))

				time.Sleep(info.SleepBetweenRuns)
			}
		}()
	}

	wg.Wait()
}
