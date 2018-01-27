package mud

import (
	"encoding/json"
	"sync"
	"time"
)

// 邮件
// 邮件类型 : 脚本邮件
// 邮件标题 : 拍卖场竞价成功
// 文本内容 : 恭喜你在xxx拍卖场xxx次拍卖会,使用xx通货成功竞价获取xx钻石
//            收到邮件后, 扣除xx通货xx个, 获取钻石xx个
// 脚本内容 : AutionSuccess(xx钻石数量,xx通货ID,xx通货数量)

// Mail 邮件
type Mail struct {
	ID            int32  // ID
	RSSMailBoxID  int64  // 订阅的邮箱
	RSSID         int32  // 订阅的邮件编号
	Title         string // 标题
	SenderID      int64  // 发送者ID
	RecverID      int64  // 接收者ID
	Content       string // 文本内容(还有客户端脚本)
	Script        string // 脚本内容: 供服务端解析执行
	Read          bool   // 是否已读
	RecvAccessory bool   // 是否接收了附件
	SendTime      int32  // 发送时间
	DestoryTime   int32  // 保存截止时间
}

// Mailbox 邮箱
type Mailbox struct {
	ID      int64           // 邮箱编号
	Account string          // 帐号名
	Mails   map[int32]*Mail // 邮件
	LastID  int32           // 最后一封邮件编号
	RSSBox  map[int64]int32 // 订阅邮箱, [订阅的邮箱编号]最后一封邮件编号
}

// 邮件系统
type MailSys struct {
	mailboxs     map[int64]*Mailbox    // 邮箱
	mailChange   map[int64]interface{} // 邮件有变动
	mailAccounts map[string]int64      // 帐号名转邮箱ID
	lastID       int64                 // 最后一个邮箱编号
	skey         int32                 // 区服编号
	locker       *sync.Mutex           // 邮件锁
}

var (
	mailSys *MailSys
)

func init() {
	mailSys = &MailSys{}
	mailSys.mailboxs = make(map[int64]*Mailbox, 0)
	mailSys.mailChange = make(map[int64]interface{}, 0)
	mailSys.mailAccounts = make(map[string]int64, 0)
	mailSys.lastID = 1
	mailSys.skey = 1
	mailSys.locker = new(sync.Mutex)
}

// GetMailSys 获取邮件系统
func GetMailSys() *MailSys {
	return mailSys
}

func (m *MailSys) Update() {
	// 订阅
}

// Send 发送邮件
func (m *MailSys) SendByName(senderID int64, recerName string, title, content, script string, back func(int64, string, string), tempID string) {
	if recerID, ok := m.GetMailIDByAccount(recerName); ok {
		m.Send(senderID, recerID, title, content, script, back, tempID)
	} else {
		go back(senderID, tempID, "no recer")
	}
}

// Send 发送邮件
func (m *MailSys) Send(senderID, recerID int64, title, content, script string, back func(int64, string, string), tempID string) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if v, ok := m.mailboxs[recerID]; ok {
		mailBody := &Mail{
			ID:            v.LastID,
			RSSMailBoxID:  0,
			RSSID:         0,
			Title:         title,
			SenderID:      senderID,
			RecverID:      recerID,
			Content:       content,
			Script:        script,
			Read:          false,
			RecvAccessory: false,
			SendTime:      int32(time.Now().Unix()),
			DestoryTime:   int32(time.Now().Unix()) + int32(1*24*3600)}
		v.Mails[mailBody.ID] = mailBody
		v.LastID++

		m.mailChange[recerID] = true

		go back(senderID, tempID, "ok")
	} else {
		go back(senderID, tempID, "no recer")
	}
}

// Send 发送邮件
func (m *MailSys) SendRSS(rSSID int32, rSSMailBoxID, senderID, recerID int64, title, content, script string) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if v, ok := m.mailboxs[recerID]; ok {
		mailBody := &Mail{
			ID:            v.LastID,
			RSSMailBoxID:  rSSMailBoxID,
			RSSID:         rSSID,
			Title:         title,
			SenderID:      senderID,
			RecverID:      recerID,
			Content:       content,
			Script:        script,
			Read:          false,
			RecvAccessory: false,
			SendTime:      int32(time.Now().Unix()),
			DestoryTime:   int32(time.Now().Unix()) + int32(1*24*3600)}
		v.Mails[mailBody.ID] = mailBody
		v.LastID++

		// 更新订阅版本号
		if _, ok := v.RSSBox[rSSMailBoxID]; ok {
			if rSSID > v.RSSBox[rSSMailBoxID] {
				v.RSSBox[rSSMailBoxID] = rSSID
			}
		}

		m.mailChange[recerID] = true
	}
}

