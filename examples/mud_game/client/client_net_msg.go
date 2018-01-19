package main

func OnCMsg_AccountLogin(a *AccountConn, mt int, data *EchoProto) bool {
	return true
}

func OnCMsg_AskMatch(a *AccountConn, mt int, data *EchoProto) bool {
	return true
}

func OnCMsg_SendMail(a *AccountConn, mt int, data *EchoProto) bool {
	return true
}
