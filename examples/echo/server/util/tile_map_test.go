package util

import "testing"

type MapObj struct {
	Objs map[int16]*Obj
}

func (m *MapObj) Init() {
	m.Objs = make(map[int16]*Obj, 0)
}

func Test_Tilemap(t *testing.T) {

	mapObjs := &MapObj{}
	mapObjs.Init()

	mapObjs.Objs[1] = &Obj{Id: 1, Pos: Rect{X: 1, Y: 1, W: 20, H: 20}, Cells: []int32{}}
	mapObjs.Objs[2] = &Obj{Id: 2, Pos: Rect{X: 30, Y: 30, W: 20, H: 20}, Cells: []int32{}}

	mapA := &TileMap{}
	mapA.Init(100, 100, 300, 300)

	if !mapA.Insert(mapObjs.Objs[1]) {
		t.Error("insert 1 failed")
	}

	if !mapA.Insert(mapObjs.Objs[2]) {
		t.Error("insert 2 failed")
	}

}
