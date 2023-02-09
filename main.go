package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
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
	m.Status, m.ResponseTime = true, pinger.Statistics().AvgRtt

}

func (m *Monitor) checkHTTP() {
	if m.Config.Port == 0 {
		m.Config.Port = 80
	}
	resp, err := http.Get(fmt.Sprintf("%s:%d", m.Config.Hostname, m.Config.Port))

	if err != nil {
		log.Default().Println(err)
		m.Status, m.StatusCode = false, 0

		return
	}
	defer resp.Body.Close()
	m.Status, m.StatusCode = true, resp.StatusCode

}

func main() {
	rawConfigFile, err := os.ReadFile("./config.json") // TODO: Replace config file path using cli args -c /path/to/config/file

	if err != nil {
		panic("Cannot read config file")
	}
	var configs []MonitorConfig

	json.Unmarshal(rawConfigFile, &configs)
	for _, config := range configs {
		monitor := Monitor{Config: config}

		if monitor.Config.RequestType == "HTTP" {
			monitor.checkHTTP()
		}
		if monitor.Config.RequestType == "ICMP" {
			monitor.checkICMP()
		}

		// TODO: Setup alerting here if something is wrong
		// TODO: Expose HTTP Status page
		log.Default().Println(monitor)
	}
}
