package lib

import (
	"BiliAutoBlackList/config"
	"io"
	"log"
	"net/http"
	"strings"
)

func Request(method string, url string, datas string) string {
	client := &http.Client{}
	var data = strings.NewReader(datas)
	req, err := http.NewRequest(method, url, data)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("cookie", config.GConfig.Cookie)
	if method == "POST" {
		req.Header.Set("content-type", "application/x-www-form-urlencoded")
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyText)
}