// RegistMailBox 注册邮箱
func (m *MailSys) RegistMailBox(accID int64, accName string) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		return
	}

	if len(accName) < 1 {
		return
	}

	if _, ok := m.mailAccounts[accName]; ok {
		return
	}

	m.mailboxs[accID] = &Mailbox{
		ID:      accID,
		Account: accName,
		Mails:   make(map[int32]*Mail, 0),
		LastID:  1,
		RSSBox:  make(map[int64]int32)}

	m.mailAccounts[accName] = accID
}

// GetMailIDByAccount 通过帐号名索取邮箱ID
func (m *MailSys) GetMailIDByAccount(accName string) (int64, bool) {
	m.locker.Lock()
	defer m.locker.Unlock()
	if v, ok := m.mailAccounts[accName]; ok {
		return v, true
	}
	return 0, false
}

// GetMails 索取所有邮件(可以考虑分段索取)
func (m *MailSys) GetMails(accID int64, back func(int64, string, string)) (ret []Mail) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		for _, v := range m.mailboxs[accID].Mails {
			ret = append(ret, *v)
		}
		mails, _ := json.Marshal(&ret)
		go back(accID, "ok", string(mails))
	} else {
		go back(accID, "no mailbox", "")
	}
	return
}

// GetNextMails 索取mailID之后的所有邮件
func (m *MailSys) GetNextMails(accID int64, mailID int32) (ret []Mail) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		for _, v := range m.mailboxs[accID].Mails {
			if v.ID > mailID {
				ret = append(ret, *v)
			}
		}
	}

	go GetHall().OnRecvMails(accID, ret)
	return
}

// GetLastMailID 获取邮箱的最后一封信编号
func (m *MailSys) GetLastMailID(accID int64) int32 {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		if m.mailboxs[accID].LastID > 0 {
			return m.mailboxs[accID].LastID
		}
	}
	return 0
}

// AppendRSS 增加订阅
func (m *MailSys) AppendRSS(accID int64, rssID int64, rssLastID int32) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		m.mailboxs[accID].RSSBox[rssID] = rssLastID
	}
}

// RemoveRSS 删除订阅
func (m *MailSys) RemoveRSS(accID int64, rssID int64) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		delete(m.mailboxs[accID].RSSBox, rssID)
	}
}

// SetMailRead 设置邮件已经读取的标记
func (m *MailSys) SetMailRead(accID int64, mailID int32) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		if _, ok := m.mailboxs[accID].Mails[mailID]; ok {
			m.mailboxs[accID].Mails[mailID].Read = true
		}
	}
}

// SetMailRecvAccessory 设置附件已经索取的标记
func (m *MailSys) SetMailRecvAccessory(accID int64, mailID int32) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		if _, ok := m.mailboxs[accID].Mails[mailID]; ok {
			m.mailboxs[accID].Mails[mailID].RecvAccessory = true
		}
	}
}

// DestoryMail 摧毁指定邮件(附件已经索取)
func (m *MailSys) DestoryMail(accID int64, mailID int32) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		if _, ok := m.mailboxs[accID].Mails[mailID]; ok {
			if m.mailboxs[accID].Mails[mailID].RecvAccessory {
				delete(m.mailboxs[accID].Mails, mailID)
			}
		}
	}
}

// DestoryMailByName 通过接收者名称,类型消耗附件已经索取的邮件
func (m *MailSys) DestoryAccountMail(accID int64) {
	m.locker.Lock()
	defer m.locker.Unlock()

	if _, ok := m.mailboxs[accID]; ok {
		ret := []int32{}
		for k, v := range m.mailboxs[accID].Mails {
			if v.RecvAccessory {
				ret = append(ret, k)
			}
		}
		for i := 0; i < len(ret); i++ {
			delete(m.mailboxs[accID].Mails, ret[i])
		}
	}
}
