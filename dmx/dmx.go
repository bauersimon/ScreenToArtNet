package dmx

import (
	"fmt"
)

// Device holds the DMX data of an rgb device.
type Device struct {
	// R holds the red channel.
	R uint16 `json:"Red"`
	// G holds the green channel.
	G uint16 `json:"Green"`
	// B holds the blue channel.
	B uint16 `json:"Blue"`

	// RValue holds the red value.
	RValue uint8
	// RValue holds the green value.
	GValue uint8
	// RValue holds the blue value.
	BValue uint8

	// Statics holds the static DMX data for this device
	Statics map[uint16]uint8
}

// Verify checks if the Device is a valid DMX device.
func (d *Device) Verify() error {
	if d.R > 511 {
		return fmt.Errorf("red channel outside of DMX range (channel=%v)", d.R)
	}
	if d.G > 511 {
		return fmt.Errorf("green channel outside of DMX range (channel=%v)", d.G)
	}
	if d.B > 511 {
		return fmt.Errorf("blue channel outside of DMX range (channel=%v)", d.B)
	}

	if !(d.R != d.G && d.G != d.B && d.R != d.B) {
		return fmt.Errorf("color channels should be different (r=%v, g=%v, b=%v)", d.R, d.G, d.B)
	}

	for channel := range d.Statics {
		if channel > 511 {
			return fmt.Errorf("invalid static channel outside of DMX range (channel=%v)", channel)
		}
	}

	return nil
}

// UpdateFrame updates the given DMX frame with the current channel values.
func (d *Device) UpdateFrame(frame *DMXFrame) {
	frame[d.R] = d.RValue
	frame[d.G] = d.GValue
	frame[d.B] = d.BValue

	for channel, value := range d.Statics {
		frame[channel] = value
	}
}
