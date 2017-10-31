package game

// 房间
type Room struct {
	Name          string                  // 房间名称
	MyMap         Map                     // 地图
	Accouts       map[string]*AccountData // 帐号
	LastRoleId    int64                   // 最后一个角色ID
	LastMonsterId int64                   // 最后一个怪物ID
	// 挂靠的玩法
}

// 初始化
func (r *Room) Init(name string) {
	r.Name = name
	r.MyMap.Init(r)
	r.Accouts = make(map[string]*AccountData)
}

// 帐号的登录
func (r *Room) Login(accInfo AccountData) bool {
	if _, ok := r.Accouts[accInfo.Account]; ok {
		return false
	}
	n := new(AccountData)
	*n = accInfo
	n.MyRoom = r
	r.Accouts[accInfo.Account] = n

	r.MyMap.AccountLogin(n, true)
	return true
}

// 帐号离开
func (r *Room) Logout(acc string) bool {
	if v, ok := r.Accouts[acc]; ok {
		r.MyMap.AccountLogin(v, false)
		delete(r.Accouts, acc)
		return true
	}
	return false
}

//
func (r *Room) MakeRoleId() int64 {
	r.LastRoleId++
	return r.LastRoleId
}

func (r *Room) MakeMonsterId() int64 {
	r.LastMonsterId++
	return r.LastMonsterId
}
