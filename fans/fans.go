package fans

import (
	"BiliAutoBlackList/config"
	"BiliAutoBlackList/feishu"
	"BiliAutoBlackList/gpt"
	"BiliAutoBlackList/lib"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
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

	pageSize := 10

	url := "https://api.bilibili.com/x/relation/fans?vmid=" + config.GConfig.TargetUID + "&pn=" + strconv.Itoa(page) + "&ps=" + strconv.Itoa(pageSize)
	bodyText := lib.Request("GET", url, "")
	var fans Fans
	err := json.Unmarshal([]byte(bodyText), &fans)
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range fans.Data.List {
		logrus.Debug(feishu.IsInExceptList(v.Mid))
		if v.Uname != "" && v.Vip.VipStatus == 0 && !feishu.IsInExceptList(v.Mid) {
			logrus.Info("check user: ", v.Uname)
			if config.GConfig.Mode == "basic" && lib.IsInBlackList(v.Uname) {
				logrus.Info("[basic] add user: ", v.Uname)
				lib.AddBlackList(v.Mid)
				continue
			} else if config.GConfig.Mode == "gpt" && gpt.JudgeResult(v.Uname, v.Sign) {
				logrus.Info("[gpt] push user: ", v.Uname)
				feishu.SendQueryCard(v.Uname, v.Sign, v.Mid)
				continue
			} else if config.GConfig.Mode == "gpt-only" && gpt.JudgeResult(v.Uname, v.Sign) {
				logrus.Info("[gpt-only] add user: ", v.Uname)
				lib.AddBlackList(v.Mid)
				continue
			}

			//save tokens, the whitelist can be shared between multiple accounts
			if strings.Contains(config.GConfig.Mode, "gpt") {
				feishu.ExceptList[v.Mid] = true
				feishu.SaveExceptList()
			}

		}
	}
	if page <= config.GConfig.FansCheckPerDay/pageSize {
		time.Sleep(time.Duration(config.GConfig.TimeDelay) * time.Second)
		GetFansAndCheck(page + 1)
	} else {
		logrus.Info("check fans list complete")
	}
}
