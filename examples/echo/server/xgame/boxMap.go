package xgame

import (
	"sort"
)

// 方向定义
const (
	DirLeft  = 1 // 向左
	DirRight = 2 // 向右
	DirUp    = 3 // 向上
	DirDown  = 4 // 向下
)

// 玩家操作
const (
	POS_Move = 1 // 移动
)

// Rect 矩形
type Rect struct {
	X int32
	Y int32
	W int32 // 最小值 1
	H int32 // 最小值 1
}

// CrossEx 碰撞
func (r *Rect) CrossEx(a Rect) bool {
	return r.Cross(a.X, a.Y, a.W, a.H)
}

// Cross 碰撞
func (r *Rect) Cross(x, y, w, h int32) bool {
	if x+w-1 < r.X {
		return false
	}
	if x > r.X+r.W-1 {
		return false
	}
	if y+h-1 < r.Y {
		return false
	}
	if y > r.Y+r.H-1 {
		return false
	}
	return true
}

// Obj 地图对象
type Obj struct {
	ID      int32          // ID
	Contain bool           // 能够容纳其他对象
	Gravity bool           // 受重力影响
	Jump    bool           // 跳动
	Pos     Rect           // 位置和体量
	Moved   bool           // 本帧发生移动
	SpeedX  int32          // 水平速度
	SpeedY  int32          // 垂直速度
	Cells   []int32        // 最后一次染色单元格
	Parent  *Obj           // 父对象
	Childs  map[int32]*Obj // 子对象
}

// 擦除一个子对象
func (o *Obj) EraseChild(cID int32) {
	delete(o.Childs, cID)
}

// 获取实际水平速度
func (o *Obj) GetRealSpeedX() int32 {
	if o.Parent == nil {
		return o.SpeedX
	}
	return o.SpeedX + o.Parent.GetRealSpeedX()
}

// 获取实际垂直速度
func (o *Obj) GetRealSpeedY() int32 {
	if o.Parent == nil {
		return o.SpeedY
	}
	return o.SpeedY + o.Parent.GetRealSpeedY()
}

// Cell 地图单元格
type Cell struct {
	Objs map[int32]*Obj
}

// FrameInput 玩家帧输入
type FrameInput struct {
	ID       int32 // 玩家ID
	OpID     int32 // 玩家操作
	OpParam1 int32 // 操作参数1
	OpParam2 int32 // 操作参数2
}

// FrameData 帧输入
type FrameData struct {
	Sn  int32         // 帧序号
	Ops []*FrameInput // 玩家帧输入
}

// BoxMap 盒子地图
type BoxMap struct {
	initOk     bool
	W          int32
	H          int32
	TileW      int32
	TileH      int32
	RealW      int32
	RealH      int32
	CellCount  uint32
	Cells      []Cell
	Objs       map[int32]*Obj       // 地图所有对象
	LastObjID  int32                // 最后一个对象ID
	GravityVal int32                // 本地图基本重力
	FrameDatas map[int32]*FrameData // 帧数据
}

// InitMap 初始化盒子地图
func (b *BoxMap) InitMap(tileW, tileH, realW, realH, gravity int32) bool {
	b.Objs = make(map[int32]*Obj, 0)
	b.LastObjID = 0
	b.GravityVal = gravity

	if tileW < 1 || tileH < 1 || realW < 0 || realH < 0 || tileW > realW || tileH > realH {
		return false
	}

	b.Destroy()
	b.RealW = realW
	b.RealH = realH
	b.TileW = tileW
	b.TileH = tileH
	b.W = (realW + tileW - 1) / tileW
	b.H = (realH + tileH - 1) / tileH

	b.CellCount = uint32(b.W * b.H)
	b.Cells = make([]Cell, b.CellCount)

	b.FrameDatas = make(map[int32]*FrameData, 0)

	b.initOk = true

	b.Clear()

	return true
}

// GetCrossCells 获取矩形交叉的Cell
func (b *BoxMap) GetCrossCells(x, y, w, h int32) (rets []int32) {
	if x < 0 || y < 0 || w < 0 || h < 0 || x+w > b.RealW || y+h > b.RealH {
		return []int32{}
	}
	x2 := x + w
	y2 := y + h
	cx := x / b.TileW
	cy := y / b.TileH
	cx2 := x2 / b.TileW
	cy2 := y2 / b.TileH
	for r := cy; r <= cy2; r++ {
		for c := cx; c <= cx2; c++ {
			rets = append(rets, r*b.W+c)
		}
	}
	return
}

// CanInsert 检查对象使用新矩形能否插入
func (b *BoxMap) CanInsert(x, y, w, h int32, o *Obj) bool {
	if o == nil {
		return false
	}

	newCells := b.GetCrossCells(x, y, w, h)
	if len(newCells) == 0 {
		return false
	}

	for _, v := range newCells {
		for _, cv := range b.Cells[v].Objs {
			if o != cv && cv.Pos.Cross(x, y, w, h) {
				return false
			}
		}
	}

	return true
}

// Insert 进行插入
func (b *BoxMap) Insert(o *Obj, newRect *Rect) bool {
	if o != nil {
		tmpRect := o.Pos
		if newRect != nil {
			tmpRect = *newRect
		}
		if b.CanInsert(tmpRect.X, tmpRect.Y, tmpRect.W, tmpRect.H, o) {
			o.Pos = tmpRect
			// 褪色
			for _, v := range o.Cells {
				delete(b.Cells[v].Objs, o.ID)
			}
			// 染色
			newCells := b.GetCrossCells(tmpRect.X, tmpRect.Y, tmpRect.W, tmpRect.H)
			for _, v := range newCells {
				b.Cells[v].Objs[o.ID] = o
			}
			o.Cells = newCells
			return true
		}
	}
	return false
}

// Erase 擦除一个对象
func (b *BoxMap) Erase(o *Obj) {
	if o != nil {
		for _, v := range o.Cells {
			delete(b.Cells[v].Objs, o.ID)
		}
		o.Cells = []int32{}
	}
}

// Clear 清除所有对象痕迹
func (b *BoxMap) Clear() {
	if b.initOk {
		for i := uint32(0); i < b.CellCount; i++ {
			b.Cells[i].Objs = make(map[int32]*Obj, 0)
		}
	}
}

// Destroy 摧毁Box
func (b *BoxMap) Destroy() {
	if b.initOk {
		b.W = 0
		b.H = 0
		b.TileW = 0
		b.TileH = 0
		b.RealW = 0
		b.RealH = 0
		b.Cells = []Cell{}
		b.CellCount = 0
		b.initOk = false
	}
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

// 玩家输入操作
func (b *BoxMap) InputOp(f *FrameInput) {
	if f != nil {
		switch f.OpID {
		case POS_Move:
			o := b.GetObjByID(f.ID)
			if o != nil {
				b.ObjMove(o, int(f.OpParam1), f.OpParam2)
			}
			break
		}
	}
}

// Update 刷新BoxMap
func (b *BoxMap) Update(frameSn int32) {
	var keys []int
	for k := range b.Objs {
		keys = append(keys, int(k))
	}
	sort.Ints(keys)

	for _, v := range keys {
		o := b.Objs[int32(v)]
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
	if _, ok := b.FrameDatas[frameSn]; ok {
		for k, _ := range b.FrameDatas[frameSn].Ops {
			b.InputOp(&b.FrameDatas[frameSn].Ops[k])
		}
	}
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

	if b.Insert(o, &newRect) {
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
