package util

import "testing"

type MapObj struct {
	Objs map[int16]*Obj
}

func (m *MapObj) Init() {
	m.Objs = make(map[int16]*Obj, 0)
}

type TRand struct {
	Seed int64
}

func (t *TRand) Random(maxValue int64) int64 {
	t.Seed = int64(int32(t.Seed*0x08088405 + 1))
	return ((t.Seed * maxValue) >> 32) + maxValue/2
}

func Test_Random(t *testing.T) {

	s := &TRand{}
	s.Seed = 99
	t.Error(s.Random(100))
	t.Error(s.Random(100))
	t.Error(s.Random(100))
	t.Error(s.Random(100))
	t.Error(s.Random(100))
	t.Error(s.Random(100))
	t.Error(s.Random(100))
}

func Test_Box(t *testing.T) {

	mapObjs := &MapObj{}
	mapObjs.Init()

	mapObjs.Objs[1] = &Obj{ID: 1, Pos: Rect{X: 0, Y: 0, W: 30, H: 20}, Cells: []int32{}}
	mapObjs.Objs[2] = &Obj{ID: 2, Pos: Rect{X: 260, Y: 40, W: 20, H: 20}, Cells: []int32{}}
	mapObjs.Objs[3] = &Obj{ID: 3, Pos: Rect{X: 300, Y: 40, W: 180, H: 20}, Cells: []int32{}}

	mapA := &Box{}
	mapA.Init(100, 100, 500, 500)

	if !mapA.Insert(mapObjs.Objs[1], nil) {
		t.Error("insert 1 failed")
	}

	if !mapA.Insert(mapObjs.Objs[2], nil) {
		t.Error("insert 2 failed")
	}

	if !mapA.Insert(mapObjs.Objs[3], nil) {
		t.Error("insert 3 failed")
	}

	// mapA.ObjMove(mapObjs.Objs[1], DirDown, 10)
	// mapA.ObjMove(mapObjs.Objs[1], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[1], DirDown, 10)
	// mapA.ObjMove(mapObjs.Objs[1], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[1], DirDown, 10)
	// mapA.ObjMove(mapObjs.Objs[1], DirDown, 10)
	// mapA.ObjMove(mapObjs.Objs[1], DirDown, 10)

	// t.Error("rect ", mapObjs.Objs[1].Pos)

	// mapA.ObjMove(mapObjs.Objs[3], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[3], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[3], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[3], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[3], DirRight, 10)

	// t.Error("rect3 ", mapObjs.Objs[3].Pos)

	// mapA.ObjMove(mapObjs.Objs[2], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[2], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[2], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[2], DirRight, 10)
	// mapA.ObjMove(mapObjs.Objs[2], DirRight, 10)
	// t.Error("rect2 ", mapObjs.Objs[2].Pos)

}
