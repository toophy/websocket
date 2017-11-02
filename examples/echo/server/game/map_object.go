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
	Default_Turn_Time = 3
	Default_Save_Pos  = 10
)

// 地图对象
type MapObj struct {
	ID        int64             // ID
	CId       int16             // child ID
	Pos       Vec4              // 位置
	Type      int32             // 类型
	Attrs     RoleAttr          // 属性
	Myroom    *Room             // 房间
	MyAccount *AccountData      // 帐号
	MyMover   Mover             // 移动子
	OldPoss   map[int64]Vec2    // 旧位置(每帧), 可以保留之前的 Max_OldPos 帧位置
	Childs    map[int64]*MapObj // 子对象
	Orders    []*MapObj         // 直线段队形
}

// 初始化
func (m *MapObj) Init() {
	m.Attrs.Init()
	m.OldPoss = make(map[int64]Vec2, Max_OldPos)
	m.Childs = make(map[int64]*MapObj, 0)
	m.Orders = make([]*MapObj, 0)
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

// 初始化
func (m *Mover) Init() {
	m.SrcX = 0
	m.SrcY = 0
	m.Left = false
	m.BeginFrame = 0
	m.EndFrame = 0
	m.Speed = 0
	m.Moving = false
	m.TurnState = 0
	m.TurnFrame = 0
	m.TurnSpeed = 0
}

// 移动
func (m *MapObj) Move(srcX, srcY int32, left bool, beginFrame int64) bool {
	if !m.MyMover.Moving {
		// [广播消息]
		m.MyMover.SrcX = srcX
		m.MyMover.SrcY = srcY
		m.MyMover.Left = left
		m.MyMover.BeginFrame = beginFrame
		m.MyMover.EndFrame = m.MyMover.EndFrame + Default_Move_Time
		m.MyMover.Speed = m.Attrs.GetAttrVal32(Attr_speed)
		m.MyMover.Moving = true
		m.MyMover.TurnState = Move_None
		m.MyMover.TurnFrame = 0
		m.MyMover.TurnSpeed = 0
	} else if m.MyMover.TurnState == Move_None {
		if m.MyMover.Left != left {
			// 方向相反 [广播消息]
			m.MyMover.TurnState = Move_Turn
			m.MyMover.TurnFrame = beginFrame
			m.MyMover.TurnSpeed = m.Attrs.GetAttrVal32(Attr_speed)
		} else {
			// 方向一致, 可能速度有变动
			if m.MyMover.Speed != m.Attrs.GetAttrVal32(Attr_speed) {
				// 仅仅是速度变化, 产生新的移动, 当前移动要保持惯性 [广播消息]
				m.MyMover.TurnState = Move_Speed
				m.MyMover.TurnFrame = beginFrame
				m.MyMover.TurnSpeed = m.Attrs.GetAttrVal32(Attr_speed)
			} else {
				// 拉长移动距离, [广播消息]
				m.MyMover.EndFrame = m.MyMover.EndFrame + m.Myroom.FrameSn - m.MyMover.BeginFrame
			}
		}
	}

	return true
}

// 移动
func (m *MapObj) changeMove(left bool) {
	m.MyMover.SrcX = m.Pos.X
	m.MyMover.SrcY = m.Pos.Y
	m.MyMover.Left = left
	m.MyMover.BeginFrame = m.Myroom.FrameSn
	m.MyMover.EndFrame = m.MyMover.BeginFrame + Default_Move_Time
	m.MyMover.Speed = m.MyMover.TurnSpeed
	m.MyMover.Moving = true
	m.MyMover.TurnState = Move_None
	m.MyMover.TurnFrame = 0
	m.MyMover.TurnSpeed = 0
}

// 刷新移动
func (m *MapObj) UpdateMove() {

	if len(m.OldPoss) > 0 {
		delete(m.OldPoss, m.Myroom.FrameSn-Default_Save_Pos)
	}

	if !m.MyMover.Moving {
		return
	}

	m.OldPoss[m.Myroom.FrameSn] = Vec2{m.Pos.X, m.Pos.Y}

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
	// 位置 --> 查找陆地, 地图对象, 碰撞?
	// 所有陆地查询一次
	// 陆地有长度
	//

	switch m.MyMover.TurnState {
	case Move_None:
		if m.Myroom.FrameSn >= m.MyMover.EndFrame {
			m.MyMover.Init()
			return
		}
		break
	case Move_Turn:
		// 要求转向
		if m.Myroom.FrameSn >= m.MyMover.TurnFrame+Default_Turn_Time {
			m.changeMove(!m.MyMover.Left)
		}
		break
	case Move_Speed:
		// 要求变速
		if m.Myroom.FrameSn >= m.MyMover.TurnFrame+Default_Turn_Time {
			m.changeMove(m.MyMover.Left)
		}
		break
	}
}
