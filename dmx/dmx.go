package dmx

import (
	"fmt"
	"image/color"
)

// Device holds the artnet address of an rgb device.
type Device struct {
	// R holds the red channel.
	R uint16 `json:"red"`
	// G holds the green channel.
	G uint16 `json:"green"`
	// B holds the blue channel.
	B uint16 `json:"blue"`

	// Statics holds the static DMX data for this device
	Statics map[uint16]uint8 `json:"statics"`

	// Net holds the ArtNet net, a group of 16 consecutive sub-nets or 256 consecutive universes.
	Net uint8 `json:"net"`
	// SubNet holds the ArtNet sub-net, a group of 16 consecutive universes.
	SubNet uint8 `json:"subnet"`
}

// Verify checks if the Device is a valid ArtNet device.
func (d *Device) Verify() error {
	if d.R > 511 {
		return fmt.Errorf("rgb device red channel outside of DMX frame (channel=%v)", d.R)
	}
	if d.G > 511 {
		return fmt.Errorf("rgb device green channel outside of DMX frame (channel=%v)", d.G)
	}
	if d.B > 511 {
		return fmt.Errorf("rgb device blue channel outside of DMX frame (channel=%v)", d.B)
	}

	if !(d.R != d.G && d.G != d.B && d.R != d.B) {
		return fmt.Errorf("channels should be different (r=%v, g=%v, b=%v)", d.R, d.G, d.B)
	}

	if d.Net > 127 {
		return fmt.Errorf("invalid ArtNet net (net=%v)", d.Net)
	}
	if d.SubNet > 15 {
		return fmt.Errorf("invalid ArtNet subnet (subnet=%v)", d.SubNet)
	}

	for channel, value := range d.Statics {
		if channel > 511 {
			return fmt.Errorf("invalid static channel outside of DMX frame (channel=%v)", channel)
		} else if value > 255 {
			return fmt.Errorf("invalid static channel value of channel %v (channel=%v)", channel, value)
		}
	}

	return nil
}

// SendColorUpdate sends a color update to the given device with the given artnet controller.
func SendColorUpdate(controller *ArtNetController, device *Device, color color.RGBA) error {
	var frame DMXFrame
	frame[device.R] = color.R
	frame[device.G] = color.G
	frame[device.B] = color.B

	for channel, value := range device.Statics {
		frame[channel] = value
	}

	return controller.SendDMX(frame, device.Net, device.SubNet)
}
