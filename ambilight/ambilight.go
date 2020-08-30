package ambilight

import (
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
	for {
		colors, err := a.Screen.GetColors()
		if err != nil {
			panic(err)
		}

		for i, c := range colors {
			devices, ok := a.Mappings[a.Screen.Areas[i]]
			if !ok {
				// This area has no devices mapped.
				continue
			}

			for _, d := range devices {
				d.RValue = c.R
				d.GValue = c.G
				d.BValue = c.B
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
	}
}

// Mapping holds a mapping from screen areas to DMX devices.
type Mapping map[*image.Rectangle][]*dmx.Device
