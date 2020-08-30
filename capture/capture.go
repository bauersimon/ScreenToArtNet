package capture

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/kbinani/screenshot"
)

// Screen represents a tiled screen.
type Screen struct {
	// Areas holds the screen areas.
	Areas []*image.Rectangle
	// Borders holds the capturing borders.
	Borders image.Rectangle

	// Config holds the configuration for the screen capture.
	Config CaptureConfig
}

// CaptureConfig holds the configuration for the screen capture.
type CaptureConfig struct {
	// Spacing holds the averaging spacing.
	Spacing int
	// Threshold holds the averaged color threshold.
	Threshold int

	// Monitor holds the monitor used for capture.
	Monitor int
}

// NewScreen returns a new screen, tiled with the given configuration.
func NewScreen(areas []*image.Rectangle, config CaptureConfig) *Screen {
	return &Screen{
		Areas:   areas,
		Borders: screenshot.GetDisplayBounds(config.Monitor),
		Config:  config,
	}
}

func (s *Screen) capture() (areas []*image.RGBA, monitor *image.RGBA, err error) {
	monitor, err = screenshot.CaptureRect(s.Borders)
	if err != nil {
		return nil, nil, err
	}

	areas = make([]*image.RGBA, len(s.Areas))
	for i, b := range s.Areas {
		areas[i] = monitor.SubImage(*b).(*image.RGBA)
	}

	return areas, monitor, nil
}

// GetColors returns an averaged color per screen tile.
func (s *Screen) GetColors() ([]color.RGBA, error) {
	var colors []color.RGBA

	areas, _, err := s.capture()
	if err != nil {
		return nil, err
	}

	for _, a := range areas {
		c, err := averageRGBA(a, s.Config.Spacing, s.Config.Threshold)
		if err != nil {
			return nil, err
		}

		colors = append(colors, c)
	}

	return colors, nil
}

// SavePreview saves the current capture configurations as multiple ".png" images at the given path.
func (s *Screen) SavePreview(dst string) error {
	areas, monitor, err := s.capture()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}

	err = saveArea(filepath.Join(dst, "monitor.png"), monitor)
	if err != nil {
		return err
	}

	for i, a := range areas {
		err = saveArea(filepath.Join(dst, fmt.Sprintf("area%d.png", i)), a)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveArea(dst string, area *image.RGBA) error {
	outputFile, err := os.Create(dst)
	if err != nil {
		return err
	}

	err = png.Encode(outputFile, area)
	if err != nil {
		return err
	}

	return outputFile.Close()
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
