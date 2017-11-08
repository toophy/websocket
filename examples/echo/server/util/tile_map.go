package util

type Rect struct {
	X int32
	Y int32
	W int32
	H int32
}

// 碰撞
func (r *Rect) CrossEx(a Rect) bool {
	return r.Cross(a.X, a.Y, a.W, a.H)
}

// 碰撞
func (r *Rect) Cross(x, y, w, h int32) bool {

	aMaxX := x + w
	aMaxY := y + h
	bMaxX := r.X + r.W
	bMaxY := r.Y + r.H

	// a is left of b
	if aMaxX < r.X {
		return false
	}

	// a is right of b
	if x > bMaxX {
		return false
	}

	// a is above b
	if aMaxY < r.Y {
		return false
	}

	// a is below b
	if y > bMaxY {
		return false
	}

	// The two overlap
	return true
}

type Obj struct {
	Id    int16
	Pos   Rect
	Cells []int32 // 最后一次染色单元格
}

type Cell struct {
	Objs map[int16]*Obj
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

// 由具体的坐标获取cell坐标
func (q *TileMap) GetCellId(x, y int32) (int32, int32, bool) {
	if x < 0 || y < 0 || x > q.realW || y > q.realH {
		return 0, 0, false
	}
	return x / q.tileW, y / q.tileH, true
}

func (q *TileMap) GetPointCell(x, y int32) int32 {
	if x < 0 || y < 0 || x > q.realW || y > q.realH {
		return -1
	}
	return y/q.tileH*q.W + x/q.tileW
}

func (q *TileMap) GetRectCell(x, y, w, h int32) []int32 {
	if x < 0 || y < 0 || w < 0 || h < 0 || x+w >= q.realW || y+h >= q.realH {
		return []int32{}
	}
	rets := []int32{}
	x2 := x + w
	y2 := y + h
	cx := x / q.tileW
	cy := y / q.tileH
	cx2 := x2 / q.tileW
	cy2 := y2 / q.tileH
	for r := cy; r < cy2; r++ {
		for c := cx; c <= cx2; c++ {
			rets = append(rets, r*q.W+c)
		}
	}
	return rets
}

func (q *TileMap) CanInsert(x, y, w, h int32) bool {
	newCells := q.GetRectCell(x, y, w, h)
	for _, v := range newCells {
		for _, cv := range q.Cells[v-1].Objs {
			if cv.Pos.Cross(x, y, w, h) {
				return false
			}
		}
	}

	return true
}

func (q *TileMap) Insert(o *Obj) bool {
	if o == nil {
		return false
	}

	if !q.CanInsert(o.Pos.X, o.Pos.Y, o.Pos.W, o.Pos.H) {
		return false
	}

	// 褪色
	for _, v := range o.Cells {
		delete(q.Cells[v-1].Objs, o.Id)
	}

	// 染色
	newCells := q.GetRectCell(o.Pos.X, o.Pos.Y, o.Pos.W, o.Pos.H)
	for _, v := range newCells {
		q.Cells[v-1].Objs[o.Id] = o
	}

	o.Cells = newCells

	return true
}

func (q *TileMap) Erase(o *Obj) {
	// 褪色
	if o != nil {
		for _, v := range o.Cells {
			delete(q.Cells[v-1].Objs, o.Id)
		}
	}
	o.Cells = []int32{}
}

func (q *TileMap) Clear() {
	if q.initOk {
		for i := uint32(0); i < q.CellCount; i++ {
			q.Cells[i].Objs = make(map[int16]*Obj, 0)
		}
	}
}

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
