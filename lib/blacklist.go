package lib

import (
	"BiliAutoBlackList/config"
	"github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

func IsInBlackList(content string) bool {
	for _, word := range config.GConfig.BlackListWord {
		if strings.Contains(content, word) {
			return true
		}
	}
	return false
}

func AddBlackList(mid int) {
	csrfToken := GetCsrfToken()
	res := Request("POST", "https://api.bilibili.com/x/relation/modify", "fid="+strconv.Itoa(mid)+"&act=5&re_src=11&csrf="+csrfToken)
	if CheckResult(res) {
		logrus.Info("add user success")
	} else {
		logrus.Fatal("add user failed")
	}
}
