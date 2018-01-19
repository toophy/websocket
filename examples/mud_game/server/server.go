// Copyright 2015 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"flag"

	"log"
	"net/http"

	"github.com/gorilla/websocket/examples/mud_game/server/mud"
)

var addr = flag.String("addr", "0.0.0.0:1888", "http service address")

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
