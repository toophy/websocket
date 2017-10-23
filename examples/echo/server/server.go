// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// 房间
type Room struct {
	// 场景
	// 玩家
	// 积分
	// 挂靠的玩法
}

// 场景
type Scene struct {
	// 地图对应的文件(json格式,有版本号)
	// 地图内存数据
	// 开始时间
	// 存活时间
	// 不同类型的刷新区(道具,宝物,积分)
}

type AccountData struct {
	Account string
	Pwd     string
	Mt      int
	C       *websocket.Conn
	PosX    float64
	PosY    float64
}

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

var upgrader = websocket.Upgrader{} // use default options

// EchoProto echo协议
type EchoProto struct {
	C    string                 `json:"c"`
	M    string                 `json:"m"`
	Data map[string]interface{} `json:"data"`
}

type TMsgFunc func(a *AccountData, mt int, data *EchoProto) bool

type MessageFunc struct {
	CM   string
	Proc TMsgFunc
}

var gMsgFuncs map[string]*MessageFunc

var gAccounts map[string]*AccountData

func echo(w http.ResponseWriter, r *http.Request) {

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	a := &AccountData{}
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	a.C = c

	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		var em EchoProto
		json.Unmarshal(message, &em)
		if _, ok := gMsgFuncs[em.C+"."+em.M]; ok {
			gMsgFuncs[em.C+"."+em.M].Proc(a, mt, &em)
		}
	}

	if len(a.Account) > 0 {
		Broadcast_Scene_PlayerLeave(a)
		delete(gAccounts, a.Account)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {
	gMsgFuncs = make(map[string]*MessageFunc, 0)
	gAccounts = make(map[string]*AccountData, 0)

	gMsgFuncs["Index.Login"] = &MessageFunc{CM: "Index.Login", Proc: Index_Login}
	gMsgFuncs["Scene.Skill"] = &MessageFunc{CM: "Scene.Skill", Proc: Scene_Skill}
	gMsgFuncs["Scene.PlayerPoint"] = &MessageFunc{CM: "Scene.PlayerPoint", Proc: Scene_PlayerPoint}

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

var homeTemplate = template.Must(template.New("").Parse(`
    <!DOCTYPE html>
    <head>
    <meta charset="utf-8">
    <script>  
    window.addEventListener("load", function(evt) {
    
        var output = document.getElementById("output");
        var input = document.getElementById("input");
        var ws;
    
        var print = function(message) {
            var d = document.createElement("div");
            d.innerHTML = message;
            output.appendChild(d);
        };
    
        document.getElementById("open").onclick = function(evt) {
            if (ws) {
                return false;
            }
            ws = new WebSocket("{{.}}");
            ws.onopen = function(evt) {
                print("OPEN");
            }
            ws.onclose = function(evt) {
                print("CLOSE");
                ws = null;
            }
            ws.onmessage = function(evt) {
                print("RESPONSE: " + evt.data);
            }
            ws.onerror = function(evt) {
                print("ERROR: " + evt.data);
            }
            return false;
        };
    
        document.getElementById("send").onclick = function(evt) {
            if (!ws) {
                return false;
            }
            print("SEND: " + input.value);
            ws.send(input.value);
            return false;
        };
    
        document.getElementById("close").onclick = function(evt) {
            if (!ws) {
                return false;
            }
            ws.close();
            return false;
        };
    
    });
    </script>
    </head>
    <body>
    <table>
    <tr><td valign="top" width="50%">
    <p>Click "Open" to create a connection to the server, 
    "Send" to send a message to the server and "Close" to close the connection. 
    You can change the message and send multiple times.
    <p>
    <form>
    <button id="open">Open</button>
    <button id="close">Close</button>
    <p><input id="input" type="text" value="Hello world!">
    <button id="send">Send</button>
    </form>
    </td><td valign="top" width="50%">
    <div id="output"></div>
    </td></tr></table>
    </body>
    </html>
    `))
