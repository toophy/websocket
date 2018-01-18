package mud

// Item 道具
type Item struct {
	ID    int32           // 唯一ID
	Type  int32           // 类型ID
	Count int32           // 数量, 为0表示不能重叠
	Data  map[int32]int32 // 自定义数据
}
