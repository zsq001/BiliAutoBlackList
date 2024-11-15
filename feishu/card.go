package feishu

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/larksuite/oapi-sdk-go/v3/core"
	larkim "github.com/larksuite/oapi-sdk-go/v3/service/im/v1"
	"github.com/sirupsen/logrus"
)

type CardContent struct {
	Type string `json:"type"`
	Data struct {
		TemplateID          string            `json:"template_id"`
		TemplateVersionName string            `json:"template_version_name"`
		TemplateVariable    map[string]string `json:"template_variable"`
	} `json:"data"`
}

func SendQueryCard(un string, sign string, uid int) {
	origin := "{\"config\":{\"update_multi\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"markdown\",\"content\":\"用户：\\\"${username}\\\" 被识别为可疑账号；\\n该用户个性签名为：\\\"${sign}\\\";\\n该用户个人主页为  <a href=https://space.bilibili.com/${uid}>https://space.bilibili.com/${uid}</a>\",\"text_align\":\"left\",\"text_size\":\"normal\"},{\"tag\":\"column_set\",\"flex_mode\":\"none\",\"horizontal_spacing\":\"8px\",\"horizontal_align\":\"left\",\"columns\":[{\"tag\":\"column\",\"width\":\"auto\",\"vertical_align\":\"top\",\"vertical_spacing\":\"8px\",\"elements\":[{\"tag\":\"button\",\"text\":{\"tag\":\"plain_text\",\"content\":\"通过\"},\"type\":\"primary_filled\",\"width\":\"default\",\"size\":\"medium\",\"icon\":{\"tag\":\"standard_icon\",\"token\":\"list-check-bold_outlined\"},\"confirm\":{\"title\":{\"tag\":\"plain_text\",\"content\":\"二次确认\"},\"text\":{\"tag\":\"plain_text\",\"content\":\"确认放行？\"}},\"behaviors\":[{\"type\":\"callback\",\"value\":\"pass\"}]}]},{\"tag\":\"column\",\"width\":\"auto\",\"vertical_align\":\"top\",\"vertical_spacing\":\"8px\",\"elements\":[{\"tag\":\"button\",\"text\":{\"tag\":\"plain_text\",\"content\":\"拉黑\"},\"type\":\"danger_filled\",\"width\":\"default\",\"size\":\"medium\",\"icon\":{\"tag\":\"standard_icon\",\"token\":\"close-bold_outlined\"},\"confirm\":{\"title\":{\"tag\":\"plain_text\",\"content\":\"二次确认\"},\"text\":{\"tag\":\"plain_text\",\"content\":\"确认拉黑？\"}},\"behaviors\":[{\"type\":\"callback\",\"value\":\"reject\"}]}]}],\"margin\":\"16px 0px 0px 0px\"}]},\"i18n_header\":{\"zh_cn\":{\"title\":{\"tag\":\"plain_text\",\"content\":\"可疑账号需要判断\"},\"subtitle\":{\"tag\":\"plain_text\",\"content\":\"\"},\"template\":\"blue\"}}}"
	content := strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(origin, "${username}", un), "${sign}", sign), "${uid}", strconv.Itoa(uid))
	strconv.Quote(content)
	logrus.Debug(content)
	//send card via openid
	req := larkim.NewCreateMessageReqBuilder().
		ReceiveIdType(`open_id`).
		Body(larkim.NewCreateMessageReqBodyBuilder().
			ReceiveId(userOpenId).
			MsgType(`interactive`).
			Content(content).
			Build()).
		Build()

	resp, err := client.Im.Message.Create(context.Background(), req)
	if err != nil {
		logrus.Error(err)
	}

	if !resp.Success() {
		fmt.Printf("logId: %s, error response: \n%s", resp.RequestId(), larkcore.Prettify(resp.CodeError))
		logrus.Error(errors.New(resp.CodeError.Msg))
	}
	logrus.Debug("Query message openid: ", *resp.Data.MessageId)

	//add uid to pending confirm queue
	messageId2Uid[*resp.Data.MessageId] = Pair{uid, un}
	logrus.Debug("user stored in messageId2Uid pair: ", messageId2Uid[*resp.Data.MessageId])
}

func GetRejCard(un string) string {
	origin := "{\"config\":{\"update_multi\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"markdown\",\"content\":\"用户：\\\"${username}\\\" 已被拉黑\",\"text_align\":\"left\",\"text_size\":\"normal\"}]},\"i18n_header\":{\"zh_cn\":{\"title\":{\"tag\":\"plain_text\",\"content\":\"已拉黑\"},\"subtitle\":{\"tag\":\"plain_text\",\"content\":\"\"},\"template\":\"red\"}}}"
	content := strings.ReplaceAll(origin, "${username}", un)
	strconv.Quote(content)
	logrus.Debug("rej card content:", content)
	return content
}

func GetPassCard(un string) string {
	origin := "{\"config\":{\"update_multi\":true},\"i18n_elements\":{\"zh_cn\":[{\"tag\":\"markdown\",\"content\":\"用户：\\\"${username}\\\" 已放行\",\"text_align\":\"left\",\"text_size\":\"normal\"}]},\"i18n_header\":{\"zh_cn\":{\"title\":{\"tag\":\"plain_text\",\"content\":\"已放行\"},\"subtitle\":{\"tag\":\"plain_text\",\"content\":\"\"},\"template\":\"green\"}}}"
	content := strings.ReplaceAll(origin, "${username}", un)
	strconv.Quote(content)
	logrus.Debug("pass card content:", content)
	return content
}
