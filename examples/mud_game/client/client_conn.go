package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	// "encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// EchoProto echo协议
type EchoProto struct {
	C    string                 `json:"c"`
	M    string                 `json:"m"`
	Data map[string]interface{} `json:"data"`
}

type TMsgFunc func(a *AccountConn, mt int, data *EchoProto) bool

type MessageFunc struct {
	CM   string
	Proc TMsgFunc
}

var (
	gMsgFuncs map[string]*MessageFunc
)

func init() {
	gMsgFuncs = make(map[string]*MessageFunc, 0)

	gMsgFuncs["Index.Login"] = &MessageFunc{CM: "Index.Login", Proc: OnCMsgAccountLogin}
	gMsgFuncs["Index.AskMatch"] = &MessageFunc{CM: "Index.AskMatch", Proc: OnCMsgAskMatch}
	gMsgFuncs["Index.SendMail"] = &MessageFunc{CM: "Index.SendMail", Proc: OnCMsgSendMail}
}

// ClientNetConn 处理客户端网络连接消息
func ClientNetConn(w http.ResponseWriter, r *http.Request) {

	// a := &AccountConn{Temp: true}
	// c, err := upgrader.Upgrade(w, r, nil)
	// if err != nil {
	// 	fmt.Println("[E] websocket升级错误:", err)
	// 	return
	// }
	// a.C = c

	// defer c.Close()
	// for {
	// 	mt, message, err := c.ReadMessage()
	// 	if err != nil {
	// 		fmt.Println("[E] 网络连接读取错误:", err)
	// 		break
	// 	}
	// 	var em EchoProto
	// 	json.Unmarshal(message, &em)
	// 	if _, ok := gMsgFuncs[em.C+"."+em.M]; ok {
	// 		gMsgFuncs[em.C+"."+em.M].Proc(a, mt, &em)
	// 	}
	// }
}

func AccountLogin(account string, pwd string) {
	retData := EchoProto{
		"Index",
		"Login",
		map[string]interface{}{}}

	retData.Data["account"] = account
	retData.Data["pwd"] = pwd

	if GetAccount(account) != nil {
		fmt.Printf("[I] 帐号%s已经在登录中\n", account)
		return
	}

	// newAcc := NewAccount(account)

	u := url.URL{Scheme: "ws", Host: "localhost:1888", Path: "/echo"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}

	go func() {
		defer c.Close()
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				log.Println("read:", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ret, _ := json.Marshal(retData)

	err = c.WriteMessage(websocket.TextMessage, ret)
	if err != nil {
		log.Println("write:", err)
		return
	}
}
