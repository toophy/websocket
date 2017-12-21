package mud

import (
	"time"
)

const (
	RT_None   = int32(0) // 没有关系
	RT_Friend = int32(1) // 好友
	RT_Black  = int32(2) // 黑名单
)

// Friend 好友
type Friend struct {
	ID        int64  // ID
	Name      string // 名称
	StartTime int32  // 开始关系的时间
	Relation  int32  // 关系
}

// FriendSys 朋友之间做的事情
type FriendSys struct {
	Friends map[int64]*Friend
}

// Add 新增好友
func (f *FriendSys) Add(id int64, name string) bool {
	if _, ok := f.Friends[id]; !ok {
		f.Friends[id] = &Friend{id, name, int32(time.Now().Unix()), RT_Friend}
		return true
	}
	return false
}

// Del 删除好友
func (f *FriendSys) Del(id int64) bool {
	delete(f.Friends, id)
	return true
}

// Has 有relation关系
func (f *FriendSys) Has(id int64, relation int32) bool {
	if _, ok := f.Friends[id]; ok {
		if relation == RT_None {
			return true
		}
		return f.Friends[id].Relation == relation
	}
	return false
}

// ToRelation 改变关系
func (f *FriendSys) ToRelation(id int64, relation int32) bool {
	if relation == RT_None {
		return f.Del(id)
	}

	if relation != RT_Friend || relation != RT_Black {
		return false
	}

	if _, ok := f.Friends[id]; ok {
		if f.Friends[id].Relation != relation {
			f.Friends[id].Relation = relation
			return true
		}
	}
	return false
}
