package sample

import (
	"learngrpc/pcbook/pb"
	"math/rand"

	"github.com/google/uuid"
)

func randomKeyboardLayout() pb.Keyboard_Layout {
	switch rand.Intn(3) {
	case 0:
		return pb.Keyboard_QWERTY
	case 1:
		return pb.Keyboard_QWERTZ
	default:
		return pb.Keyboard_AZERTY
	}

}

func randomCPUBrand() string {
	return randomStringFromSet("Intel", "AMD", "ARM")
}

func randomCPUName(brand string) string {
	switch brand {
	case "Intel":
		return randomStringFromSet("i3", "i5", "i7", "i9")
	case "AMD":
		return randomStringFromSet("Ryzen 3", "Ryzen 5", "Ryzen 7", "Ryzen 9")
	case "ARM":
		return randomStringFromSet("Cortex A53", "Cortex A55", "Cortex A57", "Cortex A72", "Cortex A73", "Cortex A75")
	default:
		return ""
	}
}

func randomGPUBrand() string {
	return randomStringFromSet("NVIDIA", "AMD", "Intel")
}	

func randomGPUName(brand string) string {
	switch brand {
	case "NVIDIA":
		return randomStringFromSet("GTX 1050", "GTX 1060", "GTX 1070", "GTX 1080", "GTX 1650", "GTX 1660", "GTX 1660 Ti", "GTX 2060", "GTX 2070", "GTX 2080", "GTX 2080 Ti")
	case "AMD":
		return randomStringFromSet("Radeon 530", "Radeon 550", "Radeon 560", "Radeon 570", "Radeon 580", "Radeon 590", "Radeon 630", "Radeon 640", "Radeon 650", "Radeon 660", "Radeon 670", "Radeon 680", "Radeon 690", "Radeon 700", "Radeon 710", "Radeon 720", "Radeon 730", "Radeon 740", "Radeon 750", "Radeon 760", "Radeon 770", "Radeon 780", "Radeon 790", "Radeon 800", "Radeon 810", "Radeon 820", "Radeon 830", "Radeon 840", "Radeon 850", "Radeon 860", "Radeon 870", "Radeon 880", "Radeon 890", "Radeon 900", "Radeon 910", "Radeon 920", "Radeon 930", "Radeon 940", "Radeon 950", "Radeon 960", "Radeon 970", "Radeon 980", "Radeon 990", "Radeon 1000", "Radeon 1010", "Radeon 1020", "Radeon 1030", "Radeon 1040", "Radeon 1050", "Radeon 1060", "Radeon 1070", "Radeon 1080", "Radeon 1090", "Radeon 1100", "Radeon 1110", "Radeon 1120", "Radeon 1130", "Radeon 1140", "Radeon 1150", "Radeon 1160", "Radeon 1170", "Radeon 1180", "Radeon 1190", "Radeon 1200", "Radeon 1210", "Radeon 1220", "Radeon 1230", "Radeon 1240", "Radeon 1250", "Radeon 1260", "Radeon 1270", "Radeon 1280", "Radeon 1290", "Radeon 1300", "Radeon 1310", "Radeon 1320")
	case "Intel":
		return randomStringFromSet("UHD Graphics 620", "UHD Graphics 630", "UHD Graphics 640", "UHD Graphics 650", "UHD Graphics 660", "UHD Graphics 670", "UHD Graphics 680", "UHD Graphics 690", "UHD Graphics 700", "UHD Graphics 710", "UHD Graphics 720", "UHD Graphics 730", "UHD Graphics 740", "UHD Graphics 750", "UHD Graphics 760", "UHD Graphics 770", "UHD Graphics 780", "UHD Graphics 790", "UHD Graphics 800", "UHD Graphics 810", "UHD Graphics 820", "UHD Graphics 830", "UHD Graphics 840", "UHD Graphics 850", "UHD Graphics 860", "UHD Graphics 870", "UHD Graphics 880", "UHD Graphics 890", "UHD Graphics 900", "UHD Graphics 910", "UHD Graphics 920", "UHD Graphics 930", "UHD Graphics 940", "UHD Graphics 950", "UHD Graphics 960", "UHD Graphics 970", "UHD Graphics 980", "UHD Graphics 990", "UHD Graphics 1000", "UHD Graphics 1010", "UHD Graphics 1020", "UHD Graphics 1030", "UHD Graphics 1040", "UHD Graphics 1050", "UHD Graphics 1060", "UHD Graphics 1070", "UHD Graphics 1080", "UHD Graphics 1090", "UHD Graphics 1100", "UHD Graphics 1110", "UHD Graphics 1120", "UHD Graphics 1130", "UHD Graphics 1140", "UHD Graphics 1150", "UHD Graphics 1160", "UHD Graphics 1170", "UHD Graphics 1180", "UHD Graphics 1190", "UHD Graphics 1200", "UHD Graphics 1210", "UHD Graphics 1220", "UHD Graphics 1230", "UHD Graphics 1240", "UHD Graphics 1250", "UHD Graphics 1260", "UHD Graphics 1270")
	default:
		return ""
	}
}

func randomScreenResolution() *pb.Screen_Resolution {
	height := randomInt(13, 17)
	width := height * 16/9

	return &pb.Screen_Resolution{
		Width: uint32(width),
		Height: uint32(height),
	}
}

func randomScreenPanel() pb.Screen_Panel {
	value := rand.Intn(2)
	switch value {
	case 0:
		return pb.Screen_IPS
	default:
		return pb.Screen_OLED
	}
}

func randomLaptopBrand() string {
	return randomStringFromSet("Apple", "Dell", "Lenovo", "Microsoft", "Asus")
}

func randomLaptopName(brand string) string {
	switch brand {
	case "Apple":
		return randomStringFromSet("Macbook Air", "Macbook Pro", "Macbook Pro 16")
	case "Dell":
		return randomStringFromSet("XPS 13", "XPS 15", "XPS 17")
	case "Lenovo":
		return randomStringFromSet("Thinkpad X1 Carbon", "Thinkpad X1 Yoga", "Thinkpad X1 Extreme")
	case "Microsoft":
		return randomStringFromSet("Surface Book 2", "Surface Laptop 3", "Surface Pro 7")
	case "Asus":
		return randomStringFromSet("Zenbook", "Vivobook", "TUF Gaming")
	default:
		return ""
	}
}

func randomStringFromSet(a ...string) string {
	if len(a) == 0 {
		return ""
	}

	return a[rand.Intn(len(a))]
}

func randomBool() bool {
	return rand.Intn(2) == 1
}

func randomInt(min, max int) int {
	// retunr random int of [min, max]
	return rand.Intn(max-min + 1) + min
}

func randomFloat64(min, max float64) float64 {
	return min + rand.Float64() * (max - min)
}

func randomFloat32(min, max float32) float32 {
	return min + rand.Float32() * (max - min)
}

func randomID() string {
	return uuid.New().String()
}