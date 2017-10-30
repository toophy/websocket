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
	Childs    map[int64]*MapObj // 子对象, 决定陆地长度, 暂定为直线段
}

// 初始化
func (m *MapObj) Init() {
	m.Attrs.Init()
	m.Childs = make(map[int64]*MapObj, 0)
}
