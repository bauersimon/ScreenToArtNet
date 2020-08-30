package capture

import (
	"image"
	"testing"
)

func BenchmarkCapture(b *testing.B) {
	s := NewScreen(
		[]*image.Rectangle{
			&image.Rectangle{
				Min: image.Point{0, 0},
				Max: image.Point{800, 600},
			},
		},
		0,
	)

	spacings := map[string]int{
		"dense":   1,
		"space 2": 2,
		"space 4": 4,
	}

	for name, space := range spacings {
		b.Run(name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				s.GetColors(space, 0)
			}
		})
	}
}
