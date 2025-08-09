package syscanary_tests

import (
	"testing"
	"time"

	"github.com/Varppi/syscanary/internal/config"
	"github.com/Varppi/syscanary/internal/global"
	"github.com/Varppi/syscanary/internal/logging"
	"github.com/Varppi/syscanary/internal/monitor"
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
