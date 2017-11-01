package game

const (
	LOT_Land    = 1 // 陆地
	LOT_Brick   = 2 // 陆地对象类型 地块
	LOT_Grow    = 3 // 陆地对象类型 生长点, 生长各种临时道具
	LOT_Player  = 4 // 玩家角色
	LOT_Monster = 5 // 小怪物
)

const (
	Max_OldPos = 10
)

const (
	Move_None  = 0 // 没有移动需求
	Move_Turn  = 1 // 转向
	Move_Speed = 2 // 变速
)

const (
	Default_Move_Time = 10
)

// 地图对象
type MapObj struct {
	ID        int64             // ID
	Pos       Vec4              // 位置
	Type      int32             // 类型
	Attrs     RoleAttr          // 属性
	Myroom    *Room             // 房间
	MyAccount *AccountData      // 帐号
	MyMover   Mover             // 移动子
	OldPoss   map[int64]Vec2    // 旧位置(每帧), 可以保留之前的 Max_OldPos 帧位置
	Childs    map[int64]*MapObj // 子对象, 决定陆地长度, 暂定为直线段
}

// 初始化
func (m *MapObj) Init() {
	m.Attrs.Init()
	m.Childs = make(map[int64]*MapObj, 0)
	m.OldPoss = make(map[int64]Vec2, Max_OldPos)
}

// 坐标转换
func (m *MapObj) MyLocalToGlobal() (x2, y2 int32) {
	x2 = m.Pos.X
	y2 = m.Pos.Y
	n := m
	for {
		if n.Pos.Parent == nil {
			break
		}
		n = n.Pos.Parent
		x2 += n.Pos.X
		y2 += n.Pos.Y
	}
	return
}

// 坐标转换
func (m *MapObj) LocalToGlobal(x, y int32) (x2, y2 int32) {
	x2 = x
	y2 = y
	n := m
	for {
		x2 += n.Pos.X
		y2 += n.Pos.Y
		if n.Pos.Parent == nil {
			break
		}
		n = n.Pos.Parent
	}
	return
}

func (m *MapObj) GlobalToLocal(x, y int32) (x2, y2 int32) {
	x2 = x
	y2 = y
	n := m
	for {
		x -= n.Pos.X
		y -= n.Pos.Y
		if n.Pos.Parent == nil {
			break
		}
		n = n.Pos.Parent
	}
	return
}

// 在任意位置增加/删除对象
func (m *MapObj) InsertObj(n *MapObj) bool {
	if n == nil || n.ID < 1 {
		return false
	}
	if _, ok := m.Childs[n.ID]; ok {
		return false
	}

	x2, y2 := n.MyLocalToGlobal()
	n.Pos.X, n.Pos.Y = m.GlobalToLocal(x2, y2)
	n.Pos.Parent = m
	m.Childs[n.ID] = n
	return true
}

// 清除一个对象
func (m *MapObj) EraseObj(n *MapObj) bool {
	if n == nil || n.ID < 1 {
		return false
	}
	if _, ok := m.Childs[n.ID]; !ok {
		return false
	}
	n.Pos.X, n.Pos.Y = n.MyLocalToGlobal()
	n.Pos.Parent = nil
	delete(m.Childs, n.ID)
	return true
}

// 移动子
type Mover struct {
	SrcX       int32 // 开始移动的位置x
	SrcY       int32 // 开始移动的位置y
	Left       bool  // 朝向
	Moving     bool  // 移动中
	BeginFrame int64 // 开始行动的帧序号
	EndFrame   int64 // 结束帧
	Speed      int32 // 速度(恒定)
	TurnState  int   // 预备转向
	TurnFrame  int64 // 开始帧
	TurnSpeed  int32 // 新速度
}

// 移动
func (m *MapObj) Move(srcX, srcY int32, left bool, beginFrame int64) bool {
	if !m.MyMover.Moving {
		// [广播消息]
		m.MyMover.SrcX = srcX
		m.MyMover.SrcY = srcY
		m.MyMover.Left = left
		m.MyMover.BeginFrame = beginFrame
		m.MyMover.EndFrame = beginFrame + Default_Move_Time
		m.MyMover.Speed = int32(m.Attrs.GetAttrVal(Attr_speed))
		m.MyMover.Moving = true
		m.MyMover.TurnState = Move_None
		m.MyMover.TurnFrame = 0
		m.MyMover.TurnSpeed = 0
	} else {
		if m.MyMover.Left != left {
			// 方向相反 [广播消息]
			m.MyMover.TurnState = Move_Turn
			m.MyMover.TurnFrame = beginFrame
			m.MyMover.TurnSpeed = int32(m.Attrs.GetAttrVal(Attr_speed))
		} else {
			// 方向一致, 可能速度有变动
			if m.MyMover.Speed != int32(m.Attrs.GetAttrVal(Attr_speed)) {
				// 仅仅是速度变化, 产生新的移动, 当前移动要保持惯性 [广播消息]
				m.MyMover.TurnState = Move_Speed
				m.MyMover.TurnFrame = beginFrame
				m.MyMover.TurnSpeed = int32(m.Attrs.GetAttrVal(Attr_speed))
			} else {
				// 拉长移动距离, [广播消息]
				m.MyMover.EndFrame = m.MyMover.EndFrame + m.Myroom.FrameSn - m.MyMover.BeginFrame
			}
		}
	}

	return true
}

// 刷新移动
func (m *MapObj) UpdateMove() {
	if !m.MyMover.Moving {
		return
	}

	// 终止时间到了
	// 要求转向
	// 要求变速

	if m.MyMover.Left {
		m.Pos.X -= m.MyMover.Speed
	} else {
		m.Pos.X += m.MyMover.Speed
	}

	// 离开当前parent么?
	if m.Pos.Parent == nil {
		// 在map空域中
		// 落到地图下, 结束生命
		// 左右上都没有关系,
		//
		// 到达其他陆地
		// 检查所有陆地
	} else {
		// 在一个陆地上
		// if m.Pos.X<0 || m.Pos.X> m.Pos.Parent {
		// }
		// 可能跳到另外一个陆地上
	}
}
