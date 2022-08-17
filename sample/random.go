package sample

import (
	"math/rand"
	"time"
	"workgrpc/pb"
)

const (
	Intel = "Intel"
	AMD   = "AMD"
	Apple = "Apple"
)

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 1:
		return pb.Keyboard_QWERTY
	case 2:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}
}

func randomBool() bool {
	return rand.Intn(2) == 1

}

func randCpuBrand() string {
	return randomStringFormSet(Intel, AMD, Apple)
}
func randomStringFormSet(a ...string) string {
	n := len(a)
	rand.Seed(time.Now().UnixNano())
	if n == 0 {
		return ""
	} else {
		return a[rand.Intn(n)]
	}
}
func randCpuName(brand string) string {
	if brand == Intel {
		return randomStringFormSet("i9-9880XE", "i9-9960X", "i9-9940X",
			"i9-9920X", "i9-9900X", "i9-9820X", "i7-9800X", "W-3175X",
			"i9-9900K/KF",
			"i9-9900/T",
			"i7-9700K/KF",
			"i7-9700/F/T",
		)
	} else if brand == AMD {
		return randomStringFormSet(
			"Athlon 64 X2 5200+", "AMD Phenom II X2 545AMD", "Phenom II X2 550AMD", "Phenom II X3 720AMD", "Phenom X4 9350eAMD")
	} else {
		return randomStringFormSet("Apple M1", "Apple M1 pro", "Apple M1 max", "Apple m2")
	}
}

func randomInt(min, max int) int {
	return min + rand.Intn(max-min+1)
}
func randomFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}
