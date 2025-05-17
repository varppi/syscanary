package detections

import (
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/spoofimei/syscanary/internal/global"
)

var (
	portsList map[string]bool
)

func OpenPortsDetect(output chan string, stop chan *struct{}, wg *sync.WaitGroup) {
	defer wg.Done()
	output <- "debug:ports:open port monitoring started"

	for {
		out, err := exec.Command("sh", "-c", "netstat -tulpn 2>/dev/null |awk '/:/{print $4}'").Output()
		if err != nil {
			output <- fmt.Sprintf("error:ports:%s", err.Error())
		}
		newPortsList := make(map[string]bool)
		for _, bind := range strings.Split(string(out), "\n") {
			if (strings.HasPrefix(bind, "127.0.0") || strings.HasPrefix(bind, "::")) && global.Config.Modules.Ports.Ignorelocal {
				continue
			}
			if _, ok := portsList[bind]; !ok && len(portsList) != 0 {
				output <- fmt.Sprintf("info:ports:new bind: %s", bind)
			}
			newPortsList[bind] = true
		}
		for bind := range portsList {
			if _, ok := newPortsList[bind]; !ok {
				output <- fmt.Sprintf("info:ports:removed bind: %s", bind)
			}
		}
		portsList = newPortsList

		if len(stop) > 0 {
			break
		}
		time.Sleep(time.Duration(global.Config.Modules.Ports.Interval) * time.Second)
	}
}
