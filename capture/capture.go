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
	config      imageDataConfig
	parentImage **image.RGBA
	subImage    *image.RGBA
	color       *color.RGBA
	Borders     *image.Rectangle
}

// NewScreen returns a new screen, tiled with the given configuration.
func NewScreen(areas []Area, config CaptureConfig) *Screen {

	screen := Screen{
		Areas:        areas,
		captureImage: &NO_IMAGE,
	}

	// init monitor image data with pointer to capture
	screen.ImageData = ImageData{
		config:      config.imageDataConfig,
		parentImage: &screen.captureImage,
	}

	// init area image data with pointer to capture
	borderUnion := image.Rect(0, 0, 0, 0)
	for i := range areas {
		borderUnion = borderUnion.Union(*areas[i].ImageData.Borders)
		areas[i].ImageData.config = config.imageDataConfig
		areas[i].ImageData.parentImage = &screen.captureImage
	}
	screen.ImageData.Borders = &borderUnion

	return &screen
}

func (s *Screen) Capture() error {
	monitorImage, err := screenshot.CaptureRect(*s.ImageData.Borders)
	if err != nil {
		return err
	}

	s.captureImage = monitorImage

	return nil
}

func (d *ImageData) Update() error {
	err := d.updateImage()
	if err != nil {
		return err
	}
	return d.updateColor()
}

func (d *ImageData) GetImage() (*image.RGBA, error) {

	if d.subImage == nil {
		return nil, fmt.Errorf("no image data, update first")
	}
	return d.subImage, nil

}

func (d *ImageData) GetColor() (*color.RGBA, error) {

	if d.color == nil {
		return nil, fmt.Errorf("no color data, update first")
	}
	return d.color, nil

}

func (d *ImageData) updateImage() error {
	monitorImage := *d.parentImage
	if monitorImage == &NO_IMAGE {
		return fmt.Errorf("no image data, capture screen first")
	}
	d.subImage = monitorImage.SubImage(*d.Borders).(*image.RGBA)
	return nil
}

func (d *ImageData) updateColor() error {
	var computedColor color.RGBA
	image, err := d.GetImage()
	if err != nil {
		return err
	}

	var r uint64
	var g uint64
	var b uint64

	var count uint64 = 1

	space := d.config.Spacing
	threshold := d.config.Threshold
	if space < 1 {
		return fmt.Errorf("invalid spacing for averaging (%v)", space)
	}
	if threshold < 0 || threshold > 255 {
		return fmt.Errorf("invalid threshold for averaging (%v)", threshold)
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

	d.color = &computedColor
	return nil
}
