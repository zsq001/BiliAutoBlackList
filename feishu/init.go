package feishu

import (
	"BiliAutoBlackList/config"
	"BiliAutoBlackList/lib"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	lark "github.com/larksuite/oapi-sdk-go/v3"
	larkcontact "github.com/larksuite/oapi-sdk-go/v3/service/contact/v3"
	"github.com/sirupsen/logrus"
)

type Challenge struct {
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
	Token     string `json:"token"`
}

type Pair struct {
	Uid      int
	username string
}

var userOpenId string

var client *lark.Client

// save messageid and uid and username for rej&pass card
var messageId2Uid map[string]Pair

var ExceptList map[int]bool

type RecallResponse struct {
	OpenID        string `json:"open_id"`
	UserID        string `json:"user_id"`
	OpenMessageID string `json:"open_message_id"`
	TenantKey     string `json:"tenant_key"`
	Token         string `json:"token"`
	Action        struct {
		Value     string      `json:"value"`
		Tag       string      `json:"tag"`
		Timezone  string      `json:"timezone"`
		FormValue interface{} `json:"form_value"`
		Name      string      `json:"name"`
	} `json:"action"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(messageId2Uid)
	// 确保请求方法是 POST
	if r.Method != http.MethodPost {
		logrus.Errorf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 读取请求体
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	// 解析 JSON 数据
	var data RecallResponse
	logrus.Info("body: ", string(body))
	err = json.Unmarshal(body, &data)
	//logrus.Info(data)
	if err != nil {
		logrus.Errorf("Failed to parse JSON: %v", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	mid := data.OpenMessageID
	key := data.Action.Value

	//logrus.Info("Received card action:", data.Action.Value)
	logrus.Info("mid: ", mid)
	logrus.Info("messageId2Uid[mid].username:", messageId2Uid[mid].username)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if key == "\"reject\"" {
		uid := messageId2Uid[mid].Uid
		lib.AddBlackList(uid)
		//delete(messageId2Uid, mid)
		w.Write([]byte(GetRejCard(messageId2Uid[mid].username)))
	} else if key == "\"pass\"" {
		uid := messageId2Uid[mid].Uid
		ExceptList[uid] = true
		//delete(messageId2Uid, mid)
		w.Write([]byte(GetPassCard(messageId2Uid[mid].username)))
	}
}

func Init() {
	client = lark.NewClient(config.GConfig.Feishu.AppId, config.GConfig.Feishu.AppSecret)

	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(`open_id`).
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Emails([]string{config.GConfig.Feishu.Email}).
			Build()).
		Build()
	// 发起请求
	resp, err := client.Contact.User.BatchGetId(context.Background(), req)

	messageId2Uid = make(map[string]Pair)
	ExceptList = make(map[int]bool)

	if err != nil {
		logrus.Fatalf("Failed to get user info: %v", err)
	}

	userOpenId = *resp.Data.UserList[0].UserId

	http.HandleFunc("/webhook/card", handler)

	go http.ListenAndServe(":9999", nil)
	if err != nil {
		logrus.Fatalf("Failed to start webhook server: %v", err)
	}
}

/* messageid uid做关联

如果用户拉黑，就拉黑对应uid

如果用户取消拉黑，就删除这对关联，并将对应uid加入例外列表
*/
