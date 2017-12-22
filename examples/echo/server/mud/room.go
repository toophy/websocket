package mud

// Room 房间
type Room struct {
	RoomID      int64         // 房间ID
	GameType    int32         // 游戏类型
	MapType     int32         // 地图类型
	Accounts    []AccountReal // 帐号们
	BattleInfos []BattleInfo  // 战场信息们
}

type RoomSys struct {
	Rooms map[int64]*Room
}

var (
	roomSys *RoomSys
)

func init() {
	roomSys = &RoomSys{}
	roomSys.Rooms = make(map[int64]*Room, 0)
}

func GetRoomSys() *RoomSys {
	return roomSys
}

func (r *Room) Run() {
	// 每次运行就是一帧
}

func (r *RoomSys) MakeRoom(gameType, mapType int32, accounts []AccountReal) bool {
}

// 响应帐号投递的消息
func (r *RoomSys) OnAccountPost(account string, op int32, param1 int32, param2 int32) {
}

// 战斗结束, 负责保存结果到 BattleInfos, 退回大厅,
func (r *RoomSys) Update() {
	// 检查每一场战斗?
	// 还是每一场战斗 单独协程
}
