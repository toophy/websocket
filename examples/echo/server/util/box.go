package util

// 方向定义
const (
	DirLeft  = 1 // 向左
	DirRight = 2 // 向右
	DirUp    = 3 // 向上
	DirDown  = 4 // 向下
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
	Terre   bool           // 紧贴地面
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

// Box 染色 tile, 一个物品可以挂在多个tile下面, 用来进行碰撞抽检, 用ID作为碰撞顺序标准
type Box struct {
	initOk    bool
	W         int32
	H         int32
	TileW     int32
	TileH     int32
	RealW     int32
	RealH     int32
	CellCount uint32
	Cells     []Cell
}

// Init 初始化Box
func (b *Box) Init(tileW, tileH, realW, realH int32) bool {

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

	b.initOk = true

	b.Clear()

	return true
}

// GetCrossCells 获取矩形交叉的Cell
func (b *Box) GetCrossCells(x, y, w, h int32) (rets []int32) {
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
func (b *Box) CanInsert(x, y, w, h int32, o *Obj) bool {
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
func (b *Box) Insert(o *Obj, newRect *Rect) bool {
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
func (b *Box) Erase(o *Obj) {
	// 褪色
	if o != nil {
		for _, v := range o.Cells {
			delete(b.Cells[v].Objs, o.ID)
		}
	}
	o.Cells = []int32{}
}

// Clear 清除所有对象痕迹
func (b *Box) Clear() {
	if b.initOk {
		for i := uint32(0); i < b.CellCount; i++ {
			b.Cells[i].Objs = make(map[int32]*Obj, 0)
		}
	}
}

// Destroy 摧毁Box
func (b *Box) Destroy() {
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
