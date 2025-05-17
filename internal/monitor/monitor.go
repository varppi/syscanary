package monitor

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spoofimei/syscanary/internal/config"
	"github.com/spoofimei/syscanary/internal/detections"
	"github.com/spoofimei/syscanary/internal/global"
)

type Context struct {
	Config *config.ModuleConfig
}

var (
	wg           sync.WaitGroup
	stop         bool
	detectionMap = map[string]func(chan string, chan *struct{}, *sync.WaitGroup){
		"usb":       detections.UsbDetect,
		"integrity": detections.FileIntegrityDetect,
		"ports":     detections.OpenPortsDetect,
		"internet":  detections.InternetDetect,
	}
)

func Start(context *Context) error {
	if len(context.Config.Detections) == 0 {
		return errors.New("no detection modules selected")
	}

	var detections []func(chan string, chan *struct{}, *sync.WaitGroup)
	for _, detection := range context.Config.Detections {
		fun, ok := detectionMap[detection]
		if !ok {
			return fmt.Errorf("invalid detection module \"%s\"", detection)
		}
		detections = append(detections, fun)
	}

	fmt.Println("monitoring started")
	global.Logger.Debug("monitoring started")

	outputChan := make(chan string, 1000)
	stopChan := make(chan *struct{}, 1)
	for _, detectionFunc := range detections {
		wg.Add(1)
		go detectionFunc(outputChan, stopChan, &wg)
	}

	for {
		if stop {
			break
		}
		if len(outputChan) != 0 {
			output := <-outputChan
			level := strings.Split(output, ":")[0]
			module := strings.Split(output, ":")[1]
			msg := strings.Join(strings.Split(output, ":")[2:], ":")
			switch level {
			case "debug":
				global.Logger.Debug(fmt.Sprintf("%s: %s", module, msg))
			case "info":
				global.Logger.Info(fmt.Sprintf("%s: %s", module, msg))
			case "error":
				global.Logger.Error(fmt.Errorf("%s: %s", module, msg))
			}
		}
		time.Sleep(100 * time.Millisecond)
	}
	stopChan <- nil
	global.Logger.Debug("stopping all processes")
	fmt.Println("monitoring stopped")
	return nil
}

func Stop() {
	stop = true
	wg.Wait()
}
