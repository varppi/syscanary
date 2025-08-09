package detections

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/Varppi/syscanary/internal/global"
)

var (
	internetStatus bool = true
)

func InternetDetect(output chan string, stop chan *struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	output <- "debug:internet:internet connectivity monitoring started"

	client := http.Client{
		Timeout: 5 * time.Second,
	}

	for {
		_, err := client.Get(fmt.Sprintf("https://%s", global.Config.Modules.Internet.Domain))
		if err != nil && internetStatus {
			internetStatus = false
			output <- fmt.Sprintf("info:internet:can't connect to %s", global.Config.Modules.Internet.Domain)
		} else if err == nil && !internetStatus {
			internetStatus = true
			output <- fmt.Sprintf("info:internet:connection to %s regained", global.Config.Modules.Internet.Domain)

		}

		if len(stop) > 0 {
			break
		}

		if global.Config.Modules.Ports.Interval < 10 {
			time.Sleep(10 * time.Second)
			continue
		}
		time.Sleep(time.Duration(global.Config.Modules.Internet.Interval) * time.Second)
	}
}
