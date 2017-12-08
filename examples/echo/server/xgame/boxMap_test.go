package xgame

import "testing"

type MapObj struct {
	Objs map[int16]*Obj
}

func (m *MapObj) Init() {
	m.Objs = make(map[int16]*Obj, 0)
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

	mapA := &BoxMap{}
	mapA.InitMap(100, 100, 500, 500, 0)

	if !mapA.Insert(mapObjs.Objs[1], nil) {
		t.Error("insert 1 failed")
	}

	if !mapA.Insert(mapObjs.Objs[2], nil) {
		t.Error("insert 2 failed")
	}

	if !mapA.Insert(mapObjs.Objs[3], nil) {
		t.Error("insert 3 failed")
	}

}
