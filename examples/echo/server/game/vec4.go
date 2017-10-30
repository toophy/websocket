package game

// 位置
type Vec4 struct {
	X      int32   // x
	Y      int32   // y, 角色在陆地上, 没有这个属性变化, 飘在场景中有这个属性
	Height int32   // 高度, 在陆地上的高度
	Parent *MapObj // 父对象
}
