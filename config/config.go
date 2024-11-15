package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Cookie          string   `yaml:"cookie"`
	TargetUID       string   `yaml:"targetUID"`
	BlackListWord   []string `yaml:"blackListWord"`
	TimeDelay       int64    `yaml:"timeDelay"`
	Cron            string   `yaml:"cron"`
	FansCheckPerDay int      `yaml:"fansCheckPerDay"`
	Mode            string   `yaml:"mode"`
	Feishu          struct {
		Email     string `yaml:"email"`
		AppId     string `yaml:"appId"`
		AppSecret string `yaml:"appSecret"`
		Token     string `yaml:"token"`
	} `yaml:"feishu"`
	OpenaiConfig struct {
		ApiBase string `yaml:"apiBase"`
		ApiKey  string `yaml:"apiKey"`
	} `yaml:"openaiConfig"`
}

var GConfig Config

// for base->lib.fans gpt->lark&lib.gpt gpt-only->lib.gpt(without lark)
func InitConfig() {

	if _, err := os.Stat("config.yaml"); os.IsNotExist(err) {
		defaultConfig := `cookie: "your_cookie_here"
targetUID: "your_uid_here"
timeDelay: 10
cron: "0 0 0/12 * * *"
mode: "gpt"
BlackListWord:
  - "word1"
  - "word2"
  - "word3"
feishu:
  appId: "your_app_id_here"
  appSecret: "your_app_secret_here"
  email: "your_email_here"
  token: "your_token_here"
openaiConfig:
  apiBase: "URL_ADDRESS"
  apiKey: "your_api_key_here"
fansCheckPerDay: 10
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
	if GConfig.OpenaiConfig.ApiBase == "" {
		GConfig.OpenaiConfig.ApiBase = "https://api.openai.com/v1/"
	}

	if GConfig.OpenaiConfig.ApiKey == "" {
		if os.Getenv("OPENAI_API_KEY") != "" {
			GConfig.OpenaiConfig.ApiKey = os.Getenv("OPENAI_API_KEY")
		} else {
			logrus.Panic("openai api key is empty, please check config.yaml")
		}
	}
}
