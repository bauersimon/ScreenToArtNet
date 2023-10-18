package ambilight

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"path"

	"github.com/bauersimon/ScreenToArtNet/capture"
	"github.com/bauersimon/ScreenToArtNet/dmx"
)

// rawConfig holds the complete raw configuration structure.
type rawConfig struct {
	// Areas holds area names and their respective areas.
	Areas map[string]*image.Rectangle
	// Universes hold universe names and their respective DMX universes.
	Universes map[string]*dmx.Universe
	// Devices hold device names and their respective devices.
	Devices map[string]*dmx.Device

	// UniversesToDevices maps universe names to multiple device names.
	UniversesToDevices map[string][]string
	// AreasToDevices maps area names to multiple device names.
	AreasToDevices map[string][]string
}

// ReadConfig reads the given config file.
func ReadConfig(configPath string) (areas []capture.Area, universes []*dmx.Universe, mapping Mapping, err error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, nil, nil, err
	}

	data, err := ioutil.ReadFile(path.Join(cwd, configPath))
	if err != nil {
		return nil, nil, nil, err
	}

	raw, err := parseConfig(data)
	if err != nil {
		return nil, nil, nil, err
	}

	mapping, err = raw.constructMapping()
	if err != nil {
		return nil, nil, nil, err
	}

	for areaName, areaRect := range raw.Areas {
		areas = append(areas, capture.Area{Name: areaName, ImageData: capture.ImageData{Borders: *areaRect}})
	}

	universes, err = raw.constructUniverses()
	if err != nil {
		return nil, nil, nil, err
	}

	return areas, universes, mapping, nil
}

func (r *rawConfig) constructUniverses() (universes []*dmx.Universe, err error) {
	for universeName, deviceNames := range r.UniversesToDevices {
		u, ok := r.Universes[universeName]
		if !ok {
			return nil, fmt.Errorf("unknown universe: %s", universeName)
		}

		u.Devices, err = r.getDevicesNamed(deviceNames)
		if err != nil {
			return nil, err
		}
		universes = append(universes, u)
	}

	return universes, nil
}

func (r *rawConfig) constructMapping() (mapping Mapping, err error) {
	mapping = make(Mapping)
	for areaName, deviceNames := range r.AreasToDevices {
		a, ok := r.Areas[areaName]
		if !ok {
			return nil, fmt.Errorf("unknown area: %s", areaName)
		}

		mapping[a], err = r.getDevicesNamed(deviceNames)
		if err != nil {
			return nil, err
		}
	}

	return mapping, nil
}

func (r *rawConfig) getDevicesNamed(names []string) ([]*dmx.Device, error) {
	devices := make([]*dmx.Device, len(names))
	for i, deviceName := range names {
		d, ok := r.Devices[deviceName]
		if !ok {
			return nil, fmt.Errorf("unknown device: %s", deviceName)
		}

		devices[i] = d
	}

	return devices, nil
}

func parseConfig(data []byte) (*rawConfig, error) {
	var config rawConfig
	err := json.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
