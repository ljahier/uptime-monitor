package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	monitor "github.com/ljahier/uptime-monitor/pkg/monitor"
	webserver "github.com/ljahier/uptime-monitor/pkg/webserver"
	launchpad "github.com/rakyll/launchpad"
)

func main() {
	var wg sync.WaitGroup
	rawConfigFile, err := os.ReadFile("./config.json") // TODO: Replace config file path using cli args -c /path/to/config/file
	if err != nil {
		panic("Cannot read config file")
	}
	var configs []monitor.MonitorConfig

	json.Unmarshal(rawConfigFile, &configs)

	// var data = make(chan []monitor.Monitor)

	go webserver.RunWebServer()

	pad, err := launchpad.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer pad.Close()

	pad.Clear()

	for {
		y := 0
		x := 0
		wg.Add(len(configs))
		c := make(chan monitor.Monitor)
		for i, config := range configs {
			monitor := monitor.Monitor{Config: config}
			go monitor.Check(&wg, c)
			tmp := <-c // TODO: probleme de mutex, 3 et 4 ecrivent en meme temps sur "c" ce qui fait que seulment une des deux value est prise en compte
			fmt.Println(i, " = ", tmp)
			if i > 8 {
				y += 1
				x = i - 8
			} else {
				x = i
			}
			if tmp.HasChecked {
				if tmp.Config.RequestType == "HTTP" && tmp.StatusCode >= 200 && tmp.StatusCode <= 299 {
					pad.Light(x, y, 3, 0)
				} else if tmp.Config.RequestType == "ICMP" {
					if tmp.ResponseTime > 15*time.Millisecond {
						pad.Light(x, y, 3, 3)
					} else {
						pad.Light(x, y, 3, 0)
					}
				}
			}
		}

		wg.Wait()
		time.Sleep(5 * time.Second)
	}
}
