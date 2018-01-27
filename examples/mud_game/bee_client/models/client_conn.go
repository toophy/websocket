package models

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
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
	GMsgFuncs map[string]*MessageFunc
)

func init() {
	GMsgFuncs = make(map[string]*MessageFunc, 0)
}

func AccountLeave(account string) {
	newAcc := GetAccount(account)
	if newAcc != nil {
		if newAcc.C != nil {
			go func() {
				retData := EchoProto{
					"Index",
					"Leave",
					map[string]interface{}{}}

				retData.Data["account"] = account

				ret, _ := json.Marshal(retData)

				err := newAcc.C.WriteMessage(websocket.TextMessage, ret)
				if err != nil {
					log.Println("write:", err)
					return
				}
			}()
		} else {
			LeaveAccount(account)
		}
	}
}

func AccountLogin(account string, pwd string) {
	retData := EchoProto{
		"Index",
		"Login",
		map[string]interface{}{}}

	retData.Data["account"] = account
	retData.Data["pwd"] = pwd

	newAcc := GetAccount(account)
	if newAcc != nil {
		if newAcc.C != nil {
			go func() {
				newAcc.C1 <- fmt.Sprintf("[I] 帐号%s已经在登录中", account)
			}()
		}
		return
	} else {
		newAcc = NewAccount(account)
		if newAcc == nil {
			return
		}
	}

	u := url.URL{Scheme: "ws", Host: "localhost:1888", Path: "/echo"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		go func() {
			newAcc.C1 <- fmt.Sprintf("[I] 帐号%s连接错误: %s", account, err.Error())
		}()
		return
	}

	newAcc.C = c

	go func() {
		defer c.Close()
		for {
			mt, message, err := c.ReadMessage()
			if err != nil {
				newAcc.C1 <- fmt.Sprintf("[I] 帐号%s网络连接读取错误: %s", account, err.Error())
				return
			}
			log.Printf("recv: %s", message)
			a := GetAccount(account)
			var em EchoProto
			json.Unmarshal(message, &em)
			if _, ok := GMsgFuncs[em.C+"."+em.M]; ok {
				GMsgFuncs[em.C+"."+em.M].Proc(a, mt, &em)
			}
		}
	}()

	ret, _ := json.Marshal(retData)

	err = c.WriteMessage(websocket.TextMessage, ret)
	if err != nil {
		log.Println("write:", err)
		return
	}
}

func SendMail(account, recer, title, content string) {

	newAcc := GetAccount(account)
	if newAcc == nil || newAcc.C == nil {
		go func() {
			newAcc.C1 <- fmt.Sprintf("[I] 帐号%s不在登录中", account)
		}()
		return
	}

	if ret, ok := json.Marshal(&EchoProto{
		C: "Index",
		M: "SendMail",
		Data: map[string]interface{}{
			"recer":   recer,
			"title":   title,
			"content": content}}); ok == nil {
		if err := newAcc.C.WriteMessage(websocket.TextMessage, ret); err != nil {
			go func(retErr string) {
				newAcc.C1 <- fmt.Sprintf("[I] 帐号%s发送邮件失败:%s", account, retErr)
			}(err.Error())
		}
	}
}

func GetMails(account string) {

	newAcc := GetAccount(account)
	if newAcc == nil || newAcc.C == nil {
		go func() {
			newAcc.C1 <- fmt.Sprintf("[I] 帐号%s不在登录中", account)
		}()
		return
	}

	if ret, ok := json.Marshal(&EchoProto{
		C:    "Index",
		M:    "GetMails",
		Data: map[string]interface{}{}}); ok == nil {
		if err := newAcc.C.WriteMessage(websocket.TextMessage, ret); err != nil {
			go func(retErr string) {
				newAcc.C1 <- fmt.Sprintf("[I] 帐号%s收取邮件失败:%s", account, retErr)
			}(err.Error())
		}
	}
}
