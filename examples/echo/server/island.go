//
package main

import "github.com/gorilla/websocket"

type AccountData struct {
	ID      int64
	Account string
	Pwd     string
	Mt      int
	C       *websocket.Conn
	PosX    float64
	PosY    float64
}

// 房间
type Room struct {
	TheScene Scene                   // 场景
	Accouts  map[string]*AccountData // 帐号
	// 排行榜
	// 挂靠的玩法
}

// 场景
type Scene struct {
	ConfigFile            string            // 地图对应的文件(json格式,有版本号)
	Version               int64             // 版本号
	IsLandBlockTileWidth  int16             // 浮岛块的tile宽
	IsLandBlockTileHeight int16             // 浮岛块的tile高
	RoleTileWidth         int16             // role的tile宽
	RoleTileHeight        int16             // role的tile高
	StartTime             int64             // 开始时间
	EndTime               int64             // 存活时间
	FrameTime             int64             // 一帧的时间
	FrameID               int64             // 当前是第几帧, 一般是1秒5帧
	IsLands               map[int64]*IsLand // 浮岛列表
	Roles                 map[int64]*Role   // 角色列表
}

// 位置
type Vec4 struct {
	X        int32 // x
	Y        int32 // y, 角色在浮岛上, 没有这个属性变化, 飘在场景中有这个属性
	Height   int32 // 高度, 在浮岛上的高度
	IsLandID int64 // 浮岛, 浮岛为空, 角色/浮岛就漂浮在场景中, 如果一直没有落在浮岛上, 角色落下场景死亡, 浮岛也有可能也会落下场景
}

// 浮岛地块, 也是一种角色, 具备 血量等属性
type IsLandBlock struct {
	Pos   Vec4  // 位置
	Type  int32 // 类型
	Param int64 // 参数
}

// 浮岛
type IsLand struct {
	Pos    Vec4                   // 位置
	Blocks map[int16]*IsLandBlock // 每个地块属性
	// 浮岛长度
	// 浮岛重量
	// 浮岛移动相关(速度, 加速度, 惯性)
	// 浮岛上的对象层(草丛,刷新点,角色)
}

// 角色
type Role struct {
	Account string // 帐号
	Name    string // 角色名
	ID      int64  // 角色ID
	Pos     Vec4   // 位置
	SceneID int64  // 角色所在场景ID
	// 角色重量
	// 角色移动相关(速度, 加速度, 惯性)
	// 角色血量, 体力, 攻击力...
}

// 浮岛对象
type IsLandObj struct {
	// 浮岛
	// 序号(x坐标)
	// 类型
	// 参数
	//
}

// 移动具有格子属性, 不按照像素点移动, 更好把握同步运算
// 浮岛上可以再加上浮岛, 比如玩家的交通工具(暂时不支持)
// 角色的位置, 从浮岛到场景, 从场景到浮岛

// 对象基本属性
// ID, 名称, 位置, 血量, 体力, 速度, 惯性(固定规则), 攻击力, 防御力,

// 速度
// 帧/每块, 精确每帧的位置
// 除了浮岛块的坐标, 还有一层角色坐标么?
// 角色比浮岛块要小么? 小多少? 1:3?
//
