package xgame

import "sync"

// MapConfig 地图配置信息
type MapConfig struct {
	Name    string // 地图模板名称
	ID      int32  // 地图模板ID
	TileW   int32  // 九宫格宽
	TileH   int32  // 九宫格高
	Width   int32  // 实际宽
	Height  int32  // 实际高
	Gravity int32  // 重力
}

// 初始化位置
// 初始化建筑
// 初始化地块
// 初始化怪物
// 初始化NPC
// 初始化陷阱
// 初始化xxx
// 初始化季节
// 初始化河流
// ...

var (
	mapConfigMutex *sync.Mutex           //
	mapConfig      map[string]*MapConfig // 地图配置中心
	mapConfigID    map[int32]*MapConfig  // 地图配置中心
)

// 初始化
func init() {
	mapConfigMutex = new(sync.Mutex)
	mapConfig = make(map[string]*MapConfig, 0)
	mapConfigID = make(map[int32]*MapConfig, 0)
}

func addMapConfig(name string, id int32, tileW, tileH, realW, realH int32) bool {
	mapConfigMutex.Lock()
	defer mapConfigMutex.Unlock()

	if tileW < 1 || tileH < 1 || realW < 0 || realH < 0 || tileW > realW || tileH > realH {
		return false
	}

	n := &MapConfig{Name: name, ID: id, TileW: tileW, TileH: tileH, Width: realW, Height: realH}
	mapConfig[name] = n
	mapConfigID[id] = n

	return true
}

func delMapConfig(name string, id int32) {
	mapConfigMutex.Lock()
	defer mapConfigMutex.Unlock()

	if _, ok := mapConfigID[id]; ok {
		delete(mapConfig, mapConfigID[id].Name)
		delete(mapConfigID, id)
	} else if _, ok := mapConfig[name]; ok {
		delete(mapConfigID, mapConfig[name].ID)
		delete(mapConfig, name)
	}
}

// GetMapConfig 获取地图配置信息
func GetMapConfig(name string, id int32) (mc MapConfig) {
	mapConfigMutex.Lock()
	defer mapConfigMutex.Unlock()

	if _, ok := mapConfigID[id]; ok {
		mc = *mapConfigID[id]
	} else if _, ok := mapConfig[name]; ok {
		mc = *mapConfig[name]
	}
	return
}
