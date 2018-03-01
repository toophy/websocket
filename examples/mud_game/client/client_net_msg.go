package main

func OnCMsgAccountLogin(a *AccountConn, mt int, data *EchoProto) bool {
	return true
}

func OnCMsgAskMatch(a *AccountConn, mt int, data *EchoProto) bool {
	return true
}

func OnCMsgSendMail(a *AccountConn, mt int, data *EchoProto) bool {
	return true
}
