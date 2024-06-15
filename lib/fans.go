package lib

import (
	"BiliAutoBlackList/config"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
)

type Fans struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    struct {
		List []struct {
			Mid            int         `json:"mid"`
			Attribute      int         `json:"attribute"`
			Mtime          int         `json:"mtime"`
			Tag            interface{} `json:"tag"`
			Special        int         `json:"special"`
			Uname          string      `json:"uname"`
			Face           string      `json:"face"`
			Sign           string      `json:"sign"`
			OfficialVerify struct {
				Type int    `json:"type"`
				Desc string `json:"desc"`
			} `json:"official_verify"`
			Vip struct {
				VipType       int    `json:"vipType"`
				VipDueDate    int    `json:"vipDueDate"`
				DueRemark     string `json:"dueRemark"`
				AccessStatus  int    `json:"accessStatus"`
				VipStatus     int    `json:"vipStatus"`
				VipStatusWarn string `json:"vipStatusWarn"`
				ThemeType     int    `json:"themeType"`
				Label         struct {
					Path string `json:"path"`
				} `json:"label"`
			} `json:"vip"`
		} `json:"list"`
		ReVersion int64 `json:"re_version"`
		Total     int   `json:"total"`
	} `json:"data"`
}

func GetFansAndCheck(page int) {
	logrus.Info("start check fans list page " + strconv.Itoa(page))
	pageSize := 50
	url := "https://api.bilibili.com/x/relation/fans?vmid=" + config.GConfig.TargetUID + "&pn=" + strconv.Itoa(page) + "&ps=" + strconv.Itoa(pageSize)
	bodyText := Request("GET", url, "")
	var fans Fans
	err := json.Unmarshal([]byte(bodyText), &fans)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range fans.Data.List {
		if v.Uname != "" {
			logrus.Info("check user: ", v.Uname)
			if IsInBlackList(v.Uname) {
				logrus.Info("add user: ", v.Uname)
				AddBlackList(v.Mid)
			}
		}
	}
	if page != 5 {
		time.Sleep(time.Duration(config.GConfig.TimeDelay) * time.Second)
		GetFansAndCheck(page + 1)
	} else {
		logrus.Info("check fans list complete")
	}
}
