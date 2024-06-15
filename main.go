package main

import (
	"BiliAutoBlackList/config"
	"BiliAutoBlackList/lib"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func main() {
	config.InitConfig()
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(config.GConfig.Cron, func() {
		lib.GetFansAndCheck(1)
	})
	logrus.Info("starting cron")
	if err != nil {
		panic(err)
	}
	c.Start()
	logrus.Info("program running")
	select {}
}
