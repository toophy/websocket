package mud

func OnCMsg_AccountLeave(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().AccountLeave(a,mt,data)
	return true
}

func OnCMsg_AccountLogin(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().AccountLogin(a, mt, data)
	return true
}

func OnCMsg_AskMatch(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().AskMatch(a, mt, data)
	return true
}

func OnCMsg_SendMail(a *AccountConn, mt int, data *EchoProto) bool {
	GetHall().SendMail(a, mt, data)
	return true
}
