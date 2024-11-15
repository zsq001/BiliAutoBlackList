package main

import (
	"BiliAutoBlackList/config"
	"BiliAutoBlackList/fans"
	"BiliAutoBlackList/feishu"
	"BiliAutoBlackList/gpt"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"strings"
)

func main() {
	config.InitConfig()

	if strings.Contains(config.GConfig.Mode, "gpt") {
		gpt.InitGPTClient()
	}
	feishu.Init()

	c := cron.New(cron.WithSeconds())

	_, err := c.AddFunc(config.GConfig.Cron, func() {
		fans.GetFansAndCheck(1)
	})

	fans.GetFansAndCheck(1)

	logrus.Info("starting cron")
	if err != nil {
		panic(err)
	}
	c.Start()
	logrus.Info("program running")
	select {}

}
