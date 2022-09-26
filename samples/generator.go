package sample

import (
	"learngrpc/pcbook/pb"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// NewKeyboard returns a sample keyboard.
func NewKeyboard() *pb.Keyboard {
	keyboard := &pb.Keyboard{
		Layout: randomKeyboardLayout(),
		Backlit: randomBool(),
	}
	return keyboard
}

// NewCPU returns a sample CPU.
func NewCPU() *pb.CPU {
	brand := randomCPUBrand()
	name := randomCPUName(brand)
	numberCores := randomInt(2, 8)
	numberThreads := randomInt(numberCores, 12)
	minxGhz := randomFloat64(2.0, 3.5)
	maxGhz := randomFloat64(minxGhz, 5.0)

	cpu := &pb.CPU{
		Brand: brand,
		Name: name,
		NumCores: uint32(numberCores),
		NumThreads: uint32(numberThreads),
		MinGhz: minxGhz,
		MaxGhz: maxGhz,
	}
	return cpu
}

// NewGPU returns a sample GPU.
func NewGPU() *pb.GPU {
	brand := randomGPUBrand()
	name := randomGPUName(brand)
	minxGhz := randomFloat64(1.0, 1.5)
	maxGhz := randomFloat64(minxGhz, 2.0)
	memory := &pb.Memory {
		Value: uint64(randomInt(2, 6)),
		Unit: pb.Memory_GIGABYTE,
	}

	gpu := &pb.GPU{
		Brand: brand,
		Name: name,
		MinGhz: minxGhz,
		MaxGhz: maxGhz,
		Memory: memory,
	}
	return gpu
}

// NewRAM returns a sample RAM.
func NewRAM() *pb.Memory {
	memory := &pb.Memory {
		Value: uint64(randomInt(4, 64)),
		Unit: pb.Memory_GIGABYTE,
	}
	return memory
}

// NewSSD returns a sample SSD.
func NewSSD() *pb.Storage {
	storage := &pb.Storage {
		Driver: pb.Storage_SSD,
		Memory: &pb.Memory {
			Value: uint64(randomInt(128, 1024)),
			Unit: pb.Memory_GIGABYTE,
		},
	}
	return storage
}

// NewHDD returns a sample HDD.
func NewHDD() *pb.Storage {
	storage := &pb.Storage {
		Driver: pb.Storage_HDD,
		Memory: &pb.Memory {
			Value: uint64(randomInt(1, 6)),
			Unit: pb.Memory_TERABYTE,
		},
	}
	return storage
}

// NewScreen returns a sample screen.
func NewScreen() *pb.Screen {
	screen := &pb.Screen {
		Resolution: randomScreenResolution(),
		SizeInch: randomFloat32(13, 17),
		Panel: randomScreenPanel(),
		MultiTouch: randomBool(),
	}
	return screen
}

// NewLaptop returns a sample laptop.
func NewLaptop() *pb.Laptop{
	brand := randomLaptopBrand()
	name := randomLaptopName(brand)

	laptop := &pb.Laptop {
		Id: randomID(),
		Brand: brand,
		Name: name,
		Cpu: NewCPU(),
		Gpu: []*pb.GPU{NewGPU()},
		Memory: NewRAM(),
		Storage: []*pb.Storage{ NewSSD(), NewHDD()},
		Screen: NewScreen(),
		Keyboard: NewKeyboard(),
		Weight: &pb.Laptop_WeightKg {
			WeightKg: randomFloat32(1.0, 3.0),
		},
		PriceUsd: randomFloat64(1500, 3500),
		ReleaseYear: uint32(randomInt(2015, 2022)),
		ReleaseDate: timestamppb.Now(),
	}	

	return laptop
}