package logging

import (
	"errors"
	"io"
	"log"
	"os"
)

type LoggerConf struct {
	OutputFile   string
	CustomWriter io.Writer
	Level        int // 1=debug 2=info 3=error
}

type Logger struct {
	Conf   *LoggerConf
	logger *log.Logger
}

func InitLogger(conf *LoggerConf) (*Logger, error) {
	loggerInstance := &Logger{
		Conf:   conf,
		logger: log.Default(),
	}
	if conf.CustomWriter != nil && conf.OutputFile != "" {
		return nil, errors.New("can't have both CustomWriter and OutputFile defined")
	}
	if conf.OutputFile != "" {
		outputHandle, err := os.OpenFile(conf.OutputFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			return nil, err
		}
		loggerInstance.logger.SetOutput(outputHandle)
	}
	if conf.CustomWriter != nil {
		loggerInstance.logger.SetOutput(conf.CustomWriter)
	}
	return loggerInstance, nil
}

func (l *Logger) Debug(msg string) {
	if l.Conf.Level > 0 {
		return
	}
	l.logger.Println(msg)
}
func (l *Logger) Info(msg string) {
	if l.Conf.Level > 1 {
		return
	}
	l.logger.Println(msg)
}
func (l *Logger) Error(err error) {
	if l.Conf.Level > 2 {
		return
	}
	l.logger.Printf("Error: %s", err.Error())
}
func (l *Logger) Fatal(err error) {
	l.logger.Fatalf("Fatal Error: %s", err.Error())
}
