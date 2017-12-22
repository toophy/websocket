package mud

import (
	"math"
	"testing"
)

func CalcScore(step int32, winRate, kda float64, elo int32) int64 {
	return (int64(step) * 10000000000) + (int64(elo) * 1000000) + (int64(math.Ceil(kda*10)) * 1000) + (int64(math.Ceil(winRate * 10)))
}

func Test_match_elo(t *testing.T) {
	a := 1500.0
	b := 1600.0

	newA := CalcResult(a, b, false, false)
	newB := CalcResult(b, a, true, false)

	t.Error(newA, newB)

	a = newA
	b = newB

	newA = CalcResult(a, b, true, false)
	newB = CalcResult(b, a, false, false)

	t.Error(newA, newB)

	step := 40
	elo := 9999
	kda := 10.21
	winR := 30.3

	t.Error(CalcScore(int32(step), winR, kda, int32(elo)))
}
