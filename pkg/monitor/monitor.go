package monitor

import (
	"fmt"
	"log"
	"net/http"
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

func (m Monitor) Check(wg *sync.WaitGroup) {
	if m.Config.RequestType == "HTTP" {
		m.checkHTTP()
	}
	if m.Config.RequestType == "ICMP" {
		m.checkICMP()
	}

	// TODO: Setup alerting here if something is wrong
	// TODO: Expose HTTP Status page
	if m.HasChecked {
		if m.Config.RequestType == "HTTP" {
			log.Default().Println(m.Config.Name, ": status code = ", m.StatusCode)
		} else if m.Config.RequestType == "ICMP" {
			log.Default().Println(m.Config.Name, ": response time = ", m.ResponseTime)
		}
	} else {
		log.Default().Println(m.Config.Name, " has not been checked. Please check the configuration")
	}
	wg.Done()
}
