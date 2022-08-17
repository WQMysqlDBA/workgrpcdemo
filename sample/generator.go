package sample

import (
	"workgrpc/pb"
)

func NewKeyboard() *pb.Keyboard {
	return &pb.Keyboard{
		Layout:  randomKeyboardLayout(),
		Backlit: randomBool(),
	}
}
func NewCPU() *pb.CPU {
	brand := randCpuBrand()
	name:= randCpuName(brand)
	cores := randomInt(2,8)
	threads := randomInt(cores,12)
	minGhz :=randomFloat64(2.0,3.5)
	maxGhz :=randomFloat64(minGhz,5.0)
	cpu := &pb.CPU{
		Brand:         brand,
		Name:          name,
		NumberCores:   uint32(cores),
		NumberThreads: uint32(threads),
		MinGhz:        minGhz,
		MaxGhz:        maxGhz,
	}
	return cpu
}
