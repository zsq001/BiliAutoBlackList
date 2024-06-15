package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"log"
	"os"
)

type Config struct {
	Cookie        string   `yaml:"cookie"`
	TargetUID     string   `yaml:"targetUID"`
	BlackListWord []string `yaml:"blackListWord"`
	TimeDelay     int64    `yaml:"timeDelay"`
	Cron          string   `yaml:"cron"`
}

var GConfig Config

func InitConfig() {
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Panicf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(configFile, &GConfig)
	if err != nil {
		log.Panicf("Failed to unmarshal config file: %v", err)
	}

	if GConfig.Cron == "" {
		GConfig.Cron = "0 0 0/12 * * *"
	}

	if GConfig.TimeDelay < 10 {
		logrus.Warn("timeDelay is too short, may cause your account banned!")
	}
}
