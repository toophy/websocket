package mud

import (
	"errors"
	"math"
	"sync"
	"time"
)

type AuctionLog struct {
	AccName   string // 竞拍人帐号
	ItemCount int64  // 货币数量
	Time      int64  // 竞拍时间
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
	t2 := time.NewTimer(2 * time.Second)

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

		case <-t2.C:
			for i, v := range k.JewelAuctions {
				if len(v.Sessions) > 0 {
					if _, ok := v.Sessions[v.LastSessionID]; ok {
						if !v.Sessions[v.LastSessionID].AuctionOver {
							if time.Now().Unix() > v.Sessions[v.LastSessionID].EndTime {
								v.Sessions[v.LastSessionID].AuctionOver = true
								// 拍卖场次结束, 钻石归属, 玩家通用货币归属
								// 直接把钻石邮件给玩家, 同时扣除玩家的通货(ItemID,ItemCount)
								// 系统级别的脚本邮件, 也就是唯一的脚本邮件形式
								// 邮件系统接收到该邮件, 如果玩家在线, 自动执行, 不在线, 等上线时执行
								// 此类邮件永久保存
								//
								// 邮件类型 : 脚本邮件
								// 邮件标题 : 拍卖场竞价成功
								// 文本内容 : 恭喜你在xxx拍卖场xxx次拍卖会,使用xx通货成功竞价获取xx钻石
								//            收到邮件后, 扣除xx通货xx个, 获取钻石xx个
								// 脚本内容 : AutionSuccess(xx钻石数量,xx通货ID,xx通货数量)

								// 	type Mail struct {
								// 	ID            int64  // ID
								// 	Title         string // 标题
								// 	SenderType    string // 发送人类型: 玩家, 系统, 战队
								// 	SenderName    string // 发送人名称: SYSTEM,TEAM,Account
								// 	RecverType    string // 接收人类型: 玩家, 系统, 战队
								// 	RecverName    string // 接收人名称: SYSTEM,TEAM,Account
								// 	Content       string // 文本内容(还有客户端脚本)
								// 	Script        string // 脚本内容: 供服务端解析执行
								// 	Read          bool   // 是否已读
								// 	RecvAccessory bool   // 是否接收了附件
								// 	SendTime      int64  // 发送时间
								// 	DestoryTime   int64  // 保存截止时间
								// }
								newMail := &Mail{
									ID : 0,
									Title:"拍卖场竞价成功",
									SenderType:""
								}
							}
						}
					}
				}
			}

			t2.Reset(2 * time.Second)
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
			if _, ok := j.Sessions[j.LastSessionID]; ok {
				if !j.Sessions[j.LastSessionID].AuctionOver {
					return errors.New("拍卖会还在进行中")
				}
				lastItemCounts = j.Sessions[j.LastSessionID].LastCounts
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

		al := AuctionLog{AccName: accName, ItemCount: itemCount, Time: time.Now().Unix()}

		j.Sessions[as.ID] = as

		j.LastSessionID++
		return nil
	}

	return errors.New("没有这种通货拍卖场")
}

// AskAuction 请求拍卖会竞价
func (a *KuangSys) AskAuction(itemID int64, accName string, itemCount int64) error {

	if j, ok := a.JewelAuctions[itemID]; ok {

		lastItemCounts := 0 // 货币数量

		if len(accName) < 1 {
			return errors.New("帐号不存在")
		}

		if itemCount < 1 {
			return errors.New("通货数量要>0")
		}

		if len(j.Sessions) > 0 {
			if _, ok := j.Sessions[j.LastSessionID]; ok {
				if j.Sessions[j.LastSessionID].AuctionOver {
					return errors.New("拍卖会已经结束")
				}
				lastItemCounts = j.Sessions[j.LastSessionID].LastCounts

				if len(j.Sessions[j.LastSessionID].Logs) > 0 {
					if j.Sessions[j.LastSessionID].Logs[len(j.Sessions[j.LastSessionID].Logs)-1].AccName == accName {
						return errors.New("同一个账号不可以连续竞价")
					}
				}
			}
		}

		// 每次至少增加 1%
		if itemCount < int64(lastItemCounts*1.1) {
			return errors.New("至少增加1%进行竞价")
		}

		nA := &AuctionLog{AccName: accName, ItemCount: itemCount, Time: time.Now().Unix()}

		j.Sessions[j.LastSessionID].LastCounts = itemCount
		j.Sessions[j.LastSessionID].Logs = append(j.Sessions[j.LastSessionID].Logs, nA)

		// 玩家无法赖账, 因为通货有多个入口(+), 但是只有拍卖场一个出口(而且只有一个拍卖机会-当前场次), 过了当前场次, 竞价作废, 重新再来
		return nil
	}

	return errors.New("没有这种通货拍卖场")
}
