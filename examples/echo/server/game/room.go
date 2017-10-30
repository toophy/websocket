package game

// 房间
type Room struct {
	Name    string                  // 房间名称
	MyMap   Map                     // 地图
	Accouts map[string]*AccountData // 帐号
	// 挂靠的玩法
}

// 初始化
func (r *Room) Init(name string) {
	r.Name = name
	r.MyMap.Init(r)
	r.Accouts = make(map[string]*AccountData)
}
