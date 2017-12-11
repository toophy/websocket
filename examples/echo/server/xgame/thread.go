package xgame

type AccountData struct {
	Name  string
	AccID int64
	ObjID int32
}

type MapThread struct {
	FrameSN  int32                   // 帧序号
	Map      *BoxMap                 // 地图
	Accounts map[string]*AccountData // 帐号
	// 网络消息收集
	// 广播帧消息
	// 缓冲帧
}

// Init 初始化
func (m *MapThread) Init() {
	m.FrameSN = 0
}

// 刷新一帧
func (m *MapThread) UpdateFrame() {
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
func (m *MapThread) InsertAccount(acc string, accID int64) {
	if _, ok := m.Accounts[acc]; !ok {
		m.Accounts[acc] = &AccountData{Name: acc, AccID: accID, ObjID: 0}
	}
}

// GetObjIDByAccount 通过玩家账号获取对象ID
func (m *MapThread) GetObjIDByAccount(acc string) int32 {
	if _, ok := m.Accounts[acc]; ok {
		return m.Accounts[acc].ObjID
	}
	return 0
}

// 模拟玩家操作
func (m *MapThread) InputOp(acc string, opID int32, opParam1 int32, opParam2 int32) {
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
