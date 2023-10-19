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
// TODO: handle negative values
type AmbilightConfiguration struct {
	// target frames per second
	Fps int
	// max number of worker threads
	Workers int
}

// Go fires up the ambilight.
func (a *Ambilight) Go() error {
	// Data for the dynamic performance display.
	iter := 0
	maxIter := 10
	start := time.Now()
	frameDurationTarget := time.Duration(1000/a.Config.Fps) * time.Millisecond

	// Create worker pool that can perform as many tasks as possible, given the allowed workers
	maxParallelTasks := max(len(a.Screen.Areas), len(a.Universes))
	workerPool := newWorkerPool(min(a.Config.Workers, maxParallelTasks))

	for {
		frameStart := time.Now()
		err := a.Screen.Capture()
		if err != nil {
			panic(err)
		}

		var colorJobs queue
		for _, area := range a.Screen.Areas {
			colorJobs.enqueue(getColorJob(a, area))
		}
		workerPool.workOn(colorJobs)

		var networkJobs queue
		for _, u := range a.Universes {
			networkJobs.enqueue(func() {

				err := u.SendColorUpdate(a.Controller)
				if err != nil {
					panic(err)
				}
			})
		}
		workerPool.workOn(networkJobs)

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

		frameDurationCurrent := time.Since(frameStart)
		timeToSleep := frameDurationTarget - frameDurationCurrent
		time.Sleep(timeToSleep)
	}
}

func getColorJob(a *Ambilight, area capture.Area) func() {
	// TODO: maybe link directly in Area?
	devices, ok := a.Mappings[area.ImageData.Borders]
	if !ok {
		// This area has no devices mapped.
		return func() {}
	}

	return func() {
		area.ImageData.Update()
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

}

func min(a int, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a int, b int) int {
	if a > b {
		return a
	}
	return b
}

// Mapping holds a mapping from screen areas to DMX devices.
// FIXME: it makes no sense to have more than one device per area as this would overwrite colors
type Mapping map[*image.Rectangle][]*dmx.Device
