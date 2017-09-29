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

type AccountData struct {
	Account string
	Pwd     string
	Mt      int
	C       *websocket.Conn
}

var addr = flag.String("addr", "localhost:8080", "http service address")

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

	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

func Broadcast_Scene_Skill(a *AccountData, data *EchoProto) {

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		// 技能目标
		// 技能xx
		retData := struct {
			Name    string `json:"name"` // 技能name
			Account string `json:"account"`
			IsDown  bool   `json:"isDown"`
			Ret     string `json:"ret"`
		}{
			data.Data["name"].(string),
			a.Account,
			data.Data["isDown"].(bool),
			"ok"}

		ret, _ := json.Marshal(struct {
			C    string `json:"c"`
			M    string `json:"m"`
			Data struct {
				Name    string `json:"name"` // 技能name
				Account string `json:"account"`
				IsDown  bool   `json:"isDown"`
				Ret     string `json:"ret"`
			} `json:"data"`
		}{
			C:    "Scene",
			M:    "Skill",
			Data: retData})

		err := gAccounts[k].C.WriteMessage(a.Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Scene_Skill(a *AccountData, mt int, data *EchoProto) bool {

	// 技能目标
	// 技能xx
	retData := struct {
		Name    string `json:"name"` // 技能name
		Account string `json:"account"`
		Ret     string `json:"ret"`
	}{
		data.Data["name"].(string),
		a.Account,
		"ok"}

	ret, _ := json.Marshal(struct {
		C    string `json:"c"`
		M    string `json:"m"`
		Data struct {
			Name    string `json:"name"` // 技能name
			Account string `json:"account"`
			Ret     string `json:"ret"`
		} `json:"data"`
	}{
		C:    data.C,
		M:    data.M,
		Data: retData})

	err := a.C.WriteMessage(mt, ret)
	if err != nil {
		log.Println("write:", err)
		return false
	}

	Broadcast_Scene_Skill(a, data)

	return true
}

func Broadcast_Scene_PlayerEnter(a *AccountData, data *EchoProto) {

	data.C = "Scene"
	data.M = "PlayerEnter"

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		pos_x, _ := data.Data["pos_x"].(int)
		pos_y, _ := data.Data["pos_y"].(int)

		retData := struct {
			Account string `json:"account"`
			Pwd     string `json:"pwd"`
			PosX    int    `json:"pos_x"`
			PosY    int    `json:"pos_y"`
			Ret     string `json:"ret"`
			Msg     string `json:"msg"`
		}{
			data.Data["account"].(string),
			data.Data["pwd"].(string),
			pos_x,
			pos_y,
			"ok",
			""}

		ret, _ := json.Marshal(struct {
			C    string `json:"c"`
			M    string `json:"m"`
			Data struct {
				Account string `json:"account"`
				Pwd     string `json:"pwd"`
				PosX    int    `json:"pos_x"`
				PosY    int    `json:"pos_y"`
				Ret     string `json:"ret"`
				Msg     string `json:"msg"`
			} `json:"data"`
		}{
			C:    data.C,
			M:    data.M,
			Data: retData})

		err := gAccounts[k].C.WriteMessage(gAccounts[k].Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Broadcast_Scene_PlayerLeave(a *AccountData) {

	for k, _ := range gAccounts {
		if gAccounts[k].Account == a.Account {
			continue
		}

		retData := struct {
			Account string `json:"account"`
			Ret     string `json:"ret"`
			Msg     string `json:"msg"`
		}{
			a.Account,
			"ok",
			""}

		ret, _ := json.Marshal(struct {
			C    string `json:"c"`
			M    string `json:"m"`
			Data struct {
				Account string `json:"account"`
				Ret     string `json:"ret"`
				Msg     string `json:"msg"`
			} `json:"data"`
		}{
			C:    "Scene",
			M:    "PlayerLeave",
			Data: retData})

		err := gAccounts[k].C.WriteMessage(gAccounts[k].Mt, ret)
		if err != nil {
			log.Println("write:", err)
		}
	}
}

func Index_Login(a *AccountData, mt int, data *EchoProto) bool {

	pos_x, _ := data.Data["pos_x"].(int)
	pos_y, _ := data.Data["pos_y"].(int)

	retData := struct {
		Account string `json:"account"`
		Pwd     string `json:"pwd"`
		PosX    int    `json:"pos_x"`
		PosY    int    `json:"pos_y"`
		Ret     string `json:"ret"`
		Msg     string `json:"msg"`
	}{
		data.Data["account"].(string),
		data.Data["pwd"].(string),
		pos_x,
		pos_y,
		"ok",
		""}

	log.Printf("account:%s,pwd:%s", data.Data["account"], data.Data["pwd"])
	if _, ok := gAccounts[data.Data["account"].(string)]; ok {
		retData.Ret = "fail"
		retData.Msg = "重复登录"
	} else {
		a.Account = data.Data["account"].(string)
		a.Pwd = data.Data["pwd"].(string)
		a.Mt = mt
		gAccounts[data.Data["account"].(string)] = a
	}

	ret, _ := json.Marshal(struct {
		C    string `json:"c"`
		M    string `json:"m"`
		Data struct {
			Account string `json:"account"`
			Pwd     string `json:"pwd"`
			PosX    int    `json:"pos_x"`
			PosY    int    `json:"pos_y"`
			Ret     string `json:"ret"`
			Msg     string `json:"msg"`
		} `json:"data"`
	}{
		C:    data.C,
		M:    data.M,
		Data: retData})

	err := a.C.WriteMessage(mt, ret)
	if err != nil {
		log.Println("write:", err)
		return false
	}

	data.C = "Scene"
	data.M = "PlayerEnter"
	Broadcast_Scene_PlayerEnter(a, data)

	return true
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
