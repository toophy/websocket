package controllers

import (
	"github.com/astaxie/beego"
	"github.com/gorilla/websocket/examples/mud_game/bee_client/models"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) AccLogin() {
	uid := c.Input().Get("uid")
	pwd := c.Input().Get("pwd")
	models.AccountLogin(uid, pwd)
	a := models.GetAccount(uid)

	if a != nil {
		ret := <-a.C1
		c.Ctx.WriteString(ret)
	} else {
		ret := "帐号登录失败"
		c.Ctx.WriteString(ret)
		return
	}

	return
}

func (c *MainController) AccLeave() {
	uid := c.Input().Get("uid")
	models.AccountLeave(uid)
	c.Ctx.WriteString("Leave ok")
}

func (c *MainController) SendMail() {
	uid := c.Input().Get("uid")
	recer := c.Input().Get("recer")
	title := c.Input().Get("title")
	content := c.Input().Get("content")

	models.SendMail(uid, recer, title, content)
	a := models.GetAccount(uid)

	if a != nil {
		ret := <-a.C1
		c.Ctx.WriteString(ret)
	} else {
		ret := "帐号发送邮件失败"
		c.Ctx.WriteString(ret)
		return
	}

	return
}

func (c *MainController) GetMails() {
	uid := c.Input().Get("uid")

	models.GetMails(uid)
	a := models.GetAccount(uid)

	if a != nil {
		ret := <-a.C1
		c.Ctx.WriteString(ret)
	} else {
		ret := "信箱被你掏空了"
		c.Ctx.WriteString(ret)
		return
	}

	return
}
