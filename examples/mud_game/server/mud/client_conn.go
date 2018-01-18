package mud

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/gorilla/websocket/examples/echo/server/game"
	"github.com/gorilla/websocket/examples/echo/server/mud"
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
	gMsgFuncs     map[string]*MessageFunc
	gAccountConns map[string]*AccountConn
)

func init() {
	gMsgFuncs = make(map[string]*MessageFunc, 0)
	gAccountConns = make(map[string]*AccountConn, 0)

	gMsgFuncs["Index.Login"] = &MessageFunc{CM: "Index.Login", Proc: OnCMsg_AccountLogin}
	gMsgFuncs["Index.AskMatch"] = &MessageFunc{CM: "Index.AskMatch", Proc: OnCMsg_AskMatch}
	gMsgFuncs["Index.SendMail"] = &MessageFunc{CM: "Index.SendMail", Proc: OnCMsg_SendMail}
}

func ClientNetConn(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	a := &AccountConn{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		println("upgrade:", err)
		return
	}
	a.C = c

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			println("read:", err)
			break
		}
		println("recv: ", message)
		var em EchoProto
		json.Unmarshal(message, &em)
		if _, ok := gMsgFuncs[em.C+"."+em.M]; ok {
			gMsgFuncs[em.C+"."+em.M].Proc(a, mt, &em)
		}
	}

	if len(a.Account) > 0 {
		delete(gAccountConns, a.Account)
	}
}
