package syscanary_tests

import (
	"bytes"
	"errors"
	"testing"

	"github.com/Varppi/syscanary/internal/logging"
)

func TestLoggerDebug(t *testing.T) {
	buf := new(bytes.Buffer)
	config := &logging.LoggerConf{
		CustomWriter: buf,
		Level:        0,
	}
	logger, err := logging.InitLogger(config)
	if err != nil {
		t.Error(err)
	}
	logger.Debug("test")
	if buf.Len() != 25 {
		t.Error(errors.New("logger didn't write"))
	}
	logger.Info("test")
	if buf.Len() != 25*2 {
		t.Error(errors.New("logger didn't write"))
	}
	logger.Error(errors.New("test"))
	if buf.Len() != 25*3+7 {
		t.Error(errors.New("logger didn't write"))
	}
	buf.Reset()

}

func TestLoggerError(t *testing.T) {
	buf := new(bytes.Buffer)
	config := &logging.LoggerConf{
		CustomWriter: buf,
		Level:        2,
	}
	logger, err := logging.InitLogger(config)
	if err != nil {
		t.Error(err)
	}
	logger.Debug("test")
	if buf.Len() != 0 {
		t.Error(errors.New("logger wrote when it shouldn't had"))
	}
	logger.Info("test")
	if buf.Len() != 0 {
		t.Error(errors.New("logger wrote when it shouldn't had"))
	}
	logger.Error(errors.New("test"))
	if buf.Len() != 25+7 {
		t.Error(errors.New("logger didn't write"))
	}
	buf.Reset()
}
