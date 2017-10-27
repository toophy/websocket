//
package game

// 房间
type Room struct {
	MyMap   Map                     // 地图
	Accouts map[string]*AccountData // 帐号
	// 挂靠的玩法
}

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
}

// 位置
type Vec4 struct {
	X      int32   // x
	Y      int32   // y, 角色在陆地上, 没有这个属性变化, 飘在场景中有这个属性
	Height int32   // 高度, 在陆地上的高度
	Parent *MapObj // 父对象
}

const (
	LOT_Land  = 1 // 陆地
	LOT_Brick = 2 // 陆地对象类型 地块
	LOT_Grow  = 3 // 陆地对象类型 生长点, 生长各种临时道具
)

// 地图对象
type MapObj struct {
	ID     int64
	Pos    Vec4              // 位置
	Type   int32             // 类型
	Param  int64             // 参数
	Childs map[int64]*MapObj // 子对象, 决定陆地长度, 暂定为直线段
}
