package main

import (
	"BiliAutoBlackList/config"
	"BiliAutoBlackList/lib"
	"github.com/robfig/cron/v3"
)

func main() {
	config.InitConfig()
	c := cron.New(cron.WithSeconds())
	_, err := c.AddFunc(config.GConfig.Cron, func() {
		lib.GetFansAndCheck(1)
	})
	if err != nil {
		panic(err)
	}
	c.Start()
	select {}
}
