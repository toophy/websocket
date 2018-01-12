package mud

// 邮件
// 邮件类型 : 脚本邮件
// 邮件标题 : 拍卖场竞价成功
// 文本内容 : 恭喜你在xxx拍卖场xxx次拍卖会,使用xx通货成功竞价获取xx钻石
//            收到邮件后, 扣除xx通货xx个, 获取钻石xx个
// 脚本内容 : AutionSuccess(xx钻石数量,xx通货ID,xx通货数量)

// Mail 邮件
type Mail struct {
	ID            int64  // ID
	RSSMailBoxID  int64  // 订阅的邮箱
	RSSID         int64  // 订阅的邮件编号
	Title         string // 标题
	SenderType    string // 发送人类型: 玩家, 系统, 战队
	SenderName    string // 发送人名称: SYSTEM,TEAM,Account
	RecverType    string // 接收人类型: 玩家, 系统, 战队
	RecverName    string // 接收人名称: SYSTEM,TEAM,Account
	Content       string // 文本内容(还有客户端脚本)
	Script        string // 脚本内容: 供服务端解析执行
	Read          bool   // 是否已读
	RecvAccessory bool   // 是否接收了附件
	SendTime      int64  // 发送时间
	DestoryTime   int64  // 保存截止时间
}

type MailSys struct {
	Mails  map[int64]*Mail
	LastID int64
}

// Mailbox 邮箱
type Mailbox struct {
	ID     int64           // 邮箱编号
	Mails  map[int64]*Mail // 邮件
	LastID int64           // 最后一封邮件编号
	RSSBox map[int64]int64 // 订阅邮箱, [订阅的邮箱编号]最后一封邮件编号
}

var (
	mailSys *MailSys
)

func init() {
	mailSys = &MailSys{}
	mailSys.Mails = make(map[int64]*Mail, 0)
	MailSys.LastID = 1
}

// GetMailSys 获取邮件系统
func GetMailSys() *MailSys {
	return mailSys
}

// Send 发送邮件
func (m *MailSys) Send(mail *Mail) {
	mail.ID = mailSys.LastID
	mailSys.LastID++
	m.Mails[mail.ID] = mail
}

// GetMails 索取所有邮件(可以考虑分段索取)
func (m *MailSys) GetMails(recverType string, recverName string) (ret []Mail) {
	for k, v := range m.Mails {
		if v.RecverType == recverType && v.RecverName == recverName {
			ret = append(ret, *v)
		}
	}
	return
}

// SetMailRead 设置邮件已经读取的标记
func (m *MailSys) SetMailRead(mailID int64) {
	if _, ok := m.Mails[mailID]; ok {
		m.Mails[mailID].Read = true
	}
}

// SetMailRecvAccessory 设置附件已经索取的标记
func (m *MailSys) SetMailRecvAccessory(mailID int64) {
	if _, ok := m.Mails[mailID]; ok {
		m.Mails[mailID].RecvAccessory = true
	}
}

// DestoryMail 摧毁指定邮件(附件已经索取)
func (m *MailSys) DestoryMail(mailID int64) {
	if _, ok := m.Mails[mailID]; ok {
		if m.Mails[mailID].RecvAccessory {
			delete(m.Mails, mailID)
		}
	}
}

// DestoryMailByName 通过接收者名称,类型消耗附件已经索取的邮件
func (m *MailSys) DestoryMailByName(recverType string, recverName string) {
	ret := []int64{}

	for k, v := range m.Mails {
		if v.RecverType == recverType && v.RecverName == recverName {
			if v.RecvAccessory {
				ret = append(ret, v.ID)
			}
		}
	}

	for i := 0; i < len(ret); i++ {
		delete(m.Mails, ret[i])
	}
}
