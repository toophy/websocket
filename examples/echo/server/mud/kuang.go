package mud

import (
	"math"
	"time"
)

// Kuang 矿
type Kuang struct {
	Name          string  // 矿名
	Ore           float64 // 原矿石
	Jewel         int64   // 宝石, 从矿石中提取整数部分
	OreSpeed      float64 // 原矿石生产速度 Ore/秒
	Lucky         float64 // 幸运值
	LastJewelTime int64   // 最后一次宝石计算时间
}

// 每秒一次计算, 每分钟一次计算
func (k *Kuang) Update() {
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
}

// 很多种矿
type KuangSys struct {
	Kuangs map[string]*Kuang
}

var (
	kuangSys *KuangSys
)

func init() {
	kuangSys = &KuangSys{}
	kuangSys.Kuangs = make(map[string]*Kuang, 0)
}

func GetKuangSys() *KuangSys {
	return kuangSys
}

// InitSys 初始化矿
func (k *KuangSys) InitSys() {
	k.Kuangs = make(map[string]*Kuang, 0)
}

// CreateKuang 新搭建一个矿场
func (k *KuangSys) CreateKuang(name string, ore float64, jewel int64, oreSpeed float64, lucky float64, lastJewelTime int64) bool {
	if lastJewelTime == 0 {
		lastJewelTime = time.Now().Unix()
	}
	if _, ok := k.Kuangs[name]; !ok {
		k.Kuangs[name] = &Kuang{Name: name, Ore: ore, Jewel: jewel, OreSpeed: oreSpeed, Lucky: lucky, LastJewelTime: lastJewelTime}
		return true
	}
	return false
}

// 每秒一次计算, 每分钟一次计算
func (k *KuangSys) Update() {
	t := time.NewTimer(1 * time.Second)

	for {
		select {
		case <-t.C:
			for i := range k.Kuangs {
				k.Kuangs[i].Update()
			}
			t.Reset(1 * time.Second)
			break
		}
	}
}

// 请求冻结宝石
func (k *KuangSys) AskFreezeJewel(name string, count int64) bool {
	if _, ok := k.Kuangs[name]; ok {
		// 已经冻结宝石(正在拍卖)
		if k.Kuangs[name].FreezeJewel > 0 {
			return false
		}
		// 宝石不够数量
		if k.Kuangs[name].Jewel < count {
			return false
		}
		k.Kuangs[name].FreezeJewel = count
		k.Kuangs[name].Jewel = k.Kuangs[name].Jewel - count
	}
}

// 请求解冻宝石
func (k *KuangSys) AskUnfreezeJewel(name string, count int64) bool {
	if _, ok := k.Kuangs[name]; ok {
		// 已经冻结宝石(正在拍卖)
		if k.Kuangs[name].FreezeJewel != count {
			return false
		}
		k.Kuangs[name].FreezeJewel = 0
		k.Kuangs[name].Jewel = k.Kuangs[name].Jewel + count
	}
}

// 请求兑换冻结的宝石
func (k *KuangSys) AskExchangeFreezeJewel(name string, count int64, accName string) bool {
	if _, ok := k.Kuangs[name]; ok {
		// 已经冻结宝石(正在拍卖)
		if k.Kuangs[name].FreezeJewel < count {
			return false
		}
		k.Kuangs[name].FreezeJewel -= count
		// 宝石邮件投递给 accName, 对方扣除冻结的游戏币
		return true
	}
	return false
}

// 宝石拍卖
// 需要报宝石拿出来另外保存, 进行拍卖
// 每次拍卖时间短(原始5分钟拍卖时间), 有新出价就延迟1分钟, 否则结束拍卖
