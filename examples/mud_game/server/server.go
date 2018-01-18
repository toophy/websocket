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
	"github.com/gorilla/websocket/examples/echo/server/game"
	"github.com/gorilla/websocket/examples/echo/server/mud"
)

var addr = flag.String("addr", "0.0.0.0:8080", "http service address")

func home(w http.ResponseWriter, r *http.Request) {
	homeTemplate.Execute(w, "ws://"+r.Host+"/echo")
}

func main() {

	flag.Parse()
	log.SetFlags(0)

	mud.GetDBS().AccountRegist("wind1", "蓝枫1号")
	mud.GetDBS().AccountRegist("wind2", "蓝枫2号")

	go mud.GetHall().Update()
	go mud.GetMailSys().Update()

	http.HandleFunc("/echo", mud.ClientNetConn)
	http.HandleFunc("/", home)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
