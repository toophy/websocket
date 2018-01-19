package mud

import (
	"time"
)

// Hall 大厅
type Hall struct {
	Accounts   map[string]*AccountReal // 本大厅登录帐号
	AccountIDs map[int64]*AccountReal  // 本大厅登录帐号ID
}

var (
	hall *Hall
)

func init() {
	hall = &Hall{}
	hall.Accounts = make(map[string]*AccountReal, 0)
	hall.AccountIDs = make(map[int64]*AccountReal, 0)
}

func GetHall() *Hall {
	return hall
}

// Update
func (h *Hall) Update() {

}

// AccountLogin 登录
func (h *Hall) AccountLogin(name string) {
	if _, ok := h.Accounts[name]; !ok {
		go h.ToGetAccount(name)
	} else {
		h.Accounts[name].LoadOver = true
		h.Accounts[name].Online = true
		h.Accounts[name].LastTime = int32(time.Now().Unix())
		h.Accounts[name].LastMailID = 0
		h.Accounts[name].LastGetMailTime = time.Now().UnixNano()
	}
}

// ToGetAccount 验证帐号注册信息
func (h *Hall) ToGetAccount(name string) {
	if a, ok := GetDBS().GetAccount(name); ok {
		h.Accounts[name] = &AccountReal{
			AccountInfo:     a,
			LoadOver:        true,
			Online:          true,
			LastTime:        int32(time.Now().Unix()),
			LastMailID:      0,
			LastGetMailTime: time.Now().UnixNano()}

		h.AccountIDs[a.ID] = h.Accounts[name]
		println("AccountLogin Ok")
	} else {
		println("AccountLogin Failed")
	}
}

// AskMatch 请求匹配一种游戏方式
func (h *Hall) AskMatch(name string, game string) {
	if _, ok := h.Accounts[name]; ok {
		go GetMatchSys().AskMatch(name, game, h.Accounts[name].Step, h.Accounts[name].Elo, float64(h.Accounts[name].WinRate), float64(h.Accounts[name].Kda))
	}
}

// 匹配服,成功匹配
func (h *Hall) OnMatchOver(accounts []string) {
	for i := 0; i < len(accounts); i++ {
		if _, ok := h.Accounts[accounts[i]]; !ok {
			println("帐号", accounts[i], " 离线")
			break
		}
		println("匹配成功", accounts)
	}

	// 生成房间(Room), 开搞喽
}

// 房间战斗结束, 返回战斗结果, 每个玩家的信息分开
func (h *Hall) OnRoomOver(battles BattleInfo) {
	if _, ok := h.AccountIDs[battles.AccID]; ok {
		println("战报:", h.AccountIDs[battles.AccID].Name)
	} else {
		println("帐号不存在:", battles.AccID)
	}
}

// 处理返回邮件
func (h *Hall) OnRecvMails(accID int64, mails []Mail) {
	if _, ok := h.AccountIDs[accID]; ok {
		lastMailID := int32(0)
		for k, _ := range mails {
			if mails[k].ID > lastMailID {
				lastMailID = mails[k].ID
			}
			// 1. 系统的脚本邮件, a. 给客户端解释执行 b. 立即解释执行
			// 2. 普通邮件, a. 给客户端解释执行
			println("接收到邮件:", mails[k].Title)
		}
		if lastMailID > h.AccountIDs[accID].LastMailID {
			h.AccountIDs[accID].LastMailID = lastMailID
		}
	}
}
