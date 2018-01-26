package net_msg

import(
	"fmt"
	"encoding/json"
	"github.com/gorilla/websocket/examples/mud_game/bee_client/models"
)

func init(){
	 models.GMsgFuncs["Index.Login"] = &models.MessageFunc{CM: "Index.Login", Proc: OnCMsg_AccountLogin}
	 models.GMsgFuncs["Index.AskMatch"] = &models.MessageFunc{CM: "Index.AskMatch", Proc: OnCMsg_AskMatch}
	 models.GMsgFuncs["Index.SendMail"] = &models.MessageFunc{CM: "Index.SendMail", Proc: OnCMsg_SendMail}
	 models.GMsgFuncs["Index.Leave"] = &models.MessageFunc{CM: "Index.Leave", Proc: OnCMsg_AccountLeave}
}

func OnCMsg_AccountLeave(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	// 发送leave
	fmt.Printf("[I] 帐号%s正常离线", a.Account)
	models.LeaveAccount(a.Account)
	return true
}

func OnCMsg_AccountLogin(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	rets,_:=json.Marshal(data)
	a.C1 <- string(rets)
	return true
}

func OnCMsg_AskMatch(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	return true
}

func OnCMsg_SendMail(a *models.AccountConn, mt int, data *models.EchoProto) bool {
	return true
}
