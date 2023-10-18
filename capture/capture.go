package capture

import (
	"fmt"
	"image"
	"image/color"

	"github.com/kbinani/screenshot"
)

var NO_IMAGE image.RGBA

// Screen represents a tiled screen.
type Screen struct {
	// Areas holds the screen areas.
	Areas []Area
	// Last capture
	captureImage *image.RGBA
	// Whole monitor image data
	ImageData ImageData
}

type Area struct {
	Name      string
	ImageData ImageData
}

type ImageData struct {
	config  imageDataConfig
	image   **image.RGBA
	color   *color.RGBA
	Borders image.Rectangle
}

// NewScreen returns a new screen, tiled with the given configuration.
func NewScreen(areas []Area, config CaptureConfig) *Screen {

	screen := Screen{
		Areas:        areas,
		captureImage: &NO_IMAGE,
	}

	// init monitor image data with pointer to capture
	screen.ImageData = ImageData{
		config:  config.imageDataConfig,
		image:   &screen.captureImage,
		Borders: screenshot.GetDisplayBounds(config.Monitor),
	}

	// init area image data with pointer to capture
	for i := range areas {
		areas[i].ImageData.config = config.imageDataConfig
		areas[i].ImageData.image = &screen.captureImage
	}

	return &screen
}

func (s *Screen) Capture() error {
	monitorImage, err := screenshot.CaptureRect(s.ImageData.Borders)
	if err != nil {
		return err
	}

	s.captureImage = monitorImage

	return nil
}

func (d *ImageData) GetImage() (*image.RGBA, error) {

	monitorImage := *d.image
	if monitorImage == &NO_IMAGE {
		return nil, fmt.Errorf("no image data, capture screen first")
	}
	return monitorImage.SubImage(d.Borders).(*image.RGBA), nil

}

func (d *ImageData) GetColor() (color.RGBA, error) {
	var computedColor color.RGBA
	image, err := d.GetImage()
	if err != nil {
		return computedColor, err
	}
	if d.color != nil {
		return *d.color, nil
	}
	var r uint64
	var g uint64
	var b uint64

	var count uint64 = 1

	space := d.config.Spacing
	threshold := d.config.Threshold
	if space < 1 {
		return computedColor, fmt.Errorf("invalid spacing for averaging (%v)", space)
	}
	if threshold < 0 || threshold > 255 {
		return computedColor, fmt.Errorf("invalid threshold for averaging (%v)", threshold)
	}

	for x := image.Rect.Min.X; x < image.Rect.Max.X; x = x + space {
		for y := image.Rect.Min.Y; y < image.Rect.Max.Y; y = y + space {
			pixel := color.RGBAModel.Convert(image.At(x, y)).(color.RGBA)
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

	computedColor = color.RGBA{
		R: uint8(r / count),
		G: uint8(g / count),
		B: uint8(b / count),
		A: 255,
	}

	return computedColor, nil
}
