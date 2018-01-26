package main

import (
	_ "github.com/gorilla/websocket/examples/mud_game/bee_client/routers"
	_ "github.com/gorilla/websocket/examples/mud_game/bee_client/models/net_msg"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

