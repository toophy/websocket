package mud

import (
	"testing"
)

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
}
