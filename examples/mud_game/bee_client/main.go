package main

import (
	_ "github.com/gorilla/websocket/examples/mud_game/bee_client/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.Run()
}

