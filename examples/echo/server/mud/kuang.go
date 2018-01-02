package mud

import (
	"errors"
	"math"
	"sync"
	"time"
)

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
	ItemID        int64                     // 货币类型(道具ID)
	LastSessionID int64                     // 最后场次
	Sessions      map[int64]*AuctionSession // 拍卖场次
}

// 很多种矿
type KuangSys struct {
	Name          string                  // 矿名
	Ore           float64                 // 原矿石
	Jewel         int64                   // 宝石, 从矿石中提取整数部分
	OreSpeed      float64                 // 原矿石生产速度 Ore/秒
	Lucky         float64                 // 幸运值
	LastJewelTime int64                   // 最后一次宝石计算时间
	JewelAuctions map[int64]*JewelAuction // 拍卖场
}

var (
	gItems   map[int64]string
	kuangSys *KuangSys
)

func init() {
	gItems = make(map[int64]string, 0)
	gItems[1] = "第一名"
	gItems[2] = "打酱油"

	kuangSys = &KuangSys{}
	// 要用从序列化文件导入-->数据库
	kuangSys.Name = "钻石"
	kuangSys.Ore = 0.0
	kuangSys.Jewel = 0
	kuangSys.OreSpeed = 1.0 / 60
	kuangSys.Lucky = 0.0
	kuangSys.LastJewelTime = time.Now().Unix()

	kuangSys.JewelAuctions = make(map[string]*JewelAuction, 0)

}

func GetKuangSys() *KuangSys {
	return kuangSys
}

// 每秒一次计算, 每分钟一次计算
func (k *KuangSys) Update() {
	t := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-t.C:
			k.Ore += k.OreSpeed
			if time.Now().Unix() > k.LastJewelTime+60 {
				println(k.Name, k.Ore, k.Jewel, k.OreSpeed)
				k.LastJewelTime = time.Now().Unix()

				if k.Ore > 1 {
					newJewel := math.Floor(k.Ore)
					if newJewel > 0 {
						k.Jewel += int64(newJewel)
						k.Ore -= newJewel
					}
				}
			}

			t.Reset(1 * time.Second)
			break
		}
	}
}

// AccAuction 增加一个拍卖场
func (a *KuangSys) AddAuction(itemID int64) error {
	if _, ok := gItems[itemID]; ok {
		if _, ok := a.JewelAuctions[itemID]; !ok {
			j := &JewelAuction{}

			j.ItemID = itemID
			j.LastSessionID = 1
			j.Sessions = make(map[int64]*JewelAuction, 0)

			a.JewelAuctions[itemID] = j
			return nil
		} else {
			return errors.New("已经有这个拍卖场")
		}
	}

	return errors.New("没有这个通货")
}

// CreateAuctionSession 新建一场拍卖会
func (a *KuangSys) CreateAuctionSession(itemID int64, jewelCount int64, accName string, itemCount int64) error {

	if j, ok := a.JewelAuctions[itemID]; ok {

		lastItemCounts := 0

		if len(accName) < 1 {
			return errors.New("帐号不存在")
		}

		if itemCount < 1 {
			return errors.New("通货数量要>0")
		}

		if len(j.Sessions) > 0 {
			if _, ok := j.Sessions[j.LastID]; ok {
				if !j.Sessions[j.LastID].AuctionOver {
					return errors.New("拍卖会还在进行中")
				}
				lastItemCounts = j.Sessions[j.LastID].LastCounts
			}
		}

		if a.Jewel < jewelCount {
			return errors.New("钻石不够哦")
		}

		if itemCount < int64(lastItemCounts*0.9) {
			return errors.New("货币不足最后一次拍卖90%")
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
		return nil
	}

	return errors.New("没有这种通货拍卖场")
}

func (a *KuangSys) AskAuction(name string, accName string, itemCount int64) bool {

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
