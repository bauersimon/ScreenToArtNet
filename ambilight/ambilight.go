package ambilight

import (
	"fmt"
	"image"
	"time"

	"github.com/bauersimon/ScreenToArtNet/capture"
	"github.com/bauersimon/ScreenToArtNet/dmx"
)

// Ambilight holds all the information of an ambilight.
type Ambilight struct {
	// Controller holds the ArtNet controller.
	Controller *dmx.ArtNetController
	// Screen holds the screen configuration.
	Screen *capture.Screen
	// Universes hods the DMX universes.
	Universes []*dmx.Universe
	// Mappings holds the screen area to DMX devices mapping.
	Mappings Mapping

	Config AmbilightConfiguration
}

// AmbilightConfiguration holds the configuration for the ambilight.
type AmbilightConfiguration struct {
	// Sleep holds the sleep time after each update in ms.
	Sleep int
}

// Go fires up the ambilight.
func (a *Ambilight) Go() error {
	// Data for the dynamic performance display.
	iter := 0
	maxIter := 10
	start := time.Now()

	for {
		err := a.Screen.Capture()
		if err != nil {
			panic(err)
		}

		for _, area := range a.Screen.Areas {
			devices, ok := a.Mappings[&area.ImageData.Borders]
			if !ok {
				// This area has no devices mapped.
				continue
			}

			areaColor, err := area.ImageData.GetColor()
			if err != nil {
				panic(err)
			}
			for _, d := range devices {
				d.RValue = areaColor.R
				d.GValue = areaColor.G
				d.BValue = areaColor.B
			}
		}

		for _, u := range a.Universes {
			err := u.SendColorUpdate(a.Controller)
			if err != nil {
				return err
			}
		}

		if a.Config.Sleep > 0 {
			time.Sleep(time.Duration(a.Config.Sleep) * time.Millisecond)
		}

		// Handle the performance display
		if iter == maxIter {
			diff := time.Since(start).Seconds()
			updates := float64(iter) / diff
			fmt.Printf("%.2f updates/sec\r", updates)
			maxIter = int(updates) * 5 // We will perform roughly one print every 5 seconds.

			iter = 0
			start = time.Now()
		} else {
			iter++
		}
	}
}

// Mapping holds a mapping from screen areas to DMX devices.
type Mapping map[*image.Rectangle][]*dmx.Device
