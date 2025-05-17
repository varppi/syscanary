package config

import (
	"errors"

	"github.com/spf13/viper"
)

// 3rd layer
type Integrity struct {
	Paths    []string
	Interval int
}
type Usb struct {
	Interval int
}
type Ports struct {
	Ignorelocal bool
	Interval    int
}
type Internet struct {
	Domain   string
	Interval int
}

// 2nd layer
type Modules struct {
	Integrity Integrity
	Usb       Usb
	Ports     Ports
	Internet  Internet
}

// 1st layer
type ModuleConfig struct {
	Logfile    string
	Loglevel   int
	Detections []string
	Modules    Modules
}

func Parse() (*ModuleConfig, error) {
	viper.SetConfigName("syscanary")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	moduleConf := viper.GetStringMap("modules")
	if len(moduleConf) == 0 {
		return nil, errors.New("no module configurations")
	}

	moduleConfig := &ModuleConfig{}
	err = viper.Unmarshal(moduleConfig)
	if err != nil {
		return nil, err
	}
	return moduleConfig, nil
}
