package xgame

import "math/rand"

// AccountData 帐号信息
type AccountData struct {
	Name  string // 帐号
	AccID int64  // 帐号ID
	Nick  string // 昵称
	ObjID int32  // 对象ID
}

// Room 房间
type Room struct {
	FrameSN  int32                   // 帧序号
	Map      *BoxMap                 // 地图
	Accounts map[string]*AccountData // 帐号
	// 网络消息收集
	// 广播帧消息
	// 缓冲帧
}

// Init 初始化
func (m *Room) Init() {
	m.FrameSN = 0
}

// createMap 创建地图
func (m *Room) createMap(name string, id int32) bool {
	mc := GetMapConfig(name, id)
	if mc.ID > 0 {
		m.Map = new(BoxMap)
		m.Map.InitMap(mc, rand.Int63())
	}
	return false
}

// UpdateFrame 刷新一帧
func (m *Room) UpdateFrame() {
	m.FrameSN++

	// 广播玩家帧消息 --> 所有帐号

	// 处理玩家地图
	if m.Map != nil {
		m.Map.Update(m.FrameSN)
	}

	// ?
}

// 接收网络消息
//     解释,保存网络消息, 用帧保存, 刷新帧的时候, 广播帧消息

// InsertAccount 插入一个帐号
func (m *Room) InsertAccount(acc string, accID int64, nick string) {
	if _, ok := m.Accounts[acc]; !ok {
		m.Accounts[acc] = &AccountData{Name: acc, AccID: accID, Nick: nick, ObjID: 0}
	}
}

// GetObjIDByAccount 通过玩家账号获取对象ID
func (m *Room) GetObjIDByAccount(acc string) int32 {
	if _, ok := m.Accounts[acc]; ok {
		return m.Accounts[acc].ObjID
	}
	return 0
}

// NewRole 新建帐号的角色
func (m *Room) NewRole(acc string, role_type string) *Obj {
	if _, ok := m.Accounts[acc]; ok {
		//随机位置
		birthRect := m.Map.RandBirthPos(100, 100)
		return m.Map.NewObj(birthRect.X, birthRect.Y, birthRect.W, birthRect.H)
	}
	return nil
}

// 模拟玩家操作
func (m *Room) InputOp(acc string, opID int32, opParam1 int32, opParam2 int32) {
	objID := m.GetObjIDByAccount(acc)
	if objID == 0 {
		return
	}

	// 下一个帧
	nextFrameSN := m.FrameSN + 1
	data := &FrameInput{ID: objID, OpID: opID, OpParam1: opParam1, OpParam2: opParam2}
	if _, ok := m.Map.FrameDatas[nextFrameSN]; ok {
		if m.Map.FrameDatas[nextFrameSN].Ops == nil {
			m.Map.FrameDatas[nextFrameSN].Ops = make([]*FrameInput, 0)
		}
		m.Map.FrameDatas[nextFrameSN].Ops = append(m.Map.FrameDatas[nextFrameSN].Ops, data)
	} else {
		frameData := &FrameData{Sn: nextFrameSN}
		frameData.Ops = make([]*FrameInput, 0)
		frameData.Ops = append(frameData.Ops, data)
		m.Map.FrameDatas[nextFrameSN] = frameData
	}
}
