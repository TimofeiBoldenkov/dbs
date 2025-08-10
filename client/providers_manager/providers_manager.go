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

	"github.com/TimofeiBoldenkov/dbs/client/info_provider"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type ProviderInfo struct {
	Provider            infoprovider.InfoProvider
	Tag					string
	SleepBetweenRuns    time.Duration
	SleepBeforeFirstRun time.Duration
}

// Manages the providers' launch according to their SleepBetweenRuns and SleepBeforeFirstRun variables.
// Sends to server values returned by providers as json.
type ProvidersManager struct {
	providerInfos []ProviderInfo
}

// Same as AddSetDelay, but sets the sleepBeforeFirstRun variable 
// to a random value between 0 and SleepBetweenRuns
func (pm *ProvidersManager) Add(
	Provider infoprovider.InfoProvider, 
	Tag string, 
	SleepBetweenRuns time.Duration) {
	sleepBeforeFirstRun := time.Duration(rand.Float64() * float64(SleepBetweenRuns))

	pm.AddSetDelay(Provider, Tag, SleepBetweenRuns, sleepBeforeFirstRun)
}

// The Tag argument is used to identify the Provider in the database
func (pm *ProvidersManager) AddSetDelay(
	Provider infoprovider.InfoProvider,
	Tag string,
	SleepBetweenRuns time.Duration,
	SleepBeforeFirstRun time.Duration) {

	pm.providerInfos = append(
		pm.providerInfos,
		ProviderInfo{Provider, Tag, SleepBetweenRuns, SleepBeforeFirstRun})
}

// Runs providers according to their SleepBetweenRuns and SleepBeforeFirstRun variables.
// Each provider is run in a separate goroutine.
// If an error occurs, it is printed to logs, nothing is sent to the server.
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
					time.Sleep(info.SleepBetweenRuns)
					continue
				}
				body, err := json.Marshal(data)
				if err != nil {
					log.Error(err)
					time.Sleep(info.SleepBetweenRuns)
					continue
				}


				req, err := http.NewRequest("POST", API_URL + info.Tag, bytes.NewBuffer(body))
				if err != nil {
					log.Error(err)
					time.Sleep(info.SleepBetweenRuns)
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{}
				resp, err := client.Do(req)
				if err != nil {
					log.Error(err)
					time.Sleep(info.SleepBetweenRuns)
					continue
				}
				respBody, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Error(err)
					time.Sleep(info.SleepBetweenRuns)
					continue
				}

				log.Debugf("%v status: %v", info.Tag, resp.Status)
				log.Debugf("%v response: %v", info.Tag, string(respBody))

				time.Sleep(info.SleepBetweenRuns)
			}
		}()
	}

	wg.Wait()
}
