package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket/examples/mud_game/bee_client/models"
)

type MainController struct {
	beego.Controller
}

func (c* MainController) AccLogin(){
	models.AccountLogin("wind1", "123")
	c.Ctx.WriteString("呵呵")
}
