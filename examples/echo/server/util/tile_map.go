package util

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
	W int32
	H int32
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
	Id    int32
	Pos   Rect
	Cells []int32 // 最后一次染色单元格
}

type Cell struct {
	Objs map[int32]*Obj
}

// 染色 tile, 一个物品可以挂在多个tile下面, 用来进行碰撞抽检, 用ID作为碰撞顺序标准
type TileMap struct {
	initOk    bool
	W         int32
	H         int32
	tileW     int32
	tileH     int32
	realW     int32
	realH     int32
	CellCount uint32
	Cells     []Cell
}

// Init 初始化TileMap
func (q *TileMap) Init(tileW, tileH, realW, realH int32) bool {

	if tileW < 1 || tileH < 1 || realW < 0 || realH < 0 || tileW > realW || tileH > realH {
		return false
	}

	q.Destroy()
	q.realW = realW
	q.realH = realH
	q.tileW = tileW
	q.tileH = tileH
	q.W = (realW + tileW - 1) / tileW
	q.H = (realH + tileH - 1) / tileH

	q.CellCount = uint32(q.W * q.H)
	q.Cells = make([]Cell, q.CellCount)

	q.initOk = true

	q.Clear()

	return true
}

// GetCrossCells 获取矩形交叉的Cell
func (q *TileMap) GetCrossCells(x, y, w, h int32) (rets []int32) {
	if x < 0 || y < 0 || w < 0 || h < 0 || x+w > q.realW || y+h > q.realH {
		return []int32{}
	}
	x2 := x + w
	y2 := y + h
	cx := x / q.tileW
	cy := y / q.tileH
	cx2 := x2 / q.tileW
	cy2 := y2 / q.tileH
	for r := cy; r <= cy2; r++ {
		for c := cx; c <= cx2; c++ {
			rets = append(rets, r*q.W+c)
		}
	}
	return
}

// CanInsert 检查对象使用新矩形能否插入
func (q *TileMap) CanInsert(x, y, w, h int32, o *Obj) bool {
	if o == nil {
		return false
	}

	newCells := q.GetCrossCells(x, y, w, h)
	if len(newCells) == 0 {
		return false
	}

	for _, v := range newCells {
		for _, cv := range q.Cells[v].Objs {
			if o != cv && cv.Pos.Cross(x, y, w, h) {
				return false
			}
		}
	}

	return true
}

// Insert 进行插入
func (q *TileMap) Insert(o *Obj, newRect *Rect) bool {
	if o != nil {
		if newRect != nil {
			o.Pos = *newRect
		}
		if q.CanInsert(o.Pos.X, o.Pos.Y, o.Pos.W, o.Pos.H, o) {
			// 褪色
			for _, v := range o.Cells {
				delete(q.Cells[v].Objs, o.Id)
			}
			// 染色
			newCells := q.GetCrossCells(o.Pos.X, o.Pos.Y, o.Pos.W, o.Pos.H)
			for _, v := range newCells {
				q.Cells[v].Objs[o.Id] = o
			}
			o.Cells = newCells
			return true
		}
	}
	return false
}

// Erase 擦除一个对象
func (q *TileMap) Erase(o *Obj) {
	// 褪色
	if o != nil {
		for _, v := range o.Cells {
			delete(q.Cells[v].Objs, o.Id)
		}
	}
	o.Cells = []int32{}
}

// Clear 清除所有对象痕迹
func (q *TileMap) Clear() {
	if q.initOk {
		for i := uint32(0); i < q.CellCount; i++ {
			q.Cells[i].Objs = make(map[int32]*Obj, 0)
		}
	}
}

// Destroy 摧毁TileMap
func (q *TileMap) Destroy() {
	if q.initOk {
		q.W = 0
		q.H = 0
		q.tileW = 0
		q.tileH = 0
		q.realW = 0
		q.realH = 0
		q.Cells = []Cell{}
		q.CellCount = 0
		q.initOk = false
	}
}

// ObjMove 对象在TileMap上移动,step是步长
func (q *TileMap) ObjMove(o *Obj, dir int, step int32) bool {
	if o == nil {
		return false
	}
	newRect := o.Pos
	switch dir {
	case DirLeft:
		newRect.X -= step
		break
	case DirRight:
		newRect.X += step
		break
	case DirUp:
		newRect.Y -= step
		break
	case DirDown:
		newRect.Y += step
		break
	}
	return q.Insert(o, &newRect)
}
