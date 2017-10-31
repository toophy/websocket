package game

const (
	LOT_Land    = 1 // 陆地
	LOT_Brick   = 2 // 陆地对象类型 地块
	LOT_Grow    = 3 // 陆地对象类型 生长点, 生长各种临时道具
	LOT_Player  = 4 // 玩家角色
	LOT_Monster = 5 // 小怪物
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
	Childs    map[int64]*MapObj // 子对象, 决定陆地长度, 暂定为直线段
}

// 初始化
func (m *MapObj) Init() {
	m.Attrs.Init()
	m.Childs = make(map[int64]*MapObj, 0)
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
	SrcX       int32
	SrcY       int32
	DstX       int32
	DstY       int32
	BeginFrame int64
	EndFrame   int64
	Speed      int32
}

// 移动
func (m *MapObj) Move() {
	// 当前Mover走完了么?
	// 制作 Mover信息
	// 超出父对象范围
	// 1. 有目标陆地
	// 2. 没有目标标陆地
}

// 刷新移动
func (m *MapObj) UpdateMove() {

}
