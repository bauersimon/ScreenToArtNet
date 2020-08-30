package capture

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

	"github.com/kbinani/screenshot"
)

// Screen represents a tiled screen.
type Screen struct {
	Areas   []*image.Rectangle
	Borders image.Rectangle
}

// NewScreen returns a new screen, tiled with the given configuration.
func NewScreen(areas []*image.Rectangle, monitor int) *Screen {
	return &Screen{
		Areas:   areas,
		Borders: screenshot.GetDisplayBounds(monitor),
	}
}

// GetColors returns an averaged color per screen tile.
func (s Screen) GetColors(space int, threshold int) ([]color.RGBA, error) {
	var colors []color.RGBA

	monitor, err := screenshot.CaptureRect(s.Borders)
	if err != nil {
		return nil, err
	}

	for _, b := range s.Areas {
		area := monitor.SubImage(*b).(*image.RGBA)

		c, err := averageRGBA(area, space, threshold)
		if err != nil {
			return nil, err
		}

		colors = append(colors, c)
	}

	return colors, nil
}

func averageRGBA(area *image.RGBA, space int, threshold int) (color.RGBA, error) {
	var r uint64
	var g uint64
	var b uint64

	var count uint64 = 1

	if space < 1 {
		return color.RGBA{}, fmt.Errorf("invalid spacing for averaging (%v)", space)
	}
	if threshold < 0 || threshold > 255 {
		return color.RGBA{}, fmt.Errorf("invalid threshold for averaging (%v)", threshold)
	}

	for x := area.Rect.Min.X; x < area.Rect.Max.X; x = x + space {
		for y := area.Rect.Min.Y; y < area.Rect.Max.Y; y = y + space {
			pixel := color.RGBAModel.Convert(area.At(x, y)).(color.RGBA)
			lr := pixel.R
			lg := pixel.G
			lb := pixel.B

			average := (uint32(lr) + uint32(lg) + uint32(lb)) / 3
			if average < uint32(threshold) {
				continue
			}

			r = r + uint64(lr)
			g = g + uint64(lg)
			b = b + uint64(lb)

			count++
		}
	}

	return color.RGBA{
		R: uint8(r / count),
		G: uint8(g / count),
		B: uint8(b / count),
		A: 255,
	}, nil
}

func saveImage(dst string, img *image.RGBA) {
	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create(dst)
	if err != nil {
		panic(err)
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	png.Encode(outputFile, img)

	// Don't forget to close files
	outputFile.Close()
}
