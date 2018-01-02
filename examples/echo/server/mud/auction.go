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

// 玩家拍卖时冻结指定的货币, ItemID, ItemCount, 拍卖结束解冻, 可能使用, 也可能兑换了Jewel
// 玩家宝石也有冻结, 活动两种
// 玩家宝石有利息

// 全服只有一个竞拍, 全服共享竞拍
//

// // 请求冻结宝石
// func (k *KuangSys) AskFreezeJewel(name string, count int64) bool {
// 	if _, ok := k.Kuangs[name]; ok {
// 		// 已经冻结宝石(正在拍卖)
// 		if k.Kuangs[name].FreezeJewel > 0 {
// 			return false
// 		}
// 		// 宝石不够数量
// 		if k.Kuangs[name].Jewel < count {
// 			return false
// 		}
// 		k.Kuangs[name].FreezeJewel = count
// 		k.Kuangs[name].Jewel = k.Kuangs[name].Jewel - count
// 	}
// }

// // 请求解冻宝石
// func (k *KuangSys) AskUnfreezeJewel(name string, count int64) bool {
// 	if _, ok := k.Kuangs[name]; ok {
// 		// 已经冻结宝石(正在拍卖)
// 		if k.Kuangs[name].FreezeJewel != count {
// 			return false
// 		}
// 		k.Kuangs[name].FreezeJewel = 0
// 		k.Kuangs[name].Jewel = k.Kuangs[name].Jewel + count
// 	}
// }

// // 请求兑换冻结的宝石
// func (k *KuangSys) AskExchangeFreezeJewel(name string, count int64, accName string) bool {
// 	if _, ok := k.Kuangs[name]; ok {
// 		// 已经冻结宝石(正在拍卖)
// 		if k.Kuangs[name].FreezeJewel < count {
// 			return false
// 		}
// 		k.Kuangs[name].FreezeJewel -= count
// 		// 宝石邮件投递给 accName, 对方扣除冻结的游戏币
// 		return true
// 	}
// 	return false
// }
