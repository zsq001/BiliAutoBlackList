package gpt

import (
	"BiliAutoBlackList/config"
	"context"
	"io/ioutil"
	"os"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/sirupsen/logrus"
)

var gptClient *openai.Client

var sysPrompt string

func InitGPTClient() {
	gptClient = openai.NewClient(
		option.WithBaseURL(config.GConfig.OpenaiConfig.ApiBase),
		option.WithAPIKey(config.GConfig.OpenaiConfig.ApiKey),
	)
	var err error
	_, err = os.Stat("prompt")
	if os.IsNotExist(err) {
		defaultPrompt := "你需要判断一个用户是否是可疑账号。如果认为不是则返回false，否则返回true。"
		err = ioutil.WriteFile("prompt.example", []byte(defaultPrompt), 0644)
		logrus.Fatalf("prompt file not found")
	}
	sysPrompt, err = ReadFileToString("prompt")
}

func JudgeResult(username string, sign string) bool {
	userPrompt := "username: " + username + "\n sign: " + sign
	logrus.Info("GPT: " + username)
	resp, err := gptClient.Chat.Completions.New(context.TODO(), openai.ChatCompletionNewParams{
		Model: openai.F(openai.ChatModelGPT4o),
		Messages: openai.F([]openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(sysPrompt),
			openai.UserMessage(userPrompt),
		}),
	})
	if err != nil {
		logrus.Errorf("Failed to get chat completion: %v", err)
	}
	if resp.Choices[0].Message.Content == "false" {
		return false
	} else {
		return true
	}
}

func ReadFileToString(filename string) (string, error) {
	// 打开文件
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// 读取文件内容到字节切片
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return "", err
	}

	// 将字节切片转换为字符串并返回
	return string(bytes), nil
}
