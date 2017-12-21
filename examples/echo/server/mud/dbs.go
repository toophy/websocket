package mud

import (
	"time"
)

// DBS 虚拟数据库
type DBS struct {
	Accounts map[string]*AccountInfo
}

var (
	dbs *DBS
)

func init() {
	dbs = &DBS{}
}

// GetDBS 获取数据库
func GetDBS() *DBS {
	return dbs
}

// AccountRegist 注册
func (d *DBS) AccountRegist(name string, nick string) bool {
	if _, ok := d.Accounts[name]; !ok {
		d.Accounts[name] = &AccountInfo{
			ID:         0,
			Name:       name,
			Nick:       nick,
			Step:       0,
			WinRate:    0,
			RegistTime: int32(time.Now().Unix())}
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
