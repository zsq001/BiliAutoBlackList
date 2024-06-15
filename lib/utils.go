package lib

import (
	"BiliAutoBlackList/config"
	"encoding/json"
	"log"
	"strings"
)

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
}

func GetCsrfToken() string {
	cookie := config.GConfig.Cookie
	cookieParts := strings.Split(cookie, "; ")
	var biliJct string
	for _, part := range cookieParts {
		if strings.Contains(part, "bili_jct") {
			biliJct = strings.Split(part, "=")[1]
			break
		}
	}
	return biliJct // 输出: bili_jct字段的值
}

func CheckResult(res string) bool {
	var result Result
	err := json.Unmarshal([]byte(res), &result)
	if err != nil {
		log.Fatal(err)
	}
	if result.Code == 0 {
		return true
	}
	return false
}
