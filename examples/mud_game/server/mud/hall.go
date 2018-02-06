package mud

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
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
					go GetMailSys().GetNextMails(v.AccountInfo.ID, v.LastMailID, h.OnGetMails)
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
		if a.C != nil {
			retData := *data
			retData.Data["ret"] = "ok"
			retData.Data["msg"] = "正常离线"
			ret, _ := json.Marshal(retData)

			err := a.C.WriteMessage(mt, ret)
			if err != nil {
				fmt.Printf("[W] 向网络写消息失败:%s\n", err)
			}

			go func() {
				select {
				case <-time.After(1 * time.Second):
					a.C.Close()
					break
				}
			}()
		}

		delete(hall.AccountIDs, h.Accounts[name].ID)
		delete(hall.Accounts, name)
	} else {
		if a.C != nil {
			retData := *data
			retData.Data["ret"] = "failed"
			retData.Data["msg"] = "重复离线"
			ret, _ := json.Marshal(retData)

			err := a.C.WriteMessage(mt, ret)
			if err != nil {
				fmt.Printf("[W] 向网络写消息失败:%s\n", err)
			}
			go func() {
				select {
				case <-time.After(1 * time.Second):
					a.C.Close()
					break
				}
			}()
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
		retData.Data["msg"] = "重复登录"
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
			LastGetMailTime: time.Now().UnixNano(),
			Conn:            ac}

		h.AccountIDs[a.ID] = h.Accounts[name]

		fmt.Printf("[I] 帐号成功登录(%s : %s)\n", name, pwd)

		retData := *data
		retData.Data["ret"] = "ok"
		retData.Data["msg"] = ""

		ret, _ := json.Marshal(retData)
		err := ac.C.WriteMessage(mt, ret)
		if err != nil {
			fmt.Printf("[W] 向网络写消息失败:%s\n", err)
		}
	} else {
		fmt.Printf("[I] 帐号未注册:%s\n", name)

		retData := *data
		retData.Data["ret"] = "fail"
		retData.Data["msg"] = "帐号未注册"
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

// GetMails 主动收取邮件
func (h *Hall) GetMails(ac *AccountConn, mt int, data *EchoProto) {

	last_id := data.Data["last_id"].(int32)

	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if a, ok := h.Accounts[ac.Account]; ok {
		go GetMailSys().GetNextMails(a.ID, last_id, h.OnGetMails)
	}
}

// OnGetMails 主动收取邮件的反馈
func (h *Hall) OnGetMails(accID int64, ret string, msg string) {
	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if ac, ok := h.AccountIDs[accID]; ok {
		if ac.Conn == nil || ac.Conn.C == nil {
			return
		}

		retJson, _ := json.Marshal(&EchoProto{
			C: "Index",
			M: "GetMails",
			Data: map[string]interface{}{
				"ret": ret,
				"msg": msg}})

		if err := ac.Conn.C.WriteMessage(websocket.TextMessage, retJson); err != nil {
			fmt.Printf("[W] 向网络写消息失败:%s\n", err)
		}
	}
}

// SendMail 发送邮件
func (h *Hall) SendMail(ac *AccountConn, mt int, data *EchoProto) {
	tempID := data.Data["tempid"].(string)
	recer := data.Data["recer"].(string)
	title := data.Data["title"].(string)
	content := data.Data["content"].(string)

	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if a, ok := h.Accounts[ac.Account]; ok {
		go GetMailSys().SendByName(a.ID, recer, title, content, "", h.OnSendMail, tempID)
	}
}

// OnSendMail 主动发送邮件的反馈
func (h *Hall) OnSendMail(accID int64, tempID string, ret string) {
	h.AccountLocker.Lock()
	defer h.AccountLocker.Unlock()

	if ac, ok := h.AccountIDs[accID]; ok {
		if ac.Conn == nil || ac.Conn.C == nil {
			return
		}

		retJson, _ := json.Marshal(&EchoProto{
			C: "Index",
			M: "SendMail",
			Data: map[string]interface{}{
				"ret": ret,
				"msg": ""}})

		if err := ac.Conn.C.WriteMessage(websocket.TextMessage, retJson); err != nil {
			fmt.Printf("[W] 向网络写消息失败:%s\n", err)
		}
	}
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
