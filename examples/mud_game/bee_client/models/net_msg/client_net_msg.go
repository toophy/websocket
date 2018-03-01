package net_msg

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket/examples/mud_game/bee_client/models"
)

func init() {
	models.GMsgFuncs["Index.Login"] = &models.MessageFunc{CM: "Index.Login", Proc: OnCMsgAccountLogin}
	models.GMsgFuncs["Index.AskMatch"] = &models.MessageFunc{CM: "Index.AskMatch", Proc: OnCMsgAskMatch}
	models.GMsgFuncs["Index.SendMail"] = &models.MessageFunc{CM: "Index.SendMail", Proc: OnCMsgSendMail}
	models.GMsgFuncs["Index.GetMails"] = &models.MessageFunc{CM: "Index.GetMails", Proc: OnCMsg_GetMails}
	models.GMsgFuncs["Index.Leave"] = &models.MessageFunc{CM: "Index.Leave", Proc: OnCMsgAccountLeave}
}

func OnCMsgAccountLeave(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	// 发送leave
	fmt.Printf("[I] 帐号%s正常离线", a.Account)
	models.LeaveAccount(a.Account)
	return true
}

func OnCMsgAccountLogin(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	rets, _ := json.Marshal(data)
	a.C1 <- string(rets)
	return true
}

func OnCMsgAskMatch(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	return true
}

func OnCMsgSendMail(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	rets, _ := json.Marshal(data)
	a.C1 <- string(rets)
	return true
}

func OnCMsg_GetMails(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	rets, _ := json.Marshal(data)
	a.C1 <- string(rets)
	return true
}
