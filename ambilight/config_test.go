package ambilight

import (
	"image"
	"reflect"
	"testing"

	"github.com/bauersimon/ScreenToArtNet/dmx"
	"github.com/stretchr/testify/assert"
)

func TestConstructUniverses(t *testing.T) {
	type testCase struct {
		Name string

		Data     *rawConfig
		Expected []*dmx.Universe
		Error    string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			universes, err := tc.Data.constructUniverses()
			if tc.Error != "" {
				assert.EqualError(t, err, tc.Error)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Expected, universes)
			}
		})
	}

	validate(t, &testCase{
		Name: "Valid",

		Data: &rawConfig{
			Universes: map[string]*dmx.Universe{
				"universe": &dmx.Universe{
					Net:    0,
					SubNet: 0,
				},
			},
			Devices: map[string]*dmx.Device{
				"device1": &dmx.Device{
					R: 1,
					G: 2,
					B: 3,
				},
				"device2": &dmx.Device{
					R: 3,
					G: 2,
					B: 1,
				},
			},
			UniversesToDevices: map[string][]string{
				"universe": []string{
					"device1",
					"device2",
				},
			},
		},
		Expected: []*dmx.Universe{
			&dmx.Universe{
				Net:    0,
				SubNet: 0,
				Devices: []*dmx.Device{
					&dmx.Device{
						R: 1,
						G: 2,
						B: 3,
					},
					&dmx.Device{
						R: 3,
						G: 2,
						B: 1,
					},
				},
			},
		},
	})
	validate(t, &testCase{
		Name: "Unknown Universe",

		Data: &rawConfig{
			Universes: map[string]*dmx.Universe{
				"universe1": &dmx.Universe{
					Net:    0,
					SubNet: 0,
				},
			},
			UniversesToDevices: map[string][]string{
				"universe2": []string{},
			},
		},
		Error: "unknown universe: universe2",
	})
	validate(t, &testCase{
		Name: "Unknown Device",

		Data: &rawConfig{
			Universes: map[string]*dmx.Universe{
				"universe": &dmx.Universe{
					Net:    0,
					SubNet: 0,
				},
			},
			Devices: map[string]*dmx.Device{
				"device1": &dmx.Device{
					R: 1,
					G: 2,
					B: 3,
				},
			},
			UniversesToDevices: map[string][]string{
				"universe": []string{
					"device2",
				},
			},
		},
		Error: "unknown device: device2",
	})
}

func TestConstructMapping(t *testing.T) {
	type testCase struct {
		Name string

		Data     *rawConfig
		Expected Mapping
		Error    string
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			mapping, err := tc.Data.constructMapping()
			if tc.Error != "" {
				assert.EqualError(t, err, tc.Error)
			} else {
				assert.NoError(t, err)

				// Comparing maps is ugly.
				assert.Equal(t, len(tc.Expected), len(mapping))
				found := 0
				for keyExpected, expected := range tc.Expected {
					for keyActual, actual := range mapping {
						if reflect.DeepEqual(keyExpected, keyActual) {
							found++
							assert.Equal(t, expected, actual)
						}
					}
				}
				assert.Equal(t, len(tc.Expected), found)
			}
		})
	}

	validate(t, &testCase{
		Name: "Valid",

		Data: &rawConfig{
			Areas: map[string]*image.Rectangle{
				"area": &image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: 0,
					},
					Max: image.Point{
						X: 800,
						Y: 600,
					},
				},
			},
			Devices: map[string]*dmx.Device{
				"device1": &dmx.Device{
					R: 1,
					G: 2,
					B: 3,
				},
				"device2": &dmx.Device{
					R: 3,
					G: 2,
					B: 1,
				},
			},
			AreasToDevices: map[string][]string{
				"area": []string{
					"device1",
					"device2",
				},
			},
		},
		Expected: map[*image.Rectangle][]*dmx.Device{
			&image.Rectangle{
				Min: image.Point{
					X: 0,
					Y: 0,
				},
				Max: image.Point{
					X: 800,
					Y: 600,
				},
			}: []*dmx.Device{
				&dmx.Device{
					R: 1,
					G: 2,
					B: 3,
				},
				&dmx.Device{
					R: 3,
					G: 2,
					B: 1,
				},
			},
		},
	})
	validate(t, &testCase{
		Name: "Unknown Area",

		Data: &rawConfig{
			Areas: map[string]*image.Rectangle{
				"area1": &image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: 0,
					},
					Max: image.Point{
						X: 800,
						Y: 600,
					},
				},
			},
			AreasToDevices: map[string][]string{
				"area2": []string{},
			},
		},
		Error: "unknown area: area2",
	})
	validate(t, &testCase{
		Name: "Unknown Device",

		Data: &rawConfig{
			Areas: map[string]*image.Rectangle{
				"area": &image.Rectangle{
					Min: image.Point{
						X: 0,
						Y: 0,
					},
					Max: image.Point{
						X: 800,
						Y: 600,
					},
				},
			},
			Devices: map[string]*dmx.Device{
				"device1": &dmx.Device{
					R: 1,
					G: 2,
					B: 3,
				},
			},
			AreasToDevices: map[string][]string{
				"area": []string{
					"device2",
				},
			},
		},
		Error: "unknown device: device2",
	})
}
