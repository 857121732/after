package robot

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"after/chat"
)

var (
	confFile = flag.String("c", "conf.yml", "conf file path")
)

func calcSign(timestamp, secret string) string {
	signStr := timestamp + "\n" + secret
	return base64.StdEncoding.EncodeToString([]byte(hmacSha256(signStr, secret)))
}

func hmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func AitMe(w http.ResponseWriter, req *http.Request) {
	ts := req.Header.Get("timestamp")
	sign := req.Header.Get("sign")
	mySign := calcSign(ts, globalConf.AppSecret)
	if sign != mySign {
		errStr := fmt.Sprintf("签名校验失败。ts=%s&sign=%s, mySign=%s", ts, sign, mySign)
		log.Println(errStr)
		//w.Write([]byte(errStr))
		//return TODO 校验失败
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		errStr := fmt.Sprintf("读取body失败。err=%v", err)
		log.Println(errStr)
		w.Write([]byte(errStr))
		return
	}
	defer req.Body.Close()
	log.Printf("Get body:[%s]\n", string(body))

	request := &Request{}
	err = json.Unmarshal(body, request)
	if err != nil {
		errStr := fmt.Sprintf("body(%s)解析失败(%v)", string(body), err)
		log.Println(errStr)
		w.Write([]byte(errStr))
		return
	}

	response := chat.GetResponse(strings.TrimPrefix(request.Text.Content, " "))
	responseDingTalk(response)
	w.Write([]byte("ok"))
}

func responseDingTalk(text string) {
	body := fmt.Sprintf(`
{
    "msgtype":"text",
    "text":{
        "content":"TEST: %s"
    }
}
`, text)
	resp, err := http.Post(globalConf.WebHook, "application/json", strings.NewReader(body))
	if err != nil {
		log.Printf("Send ding talk body(%s) error(%v).\n", body, err)
		return
	}
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Read resp body error(%v).\n", err)
		return
	}
	_ = resp.Body.Close()
	log.Printf("Send ding talk resp: %s.\n", string(resBody))
}

func Serve() {
	flag.Parse()
	conf, err := parseConf(*confFile)
	if err != nil {
		panic(err)
	}
	chat.InitChat(conf.CorporaPath, 10)

	http.HandleFunc("/ait", AitMe)
	err = http.ListenAndServe(":"+conf.Port, nil)
	if err != nil {
		panic(err)
	}
}

type Request struct {
	ConversationId string `json:"conversationId"`
	AtUsers        []struct {
		DingTalkId string `json:"dingtalkId"`
	} `json:"atUsers"`
	ChatBotUserId             string `json:"chatbotUserId"`
	MsgId                     string `json:"msgId"`
	SenderNick                string `json:"senderNick"`
	IsAdmin                   bool   `json:"isAdmin"`
	SessionWebhookExpiredTime int64  `json:"sessionWebhookExpiredTime"`
	CreateAt                  int64  `json:"createAt"`
	ConversationType          string `json:"conversationType"`
	SenderId                  string `json:"senderId"`
	ConversationTitle         string `json:"conversationTitle"`
	IsInAtList                bool   `json:"isInAtList"`
	SessionWebhook            string `json:"sessionWebhook"`
	Text                      struct {
		Content string `json:"content"`
	} `json:"text"`
	RobotCode string `json:"robotCode"`
	MsgType   string `json:"msgtype"`
}
