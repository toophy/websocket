package mud

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"
)

// Hall 大厅
type Hall struct {
	Accounts      map[string]*AccountReal // 本大厅登录帐号
	AccountIDs    map[int64]*AccountReal  // 本大厅登录帐号ID
	AccountLocker *sync.Mutex             // 帐号列表锁
}

var (
	hall *Hall
)

func init() {
	hall = &Hall{}
	hall.Accounts = make(map[string]*AccountReal, 0)
	hall.AccountIDs = make(map[int64]*AccountReal, 0)
	hall.AccountLocker = new(sync.Mutex)
}

func GetHall() *Hall {
	return hall
}

// Update
func (h *Hall) Update() {
	t := time.NewTimer(5 * time.Second)
	for {
		select {
		case <-t.C:
			h.AccountLocker.Lock()
			for _, v := range h.Accounts {
				if !v.LoadOver {
					continue
				}
				if time.Now().Unix() > v.LastGetMailTime+10 {
					v.LastGetMailTime = time.Now().Unix()
					go GetMailSys().GetNextMails(v.AccountInfo.ID, v.LastMailID)
				}
			}
			h.AccountLocker.Unlock()
			t.Reset(1 * time.Second)
			break
		}
	}
}

// AccountLeave 登录
func (h *Hall) AccountLeave(a *AccountConn, mt int, data *EchoProto) {
	name := data.Data["account"].(string)

	fmt.Printf("[I] 帐号申请离线(%s)\n", name)

	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if _, ok := h.Accounts[name]; ok {
		if a.C!=nil {
			retData := *data
			retData.Data["ret"] = "ok"
			retData.Data["msg"] =  "正常离线"
			ret, _ := json.Marshal(retData)
			err := a.C.WriteMessage(mt, ret)
			if err != nil {
				fmt.Printf("[W] 向网络写消息失败:%s\n", err)
			}
			t:=time.Now()
			t.After(2)
			a.C.SetWriteDeadline(t)
		}
		
		delete(hall.AccountIDs,h.Accounts[name].ID)
		delete(hall.Accounts,name)
	} else {
		if a.C!=nil {
			retData := *data
			retData.Data["ret"] = "failed"
			retData.Data["msg"] =  "重复离线"
			ret, _ := json.Marshal(retData)
			err := a.C.WriteMessage(mt, ret)
			if err != nil {
				fmt.Printf("[W] 向网络写消息失败:%s\n", err)
			}
		}
	}
}

// AccountLogin 登录
func (h *Hall) AccountLogin(a *AccountConn, mt int, data *EchoProto) {
	name := data.Data["account"].(string)
	pwd := data.Data["pwd"].(string)

	a.Account = name

	fmt.Printf("[I] 帐号申请登录(%s : %s)\n", name, pwd)

	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if _, ok := h.Accounts[name]; !ok {
		go h.ToGetAccount(a, mt, data)
	} else {
		retData := *data
		retData.Data["ret"] = "fail"
		retData.Data["msg"] =  "重复登录"
		ret, _ := json.Marshal(retData)
		err := a.C.WriteMessage(mt, ret)
		if err != nil {
			fmt.Printf("[W] 向网络写消息失败:%s\n", err)
		}
	}
}

// ToGetAccount 验证帐号注册信息
func (h *Hall) ToGetAccount(ac *AccountConn, mt int, data *EchoProto) {
	name := data.Data["account"].(string)
	pwd := data.Data["pwd"].(string)

	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if a, ok := GetDBS().GetAccount(name); ok {
		h.Accounts[name] = &AccountReal{
			AccountInfo:     a,
			LoadOver:        true,
			Online:          true,
			LastTime:        int32(time.Now().Unix()),
			LastMailID:      0,
			LastGetMailTime: time.Now().UnixNano()}

		h.AccountIDs[a.ID] = h.Accounts[name]

		fmt.Printf("[I] 帐号成功登录(%s : %s)\n", name, pwd)

		retData := *data
		retData.Data["ret"] = "ok"
		retData.Data["msg"] =  ""

		ret, _ := json.Marshal(retData)
		err := ac.C.WriteMessage(mt, ret)
		if err != nil {
			fmt.Printf("[W] 向网络写消息失败:%s\n", err)
		}
	} else {
		fmt.Printf("[I] 帐号未注册:%s\n", name)

		retData := *data
		retData.Data["ret"] = "fail"
		retData.Data["msg"] =  "帐号未注册"
		ret, _ := json.Marshal(retData)
		err := ac.C.WriteMessage(mt, ret)
		if err != nil {
			fmt.Printf("[W] 向网络写消息失败:%s\n", err)
		}
	}
}

// AskMatch 请求匹配一种游戏方式
func (h *Hall) AskMatch(ac *AccountConn, mt int, data *EchoProto) {
	name := data.Data["account"].(string)
	game := data.Data["game"].(string)

	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if _, ok := h.Accounts[name]; ok {
		go GetMatchSys().AskMatch(name, game, h.Accounts[name].Step, h.Accounts[name].Elo, float64(h.Accounts[name].WinRate), float64(h.Accounts[name].Kda))
	}
}

// 匹配服,成功匹配
func (h *Hall) OnMatchOver(accounts []string) {
	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	for i := 0; i < len(accounts); i++ {
		if _, ok := h.Accounts[accounts[i]]; !ok {
			println("帐号", accounts[i], " 离线")
			break
		}
		println("匹配成功", accounts)
	}

	// 生成房间(Room), 开搞喽
}

// SendMail 发送邮件
func (h *Hall) SendMail(ac *AccountConn, mt int, data *EchoProto) {

}

// 房间战斗结束, 返回战斗结果, 每个玩家的信息分开
func (h *Hall) OnRoomOver(battles BattleInfo) {
	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if _, ok := h.AccountIDs[battles.AccID]; ok {
		println("战报:", h.AccountIDs[battles.AccID].Name)
	} else {
		println("帐号不存在:", battles.AccID)
	}
}

// 处理返回邮件
func (h *Hall) OnRecvMails(accID int64, mails []Mail) {
	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

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
