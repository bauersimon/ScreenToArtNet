package capture

import (
	"image"
	"testing"
)

func BenchmarkCapture(b *testing.B) {

	spacings := map[string]int{
		"dense":   1,
		"space 2": 2,
		"space 4": 4,
	}

	for name, space := range spacings {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s := getScreen(space)
				s.Capture()
				s.ImageData.GetImage()
				for _, area := range s.Areas {
					area.ImageData.GetImage()
					area.ImageData.GetColor()
				}
			}
		})
	}
}

func getScreen(space int) *Screen {

	config := CaptureConfig{
		imageDataConfig: imageDataConfig{
			Spacing:   space,
			Threshold: 0,
		},
		Monitor: 1,
	}

	area := Area{
		Name: "bla",
		ImageData: ImageData{
			Borders: image.Rectangle{
				Min: image.Point{0, 0},
				Max: image.Point{800, 600},
			}},
	}
	s := NewScreen(
		[]Area{area},
		config,
	)

	return s
}
