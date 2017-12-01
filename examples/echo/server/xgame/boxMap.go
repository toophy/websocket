package xgame

import (
	"sort"
)

// BoxMap 盒子地图
type BoxMap struct {
	Box                       // 盒子
	Objs       map[int32]*Obj // 地图所有对象
	LastObjID  int32          // 最后一个对象ID
	FrameSN    int32          // 帧序号
	GravityVal int32          // 本地图基本重力
}

// InitMap 初始化盒子地图
func (b *BoxMap) InitMap(tileW, tileH, realW, realH, gravity int32) bool {
	b.Objs = make(map[int32]*Obj, 0)
	b.LastObjID = 0
	b.FrameSN = 0
	b.GravityVal = gravity
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

// ObjMove 对象在Box上移动,step是步长
func (b *BoxMap) ObjMove(o *Obj, dir int, stepX int32) bool {
	if o == nil {
		return false
	}
	newRect := o.Pos
	switch dir {
	case DirLeft:
		newRect.X -= stepX
		break
	case DirRight:
		newRect.X += stepX
		break
	case DirUp:
		newRect.Y -= stepX
		break
	case DirDown:
		newRect.Y += stepX
		break
	}
	return b.Insert(o, &newRect)
}

// ObjMove 对象在Box上移动,step是步长
func (b *BoxMap) ObjJump(o *Obj, speedY int32) bool {
	if o == nil || o.Jump {
		return false
	}

	o.SpeedY = speedY
	o.SpeedX = o.GetRealSpeedX()
	o.Jump = true

	if o.Parent != nil {
		o.Parent.EraseChild(o.ID)
		o.Parent = nil
	}

	// 落地之前, 玩家都不能操作,
	return true
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
		if o.Parent == nil {
			b.updateObjMove(o)
		}
	}

	for _, k := range keys {
		o := b.Objs[int32(k)]
		if o.Parent == nil {
			b.updateObjMoveEvent(o)
		}
	}

	// 可能有对象脱离Parent, 看看是否到达新Parent,
	// 对象类别
	// 0 : 可以容纳其他对象, 不受重力影响 --> 陆地块
	// 1 : 不可以容纳其他对象, 受重力影响 --> 角色怪物

	// 也就是parent为nil的类型1对象, 检查是否有新的类型0对象降落
}

// updateObjMove 刷新对象移动
func (b *BoxMap) updateObjMove(o *Obj) {
	speedX := o.SpeedX
	speedY := o.SpeedY

	newRect := o.Pos
	newRect.X += speedX

	if o.Gravity && o.Jump {
		o.SpeedY += b.GravityVal
		speedY += b.GravityVal
		newRect.Y += speedY
	}
	if newRect.X != o.Pos.X || newRect.Y != o.Pos.Y {
		o.Moved = true
	}

	if b.Box.Insert(o, &newRect) {
		if o.Contain {
			var keys []int
			for k := range o.Childs {
				keys = append(keys, int(k))
			}
			sort.Ints(keys)

			for _, k := range keys {
				c := o.Childs[int32(k)]
				b.updateObjMove(c)
			}
		}
	} else {
		// 可能发生碰撞
		//

	}
}

// updateObjMoveEvent 刷新对象新位置事件
func (b *BoxMap) updateObjMoveEvent(o *Obj) {

	// 新位置是否有事件触发?
	if o.Moved {
		// 触发事件(切换parent,拾取到xx,陷阱xx)
		if o.Parent != nil {
			if o.Pos.X > o.Parent.Pos.W || o.Pos.X+o.Pos.W < 0 {
				o.Parent.EraseChild(o.ID)
				o.Parent = nil
			}
		}
		o.Moved = false
	}

	var keys []int
	for k := range o.Childs {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)
	for _, k := range keys {
		c := o.Childs[int32(k)]
		b.updateObjMoveEvent(c)
	}
}
