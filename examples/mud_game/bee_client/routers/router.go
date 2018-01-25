package routers

import (
	"github.com/gorilla/websocket/examples/mud_game/bee_client/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.AutoRouter(&controllers.MainController{})
}
