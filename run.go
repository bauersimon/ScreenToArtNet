package main

import (
	"github.com/bauersimon/ScreenToArtNet/ambilight"
	"github.com/bauersimon/ScreenToArtNet/capture"
	"github.com/bauersimon/ScreenToArtNet/dmx"
)

func run() error {
	areas, universes, mapping, err := ambilight.ReadConfig(*args.Config)
	if err != nil {
		return err
	}

	s := capture.NewScreen(
		areas,
		capture.CaptureConfig{
			Monitor: *args.Screen,
			ImageDataConfig: capture.ImageDataConfig{
				Spacing:   *args.Spacing,
				Threshold: *args.Threshold,
			},
		},
	)

	c, err := dmx.NewArtNetController(
		*args.Src,
		*args.Dst,
	)
	if err != nil {
		return err
	}

	a := &ambilight.Ambilight{
		Controller: c,
		Screen:     s,
		Universes:  universes,
		Mappings:   mapping,
		Config: ambilight.AmbilightConfiguration{
			Fps:     *args.Fps,
			Workers: *args.Workers,
		},
	}

	return a.Go()
}
