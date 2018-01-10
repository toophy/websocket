package mud

import (
	"errors"
)

// Goods 商品
type Goods struct {
	Item        // 道具
	Price int32 // 单价
	Money int32 // 货币ItemType
}

// Shop 商铺
type Shop struct {
	AccID   int64              // 帐号ID
	Goodses map[int32][]*Goods // 用ItemType索引所有道具
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
			for k, v := range d.Goodses[i.Type] {
				if v.Count == 0 {
					i.Count = 0
					d.Goodses[i.Type] = append(d.Goodses[i.Type], i)
					break
				} else {
					v.Count += i.Count
					break
				}
			}
		} else {
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
		for k, v := range d.Goodses[typeID] {
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
		realCount := len(d.Goodses[typeID])
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
				append(d.Goodses[typeID][0:k], d.Goodses[typeID][k+1:]...)
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
		for k, v := range d.Goodses[typeID] {
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
		for k, v := range d.Goodses[typeID] {
			if v.Count == 0 {
				ret = len(d.Goodses[typeID])
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
