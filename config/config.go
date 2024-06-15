package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
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

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		defaultConfig := `cookie: "your_cookie_here"
targetUID: "your_uid_here"
timeDelay: 10
cron: "0 0 0/12 * * *"
BlackListWord:
  - "word1"
  - "word2"
  - "word3"
`
		err = os.WriteFile("config_example.yaml", []byte(defaultConfig), 0644)
		if err != nil {
			logrus.Panicf("Failed to write default config file: %v", err)
		}
		logrus.Panic("config file not found, please create config.yaml based on config_example.yaml")
	}

	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		logrus.Panicf("Failed to read config file: %v", err)
	}

	err = yaml.Unmarshal(configFile, &GConfig)
	if err != nil {
		logrus.Panicf("Failed to unmarshal config file: %v", err)
	}

	if GConfig.Cron == "" {
		GConfig.Cron = "0 0 0/12 * * *"
	}

	if GConfig.TimeDelay < 10 {
		logrus.Warn("timeDelay is too short, may cause your account banned!")
	}
}
