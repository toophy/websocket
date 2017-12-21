package mud

import (
	"time"
)

// Hall 大厅
type Hall struct {
	Accounts map[string]*AccountReal // 本大厅登录帐号
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

}
