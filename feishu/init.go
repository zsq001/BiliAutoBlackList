package feishu

import (
	"BiliAutoBlackList/config"
	"BiliAutoBlackList/lib"
	"context"
	"encoding/json"
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

// messageId2Uid save messageid and uid and username for rej&pass card
var messageId2Uid map[string]Pair

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

type Preset struct {
	Challenge string `json:"challenge"`
	Token     string `json:"token"`
	Type      string `json:"type"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		logrus.Errorf("Method not allowed: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		logrus.Errorf("Failed to read request body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var data RecallResponse
	logrus.Debug("body: ", string(body))
	err = json.Unmarshal(body, &data)
	if err != nil || data.OpenID == "" {
		var precheck Preset
		if err = json.Unmarshal(body, &precheck); err == nil {
			if precheck.Token == config.GConfig.Feishu.Token {
				res := map[string]string{
					"challenge": precheck.Challenge,
				}
				jsonData, _ := json.Marshal(res)
				w.Write(jsonData)
				return
			}
		}
		logrus.Errorf("Failed to parse JSON: %v", err)
		http.Error(w, "Failed to parse JSON", http.StatusBadRequest)
		return
	}
	mid := data.OpenMessageID
	key := data.Action.Value

	logrus.Debug("mid: ", mid)
	logrus.Debug("messageId2Uid[mid].username:", messageId2Uid[mid].username)
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
		saveExceptList()
		//delete(messageId2Uid, mid)
		w.Write([]byte(GetPassCard(messageId2Uid[mid].username)))
	}
}

func Init() {
	client = lark.NewClient(config.GConfig.Feishu.AppId, config.GConfig.Feishu.AppSecret)

	initExceptList()
	req := larkcontact.NewBatchGetIdUserReqBuilder().
		UserIdType(`open_id`).
		Body(larkcontact.NewBatchGetIdUserReqBodyBuilder().
			Emails([]string{config.GConfig.Feishu.Email}).
			Build()).
		Build()
	// request
	resp, err := client.Contact.User.BatchGetId(context.Background(), req)

	messageId2Uid = make(map[string]Pair)

	if err != nil {
		logrus.Fatalf("Failed to get user info: %v", err)
	}
	// get user's openid
	userOpenId = *resp.Data.UserList[0].UserId

	http.HandleFunc("/", handler)

	go http.ListenAndServe(":9999", nil)
	if err != nil {
		logrus.Fatalf("Failed to start webhook server: %v", err)
	}
}
