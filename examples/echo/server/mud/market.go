package mud

import (
	"fmt"
)

// 市场则是 管理 商铺
// 搜索 商品
// 购买 商品

// 商铺应该是在另一个地方, GetShopSys()
// 市场远程调用,
//

// Market 市场
type Market struct {
}

var (
	market *Market
)

func init() {
	market = &Market{}
}

// GetShopSys 获取商铺系统
func GetMarket() *ShopSys {
	return market
}

// SearchGoods 搜索商品
func (m *Market) SearchGoods(typeID int32, moneyType int32, moneyMax int32) {
	go GetShopSys().OnSearchGoods(typeID, moneyType, moneyMax)
}

// OnRetSearchGoods 反馈搜索商品
func (m *Market) OnRetSearchGoods(g []*Goods) {
	fmt.Print("反馈搜索商品: \n")
	for i := 0; i < len(g); i++ {
		fmt.Printf("\t%-v\n", g[i])
	}
	fmt.Printf("\n")
}

// BuyGoods 购买商品
func (m *Market) BuyGoods(accID int64, goodsID int32, typeID int32, buyCount int32, moneyType int32) {
	go GetShopSys().OnBuyGoods(accID, goodsID, typeID, buyCount, moneyType)
}

// OnRetBuyGoods 反馈购买商品结果
func (m *Market) OnRetBuyGoods(accID int64, goodsID int32, typeID int32, buyCount int32, moneyType int32, ret int32, msg string) {
	fmt.Printf("帐号%d使用%d购买%d个%d %s\n", accID, goodsID, buyCount, typeID, msg)
}
