package main

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	monitor "github.com/ljahier/uptime-monitor/pkg/monitor"
	webserver "github.com/ljahier/uptime-monitor/pkg/webserver"
)

func main() {
	var wg sync.WaitGroup
	rawConfigFile, err := os.ReadFile("./config.json") // TODO: Replace config file path using cli args -c /path/to/config/file

	if err != nil {
		panic("Cannot read config file")
	}
	var configs []monitor.MonitorConfig

	json.Unmarshal(rawConfigFile, &configs)

	go webserver.RunWebServer()

	for {
		wg.Add(len(configs))

		for _, config := range configs {
			monitor := monitor.Monitor{Config: config}
			go monitor.Check(&wg)
		}
		wg.Wait()
		time.Sleep(1 * time.Minute)
	}
}
