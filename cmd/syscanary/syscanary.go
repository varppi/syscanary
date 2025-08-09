package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Varppi/syscanary/internal/config"
	"github.com/Varppi/syscanary/internal/global"
	"github.com/Varppi/syscanary/internal/logging"
	"github.com/Varppi/syscanary/internal/monitor"
)

func main() {
	// if os.Getuid() != 0 {
	// 	fmt.Println("Please run as root!")
	// 	return
	// }

	config, err := config.Parse()
	if err != nil {
		log.Fatal(err)
	}
	global.Config = config

	loggerConf := &logging.LoggerConf{
		Level:      config.Loglevel,
		OutputFile: config.Logfile,
	}
	logger, err := logging.InitLogger(loggerConf)
	if err != nil {
		log.Fatal(err)
	}
	global.Logger = logger

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	exit := make(chan *struct{}, 1)
	go func() {
		<-sigs
		exit <- nil
	}()

	monitorContext := &monitor.Context{
		Config: config,
	}
	go func() {
		err := monitor.Start(monitorContext)
		if err != nil {
			logger.Fatal(err)
			exit <- nil
		}
	}()

	<-exit
	monitor.Stop()
}
