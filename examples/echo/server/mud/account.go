package mud

// AccountInfo 帐号
type AccountInfo struct {
	ID         int64   // 帐号ID
	Name       string  // 帐号名
	Nick       string  // 昵称
	Step       int32   // 段位
	WinRate    float32 // 胜率
	Kda        float32 // K/D/A 值 (杀/死/助), (K+A)/D, 如果D为0, 则D为1
	Elo        int32   // Elo积分
	KillCount  int32   // 杀敌数
	DeadCount  int32   // 死亡数
	Assist     int32   // 助攻数
	RegistTime int32   // 注册时间
}

// AccountReal
type AccountReal struct {
	AccountInfo
	Online   bool  // 在线么
	LastTime int32 // 最近登录时间
}

// BattleInfo 战斗记录
type BattleInfo struct {
	AccID    int64 // 帐号ID
	RoomID   int64 // 房间ID
	GameType int32 // 游戏类型
	MapType  int32 // 地图类型
	RoleType int32 // 角色类型
	Kill     int32 // K
	Dead     int32 // D
	Assist   int32 // A
	Win      int32 // 胜利/失败 1/0
	Mvp      bool  // 是否MVP
}
