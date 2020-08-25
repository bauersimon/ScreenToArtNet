package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"
	"time"

	"github.com/bauersimon/ambilight/capture"
	"github.com/bauersimon/ambilight/dmx"
)

// Ambilight holds all the information of an ambilight.
type Ambilight struct {
	Controller *dmx.Controller
	Screen     capture.Screen
	Devices    []dmx.RGBaddress
}

// NewAmbilight creates a new ambilight with the given configuration.
func NewAmbilight(areas []image.Rectangle, devices []dmx.RGBaddress, src string, dst string, screen int) (Ambilight, error) {
	c, err := dmx.NewArtNetController(src, dst)
	if err != nil {
		return Ambilight{}, err
	}

	if len(areas) != len(devices) {
		return Ambilight{}, fmt.Errorf("areas and dmx devices don't match (%v vs. %v)", len(areas), len(devices))
	}

	s := capture.NewScreen(areas, screen)
	fmt.Printf("screen dimensions: %v\n", s.Borders)

	return Ambilight{
		Controller: c,
		Screen:     s,
		Devices:    devices,
	}, nil
}

func unmarshalConfig(data json.RawMessage) (areas []image.Rectangle, devices []dmx.RGBaddress, err error) {
	var definitions []map[string]json.RawMessage

	err = json.Unmarshal(data, &definitions)
	if err != nil {
		return nil, nil, err
	}

	for _, entry := range definitions {
		areaRaw, ok := entry["area"]
		if !ok {
			return nil, nil, fmt.Errorf("missing area entry for: %v", entry)
		}

		var area image.Rectangle
		err := json.Unmarshal(areaRaw, &area)
		if err != nil {
			return nil, nil, err
		}

		deviceRaw, ok := entry["device"]
		if !ok {
			return nil, nil, fmt.Errorf("missing device entry for: %v", entry)
		}

		var device dmx.RGBaddress
		err = json.Unmarshal(deviceRaw, &device)
		if err != nil {
			return nil, nil, err
		}

		err = device.Verify()
		if err != nil {
			return nil, nil, err
		}

		areas = append(areas, area)
		devices = append(devices, device)
	}

	return areas, devices, nil
}

func readConfig(configPath string) (json.RawMessage, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return ioutil.ReadFile(path.Join(cwd, configPath))
}

func main() {
	src := flag.String("src", "", "artnet source")
	dst := flag.String("dst", "", "artnet destination")
	pause := flag.Int("pause", 0, "pause time in ms")
	screen := flag.Int("screen", 0, "screen identifier")
	spacing := flag.Int("spacing", 1, "spacing of pixels for averaging")
	threshold := flag.Int("threshold", 0, "threshold of color (0<255)")
	config := flag.String("config", "config.json", "config file")
	flag.Parse()

	configJson, err := readConfig(*config)
	if err != nil {
		fmt.Println(err)
		return
	}

	areas, devices, err := unmarshalConfig(configJson)
	if err != nil {
		fmt.Println(err)
		return
	}

	ambi, err := NewAmbilight(areas, devices, *src, *dst, *screen)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Make sure we clean everything up.
	abort := make(chan os.Signal)
	signal.Notify(abort, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-abort
		fmt.Printf("\r%v received, stopping...\n", s)
		os.Exit(0)
	}()

	for {
		colors, err := ambi.Screen.GetColors(*spacing, *threshold)
		if err != nil {
			panic(err)
		}

		for i, d := range ambi.Devices {
			err := dmx.SendColorUpdate(ambi.Controller, d, colors[i])
			if err != nil {
				panic(err)
			}
		}

		if *pause > 0 {
			time.Sleep(time.Duration(*pause) * time.Millisecond)
		}
	}
}
