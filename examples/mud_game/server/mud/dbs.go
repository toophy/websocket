package mud

import (
	"time"
)

// DBS 虚拟数据库
type DBS struct {
	Accounts   map[string]*AccountInfo
	AccountIDs map[int64]*AccountInfo
	LastID     int64
}

var (
	dbs *DBS
)

func init() {
	dbs = &DBS{
		Accounts:   make(map[string]*AccountInfo, 0),
		AccountIDs: make(map[int64]*AccountInfo, 0),
		LastID:     int64(1)}
}

// GetDBS 获取数据库
func GetDBS() *DBS {
	return dbs
}

// AccountRegist 注册
func (d *DBS) AccountRegist(name string, nick string) bool {
	if _, ok := d.Accounts[name]; !ok {
		d.Accounts[name] = &AccountInfo{
			ID:         d.LastID,
			Name:       name,
			Nick:       nick,
			Step:       0,
			WinRate:    0,
			RegistTime: int32(time.Now().Unix())}

		d.AccountIDs[d.LastID] = d.Accounts[name]
		go GetMailSys().RegistMailBox(d.Accounts[name].ID, name)

		d.LastID++
		return true
	}
	return false
}

// GetAccount 获取帐号信息
func (d *DBS) GetAccount(name string) (a AccountInfo, ret bool) {
	ret = false
	if _, ok := d.Accounts[name]; ok {
		a = *d.Accounts[name]
		ret = true
	}
	return
}

// GetAccount 获取帐号信息
func (d *DBS) GetAccountByID(accID int64) (a AccountInfo, ret bool) {
	ret = false
	if _, ok := d.AccountIDs[accID]; ok {
		a = *d.AccountIDs[accID]
		ret = true
	}
	return
}
