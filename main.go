package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	probing "github.com/prometheus-community/pro-bing"
)

type MonitorConfig struct {
	Name            string `json:"name"`
	Hostname        string `json:"hostname"`
	Port            int16  `json:"port"`
	AlertThreshold  int8   `json:"alertThreshold"`
	PollingInterval int16  `json:"pollingInterval"`
	RequestType     string `json:"requestType"`
}

type Monitor struct {
	Status       bool
	StatusCode   int
	ResponseTime time.Duration
	Config       MonitorConfig
	HasChecked   bool
}

func (m *Monitor) checkICMP() {
	pinger, err := probing.NewPinger(m.Config.Hostname)

	if err != nil {
		log.Default().Println(err)
	}
	pinger.Count = 3
	err = pinger.Run()

	if err != nil {
		log.Default().Println(err)
	}

	m.HasChecked, m.Status, m.ResponseTime = true, true, pinger.Statistics().AvgRtt
}

func (m *Monitor) checkHTTP() {
	if m.Config.Port == 0 {
		m.Config.Port = 80
	}
	// TODO: Calcul response time
	resp, err := http.Get(fmt.Sprintf("%s:%d", m.Config.Hostname, m.Config.Port))

	if err != nil {
		m.Status, m.StatusCode = false, 0
		return
	}
	defer resp.Body.Close()
	m.HasChecked, m.Status, m.StatusCode = true, true, resp.StatusCode
}

func (m Monitor) check(wg *sync.WaitGroup) {
	if m.Config.RequestType == "HTTP" {
		m.checkHTTP()
	}
	if m.Config.RequestType == "ICMP" {
		m.checkICMP()
	}

	// TODO: Setup alerting here if something is wrong
	// TODO: Expose HTTP Status page
	if m.HasChecked {
		log.Default().Println(m.Config.Name, ": ", m.ResponseTime)
	} else {
		log.Default().Println(m.Config.Name, " has not been checked. Please verify the configuration")
	}
	wg.Done()
}

func main() {
	for {
		var wg sync.WaitGroup
		rawConfigFile, err := os.ReadFile("./config.json") // TODO: Replace config file path using cli args -c /path/to/config/file

		if err != nil {
			panic("Cannot read config file")
		}
		var configs []MonitorConfig

		json.Unmarshal(rawConfigFile, &configs)
		wg.Add(len(configs))

		for _, config := range configs {
			monitor := Monitor{Config: config}
			go monitor.check(&wg)
		}
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
