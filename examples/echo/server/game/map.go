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
	MyRoom      *Room             // 房间
}

// 初始化
func (m *Map) Init(r *Room) {
	m.Lands = make(map[int64]*MapObj, 0)
	m.Roles = make(map[int64]*MapObj, 0)
	m.MyRoom = r
}
