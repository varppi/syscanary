package syscanary_tests

import (
	"testing"

	"github.com/spoofimei/syscanary/internal/config"
)

func TestConfigParse(t *testing.T) {
	config, err := config.Parse()
	if err != nil {
		t.Error(err)
	}
	if config.Detections[0]+config.Detections[1] != "usbintegrity" {
		t.Error("invalid result from viper when parsing enabled detections")
	}
}
