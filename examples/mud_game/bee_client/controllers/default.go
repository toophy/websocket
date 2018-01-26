package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket/examples/mud_game/bee_client/models"
)

type MainController struct {
	beego.Controller
}

var (
	gAccountControllers map[string]*MainController
)

func init(){
	gAccountControllers = make(map[string]*MainController,0)
}

func GetAccountController(name string) *MainController{
	if v,ok:=gAccountControllers[name];ok{
		return v
	}
	return nil
}

func (c* MainController) AccLogin(){
	uid:=c.Input().Get("uid")
	pwd:=c.Input().Get("pwd")
	gAccountControllers[uid] = c
	models.AccountLogin(uid, pwd)
	a:=models.GetAccount(uid)
	
	if a!=nil{
		ret := <- a.C1
		c.Ctx.WriteString(ret)
	} else {
		ret := "帐号登录失败"
		c.Ctx.WriteString(ret)
		return 
	}
	
	return
}

func (c*MainController) AccLeave(){
	uid:=c.Input().Get("uid")
	models.AccountLeave(uid)
	c.Ctx.WriteString("Leave ok")
}