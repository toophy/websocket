// 四叉树扩展
package util

// 1. 静态的一维或者二维表
// 2. 用运算得到关联下标, 上级, 下级
// 3. 范围的作用, 不得超出规定范围
// 4. 非完全四叉树, 横向比纵向多很多(横版游戏), 并不是常规的正方形, 而是长方形, 但是每个item都是正方形, 部分符合四叉树
// 5. 一开始就搭建好, 没有删除item功能
// 6. 矩形和点可以在里面游荡, 就是不断的离开一个item, 在进入另外一个item
// 7. 查询, 每次都是从头开始查询, 从最上层到最下层
// 8. 插入一个矩形也是, 从最上层=>最下层
type Rect struct {
	X int32
	Y int32
	W int32
	H int32
}

type Obj struct {
	Id  int16
	Pos Rect
}

type Cell struct {
	Objs map[int16]*Obj
}

// 四叉树+9宫格
type Q9Tree struct {
	initOk    bool
	W         int32
	H         int32
	tileW     int32
	tileH     int32
	realW     int32
	realH     int32
	level     int32
	CellCount uint32
	Cells     []Cell
}

type Map struct {
	Objs map[int16]*Obj
}

func (q *Q9Tree) Init(tileW, tileH, realW, realH int32) bool {

	if tileW < 1 || tileH < 1 || realW < 0 || realH < 0 || tileW > realW || tileH > realH {
		return false
	}

	q.Destroy()
	q.realW = realW
	q.realH = realH
	q.tileW = tileW
	q.tileH = tileH
	w := (realW + tileW - 1) / tileW
	h := (realH + tileH - 1) / tileH

	wLayer := int32(0)
	for i := 0; ; i++ {
		if w/2 > 0 {
			wLayer++
		} else {
			break
		}
	}

	hLayer := int32(0)
	for i := 0; ; i++ {
		if h/2 > 0 {
			hLayer++
		} else {
			break
		}
	}
	q.level = wLayer
	if hLayer > q.level {
		q.level = hLayer
	}

	q.CellCount = uint32(0)
	for i, j := uint32(wLayer-1), uint32(hLayer-1); i > 1 && j > 1; {

		q.CellCount += (uint32(1) << i) * (uint32(1) << j)

		if i > 1 {
			i--
		}
		if j > 1 {
			j--
		}
	}

	q.Cells = make([]Cell, q.CellCount)

	q.initOk = true

	q.Clear()

	return true
}

func (q *Q9Tree) GetCell(x, y, z int32) int32 {
	return -1
}

func (q *Q9Tree) Insert(o *Obj) bool {
	return true
}

func (q *Q9Tree) Erase(o *Obj) bool {
	return true
}

func (q *Q9Tree) Retrieve(r Rect) []*Obj {
	return nil
}

func (q *Q9Tree) RetrievePoints(r Rect) []*Obj {
	return nil
}

func (q *Q9Tree) RetrieveIntersections(r Rect) []*Obj {
	return nil
}

func (q *Q9Tree) Clear() {
	if q.initOk {
		for i := uint32(0); i < q.CellCount; i++ {
			q.Cells[i].Objs = make(map[int16]*Obj, 0)
		}
	}
}

func (q *Q9Tree) Destroy() {
	if q.initOk {
		q.W = 0
		q.H = 0
		q.tileW = 0
		q.tileH = 0
		q.realW = 0
		q.realH = 0
		q.level = 0
		q.Cells = []Cell{}
		q.CellCount = 0
	}
}
