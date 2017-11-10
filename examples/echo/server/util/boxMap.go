package util

import (
	"sort"
)

// BoxMap 盒子地图
type BoxMap struct {
	Box                      // 盒子
	Objs      map[int32]*Obj // 地图所有对象
	LastObjID int32          // 最后一个对象ID
	FrameSN   int32          // 帧序号
}

// InitMap 初始化盒子地图
func (b *BoxMap) InitMap(tileW, tileH, realW, realH int32) bool {
	b.Objs = make(map[int32]*Obj, 0)
	return b.Init(tileW, tileH, realW, realH)
}

// NewObj 新建对象
func (b *BoxMap) NewObj(x, y, w, h int32) *Obj {
	b.LastObjID++
	n := &Obj{ID: b.LastObjID, Pos: Rect{X: x, Y: y, W: w, H: h}, Cells: []int32{}}
	b.Objs[b.LastObjID] = n
	return n
}

// GetObjByID 从对象ID获取对象数据
func (b *BoxMap) GetObjByID(id int32) *Obj {
	if _, ok := b.Objs[id]; ok {
		return b.Objs[id]
	}
	return nil
}

// Update 刷新BoxMap
func (b *BoxMap) Update() {
	b.FrameSN++

	var keys []int
	for k := range b.Objs {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, k := range keys {
		o := b.Objs[int32(k)]
		b.updateObj(o)
	}
}

// UpdateObj 刷新对象
func (b *BoxMap) updateObj(o *Obj) {
	// 对象内部呢, 如果一次"行走"连续多个点, 就调用多次
	// 对象瞬移, 只需要调用一次(最终位置帧)
	// 
	// 1. 移动, 速度, 开始....帧
}
