//
package game

// 地图
type Map struct {
	ConfigFile  string            // 地图对应的文件(json格式,有版本号)
	Version     int64             // 版本号
	Widht       int64             // 宽
	Height      int64             // 高
	BrickWidth  int16             // 砖块的tile宽
	BrickHeight int16             // 砖块的tile高
	RoleWidth   int16             // role的tile宽
	RoleHeight  int16             // role的tile高
	StartTime   int64             // 开始时间
	EndTime     int64             // 存活时间
	FrameTime   int64             // 一帧的时间
	FrameID     int64             // 当前是第几帧, 一般是1秒5帧
	Lands       map[int64]*MapObj // 陆地列表
	Roles       map[int64]*MapObj // 角色列表
	Monsters    map[int64]*MapObj // 怪物列表
	MyRoom      *Room             // 房间
}

// 初始化
func (m *Map) Init(r *Room) {
	m.Lands = make(map[int64]*MapObj, 0)
	m.Roles = make(map[int64]*MapObj, 0)
	m.Monsters = make(map[int64]*MapObj, 0)
	m.MyRoom = r
}

// 帐号登录/登出
func (m *Map) AccountLogin(accInfo *AccountData, login bool) bool {
	return true
}

// 角色的创建
func (m *Map) CreateRole(accInfo *AccountData, roleName string) *MapObj {
	if accInfo.RoleId != 0 {
		return nil
	}
	n := new(MapObj)
	n.Init()
	n.ID = m.MyRoom.MakeRoleId()

	n.Pos.X = 0
	n.Pos.Y = 0
	n.Pos.Parent = nil
	n.Pos.Height = 0

	n.Type = LOT_Player
	n.Myroom = m.MyRoom
	n.MyAccount = accInfo
	return n
}

// 怪物的创建
func (m *Map) CreateMonster(monsterName string, monsterType int32) *MapObj {
	if len(monsterName) < 1 || monsterType < 1 || monsterType == LOT_Player {
		return nil
	}
	n := new(MapObj)
	n.Init()
	n.ID = m.MyRoom.MakeMonsterId()

	n.Pos.X = 0
	n.Pos.Y = 0
	n.Pos.Parent = nil
	n.Pos.Height = 0

	n.Type = monsterType
	n.Myroom = m.MyRoom
	n.MyAccount = nil
	return n
}

// 陆地的创建
func (m *Map) CreateLand(landName string) *MapObj {
	if len(landName) < 1 {
		return nil
	}
	n := new(MapObj)
	n.Init()
	n.ID = m.MyRoom.MakeMonsterId()

	n.Pos.X = 0
	n.Pos.Y = 0
	n.Pos.Parent = nil
	n.Pos.Height = 0

	n.Type = LOT_Land
	n.Myroom = m.MyRoom
	n.MyAccount = nil
	return n
}

// 没有创建角色的帐号, 有浏览权限么?
// 权限的分配
// 权限在Room的帐号里面
