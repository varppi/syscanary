package syscanary_tests

import (
	"testing"
	"time"

	"github.com/spoofimei/syscanary/internal/config"
	"github.com/spoofimei/syscanary/internal/global"
	"github.com/spoofimei/syscanary/internal/logging"
	"github.com/spoofimei/syscanary/internal/monitor"
)

func TestMonitor(t *testing.T) {
	loggerConf := &logging.LoggerConf{
		Level: 1,
	}

	logger, err := logging.InitLogger(loggerConf)
	if err != nil {
		t.Error(err)
	}
	global.Logger = logger

	config, err := config.Parse()
	if err != nil {
		t.Error(err)
	}
	global.Config = config

	monitorContext := &monitor.Context{
		Config: config,
	}
	go func() {
		err = monitor.Start(monitorContext)
		if err != nil {
			t.Error(err)
		}
	}()
	time.Sleep(2 * time.Second)
	monitor.Stop()
}
