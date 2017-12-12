package xgame

// TRamd 自定义随机数
type TRand struct {
	Seed int64
}

func (t *TRand) Random64(maxValue int64) int64 {
	t.Seed = int64(int32(t.Seed*0x08088405 + 1))
	return ((t.Seed * maxValue) >> 32) + maxValue/2
}

func (t *TRand) Random32(maxValue int32) int32 {
	t.Seed = int64(int32(t.Seed*0x08088405 + 1))
	return int32(((t.Seed * int64(maxValue)) >> 32) + int64(maxValue)/2)
}
