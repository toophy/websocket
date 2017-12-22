package mud

import (
	"time"
)

// Hall 大厅
type Hall struct {
	Accounts map[string]*AccountReal // 本大厅登录帐号
}

var (
	hall *Hall
)

func init() {
	hall = &Hall{}
	hall.Accounts = make(map[string]*AccountReal, 0)
}

func GetHall() *Hall {
	return hall
}

// AccountLogin 登录
func (h *Hall) AccountLogin(name string) {
	if _, ok := h.Accounts[name]; !ok {
		go h.ToGetAccount(name)
	} else {
		h.Accounts[name].Online = true
		h.Accounts[name].LastTime = int32(time.Now().Unix())
	}
}

// ToGetAccount 验证帐号注册信息
func (h *Hall) ToGetAccount(name string) {
	if a, ok := GetDBS().GetAccount(name); ok {
		h.Accounts[name] = &AccountReal{
			AccountInfo: a,
			Online:      true,
			LastTime:    int32(time.Now().Unix())}
		println("AccountLogin Ok")
	} else {
		println("AccountLogin Failed")
	}
}

//
func (h *Hall) AskMatch(name string, game string) {
	if _, ok := h.Accounts[name]; ok {
		go GetMatchSys().AskMatch(name, game, h.Accounts[name].Step, h.Accounts[name].Elo, h.Accounts[name].WinRate, h.Accounts[name].Kda)
	}
}

// 匹配服,成功匹配
func (h *Hall) OnMatchOver(accounts []string) {
	for i := 0; i < len(accounts); i++ {
		if _, ok := h.Accounts[accounts[i]]; !ok {
			println("帐号", accounts[i], " 离线")
			break
		}
	}

	// 生成房间(Room), 开搞喽
}

// 房间战斗结束, 返回战斗结果, 每个玩家的信息分开
func (h *Hall) OnRoomOver(battles []BattleInfo) {
}
