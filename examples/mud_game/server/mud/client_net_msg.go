package mud

func OnCMsgAccountLeave(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().AccountLeave(a, mt, data)
	return true
}

func OnCMsgAccountLogin(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().AccountLogin(a, mt, data)
	return true
}

func OnCMsgAskMatch(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().AskMatch(a, mt, data)
	return true
}

func OnCMsgSendMail(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().SendMail(a, mt, data)
	return true
}

func OnCMsgGetMails(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().GetMails(a, mt, data)
	return true
}
