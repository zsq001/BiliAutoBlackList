package lib

import (
	"BiliAutoBlackList/config"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"log"
	"strconv"
	"time"
)

//curl 'https://api.bilibili.com/x/web-interface/dynamic/entrance' \
//-H $'cookie: _uuid=810149F46-7CCC-24E10-1DFE-2810833C616FE78778infoc; buvid_fp=a5b13102af87ce248b3b95aad2ceafe9; DedeUserID=381720186; DedeUserID__ckMd5=787f53a8b82b50d2; rpdid=|(YuuJk~k)u0J\'uYmYmRkJYm; LIVE_BUVID=AUTO5616969433351369; buvid4=CF6EFFE5-2873-60D9-E7D7-C8CE64A457EC80270-022012520-iMstsERVusFz8XWIwDCeww%3D%3D; enable_web_push=ENABLE; iflogin_when_web_push=1; hit-dyn-v2=1; CURRENT_BLACKGAP=0; FEED_LIVE_VERSION=V_DYN_LIVING_UP; CURRENT_FNVAL=4048; CURRENT_QUALITY=112; buvid3=C2B1A582-8E90-38A0-E3F3-EE558AF3BFCD56196infoc; b_nut=1717998056; header_theme_version=CLOSE; PVID=2; SESSDATA=2e27aea9%2C1733819763%2C0b154%2A62CjBlGz1Pz-OnKejLW5M5-xC1uU-G_3mmyumshYx4AZmQVtIrarDGFXoqXyHsKahaSX0SVnpLNDA2MDJzR2k0QXNGREI4cE9Kdm8wRFNXaW9RTjlIRFJVVDlyWEQtd0poM0Z6V05KZzJyNmhZUnlLejlXLXUycGNFMDgxZVVRV2RzNVc3Rzg1cHB3IIEC; bili_jct=93c3663cb5dd620869963c556cf471ff; sid=7tlmmvl3; bsource=search_google; fingerprint=5cba5f5cefcba8be3fb89704db79067b; bili_ticket=eyJhbGciOiJIUzI1NiIsImtpZCI6InMwMyIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTg2MzM0NjAsImlhdCI6MTcxODM3NDIwMCwicGx0IjotMX0.iVUtrAUUKPbcAcesH1mpnBbePmm90nQui4vS-brCBJI; bili_ticket_expires=1718633400; bp_t_offset_381720186=942919366444318723; b_lsid=256A85D10_1901964AD4F; home_feed_column=4; browser_resolution=948-881' \

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
	if fans.Data.List != nil {
		time.Sleep(time.Duration(config.GConfig.TimeDelay) * time.Second)
		GetFansAndCheck(page + 1)
	}
}
