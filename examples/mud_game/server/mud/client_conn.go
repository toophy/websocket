package mud

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

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

	gMsgFuncs["Index.Login"] = &MessageFunc{CM: "Index.Login", Proc: OnCMsg_AccountLogin}
	gMsgFuncs["Index.AskMatch"] = &MessageFunc{CM: "Index.AskMatch", Proc: OnCMsg_AskMatch}
	gMsgFuncs["Index.SendMail"] = &MessageFunc{CM: "Index.SendMail", Proc: OnCMsg_SendMail}
	gMsgFuncs["Index.Leave"] = &MessageFunc{CM: "Index.Leave", Proc: OnCMsg_AccountLeave}
}

// ClientNetConn 处理客户端网络连接消息
func ClientNetConn(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	a := &AccountConn{Temp: true}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("[E] websocket升级错误:", err)
		return
	}
	a.C = c

	defer c.Close()
	for {
		var em EchoProto
		mt, message, err := c.ReadMessage()
		if err != nil {
			em.Data = make(map[string]interface{},0)
			em.Data["account"] = a.Account
			GetHall().AccountLeave(a, mt, &em)
			fmt.Println("[E] 网络连接读取错误:", err)
			break
		}
		
		json.Unmarshal(message, &em)
		if _, ok := gMsgFuncs[em.C+"."+em.M]; ok {
			gMsgFuncs[em.C+"."+em.M].Proc(a, mt, &em)
		}
	}
}
