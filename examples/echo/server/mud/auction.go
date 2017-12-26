package mud

// 拍卖

// 从 KuangSys 找到一个, 进行拍卖
// 拍卖的货币使用所有的兼容货币, 使用当时的汇率进行价格比对, 然后竞价
// 每次至少加价 1%
//
// 因为货币是整数单位, 所以加价也是整数增加, 至少增加一个货币
//
// 初定货币
//
// 以下是货币存量, 存量越少, 价值越高
//
// 吃货	一场战斗中吃的最多  ....   10000
// 好基友	一场战斗中助攻最多 ....   5000
// 技术达人	一场战斗中杀人最多 ....  3000
//
// 必须知道货币存量
//
// 生成汇率表
// 汇率每5分钟一次变更?
//

// 当前货币汇率表
// 一次只能进行一场竞拍
// 起拍价 15%可以结束拍卖(封顶), 否则5分钟等待...

type AuctionLog struct {
	AccName string // 竞拍人帐号
	Count   int64  // 货币数量
	Time    int64  // 竞拍时间
}

// AuctionSession 拍卖, 拍卖时冻结指定的 Jewel 数量
type AuctionSession struct {
	ID          int64         // 拍卖序号, 第xx场拍卖
	Jewel       int64         // 宝石
	LastCounts  int64         // 货币数量
	StartTime   int64         // 起拍时间
	EndTime     int64         // 预计结拍时间
	AuctionOver bool          // 竞拍结束
	Logs        []*AuctionLog // 竞拍日志
}

// Auction 宝石拍卖系统
type JewelAuction struct {
	KuangName     string                    // 矿名
	LastItemID    int64                     // 货币类型(道具ID)
	LastSessionID int64                     // 最后场次
	Sessions      map[int64]*AuctionSession // 拍卖场次
}

// Auction 宝石拍卖系统
type AuctionSys struct {
	JewelAuctions map[int64]*JewelAuction // 拍卖场次
}

var (
	auctionSys *AuctionSys
)

func init() {
	auctionSys = &AuctionSys{}
	auctionSys.JewelAuctions = make(map[int64]*JewelAuction, 0)
}

func GetAuctionSys() *AuctionSys {
	return auctionSys
}

func (a *AuctionSys) AddAuction(name string, itemID int64) {
	if _, ok := a.JewelAuctions[name]; !ok {
		j := &JewelAuction{}

		j.KuangName = name
		j.LastItemID = itemID
		j.LastSessionID = 1
		j.Sessions = make(map[int64]*JewelAuction, 0)

		a.JewelAuctions[name] = j
	}
}

func (a *AuctionSys) CreateAuctionSession(name string, jewelCount int64, accName string, itemCount int64) bool {

	if j, ok := a.JewelAuctions[name]; ok {

		lastCounts := 0 // 货币数量

		if len(accName) < 1 {
			return false
		}
		// 对accName检查

		if itemCount < 1 {
			return false
		}

		// 1. 当前拍卖是否结束
		if len(j.Sessions) > 0 {
			if _, ok := j.Sessions[j.LastID]; ok {
				if !j.Sessions[j.LastID].AuctionOver {
					return false
				}
				lastCounts = j.Sessions[j.LastID].LastCounts
			}
		}

		// 2. jewelCount 是否正确, 到矿石系统检查
		if !GetKuangSys().AskFreezeJewel(j.KuangName, jewelCount) {
			return false
		}

		// 3. itemCount 是否合法, 基础竞价
		if itemCount < int64(lastCounts*0.9) {
			return false
		}

		//
		as := &AuctionSession{}

		as.ID = j.LastSessionID
		as.Jewel = jewelCount
		as.LastCounts = 0
		as.StartTime = time.Now().Unix()
		as.EndTime = time.Now().Unix() + 3*60
		as.AuctionOver = false
		as.Logs = make([]*AuctionLog, 0)

		al := AuctionLog{AccName: accName, Count: itemCount, Time: time.Now().Unix()}

		j.Sessions[as.ID] = as

		j.LastSessionID++
	}
}

func (a *AuctionSys) AskAuction(name string, accName string, itemCount int64) bool {

	if j, ok := a.JewelAuctions[name]; ok {

		lastCounts := 0 // 货币数量

		if len(accName) < 1 {
			return false
		}
		// 对accName检查

		if itemCount < 1 {
			return false
		}

		// 1. 当前拍卖是否结束
		if len(j.Sessions) > 0 {
			if _, ok := j.Sessions[j.LastID]; ok {
				if !j.Sessions[j.LastID].AuctionOver {
					return false
				}
				lastCounts = j.Sessions[j.LastID].LastCounts
			}
		}

		// 2. jewelCount 是否正确, 到矿石系统检查
		if !GetKuangSys().AskFreezeJewel(j.KuangName, jewelCount) {
			return false
		}

		// 3. itemCount 是否合法, 基础竞价
		if itemCount < int64(lastCounts*0.9) {
			return false
		}
	}
}

// 玩家拍卖时冻结指定的货币, ItemID, ItemCount, 拍卖结束解冻, 可能使用, 也可能兑换了Jewel
// 玩家宝石也有冻结, 活动两种
// 玩家宝石有利息

// 全服只有一个竞拍, 全服共享竞拍
//
