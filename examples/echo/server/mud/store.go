package mud

import (
	"errors"
)

// Store 仓库
type Store struct {
	AccID int64             // 帐号ID
	Items map[int32][]*Item // 用ItemType索引所有道具
}

// Insert 插入道具
func (d *Store) Insert(i *Item) bool {
	if i != nil {
		if _, ok := d.Items[i.Type]; ok {
			for k, v := range d.Items[i.Type] {
				if v.Count == 0 {
					i.Count = 0
					d.Items[i.Type] = append(d.Items[i.Type], i)
					break
				} else {
					v.Count += i.Count
					break
				}
			}
		} else {
			d.Items[i.Type] = make([]*Item, 0)
			d.Items[i.Type] = append(d.Items[i.Type], i)
		}
		return true
	}
	return false
}

// Delete 删除道具
func (d *Store) DeleteByType(typeID int32, count int32) error {
	if _, ok := d.Items[typeID]; ok {
		for k, v := range d.Items[typeID] {
			if v.Count == 0 {
				break
			} else {
				if v.Count > count {
					v.Count -= count
					return nil
				} else if v.Count == count {
					delete(d.Items, typeID)
					return nil
				} else {
					return errors.New("没有足够道具被删除")
				}
			}
		}
		// 一个个顺序删除?
		realCount := len(d.Items[typeID])
		if realCount > count {
			d.Items[typeID] = d.Items[typeID][count:realCount]
		} else if realCount == count {
			delete(d.Items, typeID)
		}
		return nil
	}
	return errors.New("你没有这种道具")
}

// DeleteByID 删除指定道具ID的物品
func (d *Store) DeleteByID(typeID, itemID int32) error {
	if itemID <= 0 {
		return errors.New("非法道具ID")
	}

	if _, ok := d.Items[typeID]; ok {
		for k, v := range d.Items[typeID] {
			if v.ID == itemID {
				append(d.Items[typeID][0:k], d.Items[typeID][k+1:]...)
				return nil
			}
		}
		return errors.New("你没有这个ID道具")
	}
	return errors.New("你没有这种道具")
}

// GetItemByID 获取指定道具ID的物品
func (d *Store) GetItemByID(typeID, itemID int32) (*Item, error) {
	if itemID <= 0 {
		return nil, errors.New("非法道具ID")
	}

	if _, ok := d.Items[typeID]; ok {
		for k, v := range d.Items[typeID] {
			if v.ID == itemID {
				return v, nil
			}
		}
		return nil, errors.New("你没有这个ID道具")
	}
	return nil, errors.New("你没有这种道具")
}

// GetItemCountByType 通货道具类型获取总量
func (d *Store) GetItemCountByType(typeID int32) (ret int32) {
	if _, ok := d.Items[typeID]; ok {
		for k, v := range d.Items[typeID] {
			if v.Count == 0 {
				ret = len(d.Items[typeID])
				return
			} else {
				ret += v.Count
			}
		}
	}
	return
}
