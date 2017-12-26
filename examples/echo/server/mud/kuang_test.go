package mud

import (
	"testing"
)

func Test_kuang(t *testing.T) {
	ks := KuangSys{}
	ks.InitSys()

	for {
		ks.Update()
		t.Error("hhh")
	}
}
