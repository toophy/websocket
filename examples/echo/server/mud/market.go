package mud

// 市场则是 管理 商铺
// 搜索 商品
// 购买 商品

// 商铺应该是在另一个地方, GetShopSys()
// 市场远程调用,
//

// Market 市场
type Market struct {
}

// SearchGoods 搜索商品
func (m *Market) SearchGoods(typeID int32) {
	// 让 ShopSys 搜索并反馈, 提供反馈函数
}
