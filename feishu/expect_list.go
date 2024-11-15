package feishu

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
)

// ExceptList when pass, save uid in except list
var ExceptList map[int]bool
var exceptListFile = "except_list.json"

func initExceptList() {
	ExceptList = make(map[int]bool)
	loadExceptList()
}

func loadExceptList() {
	file, err := os.Open(exceptListFile)
	if err != nil {
		if os.IsNotExist(err) {
			logrus.Info("ExceptList file does not exist, creating a new one.")
			return
		}
		logrus.Errorf("Failed to open ExceptList file: %v", err)
		return
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		logrus.Errorf("Failed to read ExceptList file: %v", err)
		return
	}

	err = json.Unmarshal(data, &ExceptList)
	if err != nil {
		logrus.Errorf("Failed to parse ExceptList JSON: %v", err)
		return
	}
}

func saveExceptList() {
	file, err := os.Create(exceptListFile)
	if err != nil {
		logrus.Errorf("Failed to create ExceptList file: %v", err)
		return
	}
	defer file.Close()

	data, err := json.Marshal(ExceptList)
	if err != nil {
		logrus.Errorf("Failed to marshal ExceptList to JSON: %v", err)
		return
	}

	_, err = file.Write(data)
	if err != nil {
		logrus.Errorf("Failed to write ExceptList to file: %v", err)
		return
	}
}

func IsInExceptList(uid int) bool {
	return ExceptList[uid]
}
