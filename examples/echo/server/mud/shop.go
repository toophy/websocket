package mud

import (
	"errors"
)

// Goods 商品
type Goods struct {
	Item            // 道具
	GoodsID   int32 // 商品ID
	Price     int32 // 单价
	MoneyType int32 // 货币ItemType
}

// Goods 商品
type RetGoods struct {
	Goods       // 商品
	AccID int64 //帐号
}

// Shop 商铺
type Shop struct {
	AccID       int64              // 帐号ID
	LastGoodsID int32              // 最后一个商品ID
	Goodses     map[int32][]*Goods // 用ItemType索引所有道具
}

// Init 初始化商铺
func (s *Shop) Init(accID int64) {
	s.AccID = accID
	s.Goodses = make(map[int32][]*Goods, 0)
}

// Insert 插入道具
func (d *Shop) Insert(i *Goods) bool {
	if i != nil {
		if _, ok := d.Goodses[i.Type]; ok {
			for _, v := range d.Goodses[i.Type] {
				if v.Count == 0 {
					d.LastGoodsID++
					i.Count = 0
					i.GoodsID = d.LastGoodsID
					d.Goodses[i.Type] = append(d.Goodses[i.Type], i)
					break
				} else {
					v.Count += i.Count
					break
				}
			}
		} else {
			d.LastGoodsID++
			i.GoodsID = d.LastGoodsID
			d.Goodses[i.Type] = make([]*Goods, 0)
			d.Goodses[i.Type] = append(d.Goodses[i.Type], i)
		}
		return true
	}
	return false
}

// Delete 删除道具
func (d *Shop) DeleteByType(typeID int32, count int32) error {
	if _, ok := d.Goodses[typeID]; ok {
		for _, v := range d.Goodses[typeID] {
			if v.Count == 0 {
				break
			} else {
				if v.Count > count {
					v.Count -= count
					return nil
				} else if v.Count == count {
					delete(d.Goodses, typeID)
					return nil
				} else {
					return errors.New("没有足够道具被删除")
				}
			}
		}
		// 一个个顺序删除?
		realCount := int32(len(d.Goodses[typeID]))
		if realCount > count {
			d.Goodses[typeID] = d.Goodses[typeID][count:realCount]
		} else if realCount == count {
			delete(d.Goodses, typeID)
		}
		return nil
	}
	return errors.New("你没有这种道具")
}

// DeleteByID 删除指定道具ID的物品
func (d *Shop) DeleteByID(typeID, itemID int32) error {
	if itemID <= 0 {
		return errors.New("非法道具ID")
	}

	if _, ok := d.Goodses[typeID]; ok {
		for k, v := range d.Goodses[typeID] {
			if v.ID == itemID {
				d.Goodses[typeID] = append(d.Goodses[typeID][0:k], d.Goodses[typeID][k+1:]...)
				return nil
			}
		}
		return errors.New("你没有这个ID道具")
	}
	return errors.New("你没有这种道具")
}

// GetItemByID 获取指定道具ID的物品
func (d *Shop) GetItemByID(typeID, itemID int32) (*Goods, error) {
	if itemID <= 0 {
		return nil, errors.New("非法道具ID")
	}

	if _, ok := d.Goodses[typeID]; ok {
		for _, v := range d.Goodses[typeID] {
			if v.ID == itemID {
				return v, nil
			}
		}
		return nil, errors.New("你没有这个ID道具")
	}
	return nil, errors.New("你没有这种道具")
}

// GetItemCountByType 通货道具类型获取总量
func (d *Shop) GetItemCountByType(typeID int32) (ret int32) {
	if _, ok := d.Goodses[typeID]; ok {
		for _, v := range d.Goodses[typeID] {
			if v.Count == 0 {
				ret = int32(len(d.Goodses[typeID]))
				return
			} else {
				ret += v.Count
			}
		}
	}
	return
}

// ShopSys 商铺系统
type ShopSys struct {
	Shops map[int64]*Shop
}

var (
	shopSys *ShopSys
)

func init() {
	shopSys = &ShopSys{}
	shopSys.Shops = make(map[int64]*Shop, 0)
}

// GetShopSys 获取商铺系统
func GetShopSys() *ShopSys {
	return shopSys
}

// OnSearchGoods 响应市场搜索商品
func (s *ShopSys) OnSearchGoods(typeID int32, moneyType int32, moneyMax int32) {

	g := make([]*RetGoods, 0)

	for _, vs := range s.Shops {
		if _, ok := vs.Goodses[typeID]; ok {
			for _, v := range vs.Goodses[typeID] {
				if v.MoneyType == moneyType && v.Price <= moneyMax {
					c := &RetGoods{*v, vs.AccID}
					g = append(g, c)
				}
			}
		}
	}

	go GetMarket().OnRetSearchGoods(g)
}

// OnBuyGoods 响应购买商品
func (s *ShopSys) OnBuyGoods(accID int64, goodsID int32, typeID int32, buyCount int32, moneyType int32) {
	ret := int32(0)
	msg := "购买成功"

	if _, ok := s.Shops[accID]; ok {
		vs := s.Shops[accID]
		if _, ok := vs.Goodses[typeID]; ok {
			isFind := false
			for _, v := range vs.Goodses[typeID] {
				if v.GoodsID == goodsID {
					isFind = true
					if v.MoneyType == moneyType && buyCount <= v.Count {
						// 成功购买, 发送系统脚本邮件给用户
						// 同时去掉商铺中的商品?
						// 什么时候收到玩家的付款呢?
						// 脚本邮件成功使用后, 玩家付款给商铺?
					} else {
						ret = 1
						msg = "通货不匹配或者购买量不合理"
					}
					break
				}
			}
			if !isFind {
				ret = 2
				msg = "指定商品不存在"
			}
		} else {
			ret = 3
			msg = "商品没有上架"
		}
	} else {
		ret = 4
		msg = "商铺不存在"
	}

	go GetMarket().OnRetBuyGoods(accID, goodsID, typeID, buyCount, moneyType, ret, msg)
}
