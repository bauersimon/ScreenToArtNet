package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/bauersimon/ScreenToArtNet/dmx"

	"github.com/bauersimon/ScreenToArtNet/ambilight"
	"github.com/bauersimon/ScreenToArtNet/capture"
)

func run() error {
	areas, universes, mapping, err := ambilight.ReadConfig(*args.Config)
	if err != nil {
		return err
	}

	s := capture.NewScreen(
		areas,
		capture.CaptureConfig{
			Spacing:   *args.Spacing,
			Threshold: *args.Threshold,
			Monitor:   *args.Screen,
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
			Sleep: *args.Pause,
		},
	}

	return a.Go()
}

func preview() error {
	areas, _, _, err := ambilight.ReadConfig(*args.Config)
	if err != nil {
		return err
	}

	s := capture.NewScreen(
		areas,
		capture.CaptureConfig{
			Monitor: *args.Screen,
		},
	)

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	return s.SavePreview(filepath.Join(cwd, "preview"))
}

var args = struct {
	Mode      *string
	Src       *string
	Dst       *string
	Pause     *int
	Screen    *int
	Spacing   *int
	Threshold *int
	Config    *string
}{
	flag.String("mode", "run", "tool mode {run|preview}"),
	flag.String("src", "", "artnet source"),
	flag.String("dst", "", "artnet destination"),
	flag.Int("pause", 0, "pause time in ms"),
	flag.Int("screen", 0, "screen identifier"),
	flag.Int("spacing", 1, "spacing of pixels for averaging"),
	flag.Int("threshold", 0, "threshold of color (0<255)"),
	flag.String("config", "config.json", "config file"),
}

func main() {
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		return
	}
	flag.Parse()

	// Make sure we clean everything up.
	abort := make(chan os.Signal)
	signal.Notify(abort, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-abort
		fmt.Printf("\r%v received, stopping...\n", s)
		os.Exit(0)
	}()

	switch *args.Mode {
	case "run":
		err := run()
		if err != nil {
			crash(err)
		}
	case "preview":
		err := preview()
		if err != nil {
			crash(err)
		}
	default:
		fmt.Printf("unknown mode: %s", *args.Mode)
	}
}

func crash(err error) {
	fmt.Printf("encountered error:\n%s", err.Error())
	os.Exit(-1)
}
