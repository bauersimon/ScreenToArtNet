package dmx

import (
	"fmt"
	"image/color"
)

// RGBaddress holds the artnet address of an rgb device.
type RGBaddress struct {
	R uint16 `json:"red"`
	G uint16 `json:"green"`
	B uint16 `json:"blue"`

	Net      uint8 `json:"net"`
	SubNet   uint8 `json:"subnet"`
	Universe uint8 `json:"universe"`
}

// Verify checks if the RGBaddress is a valid artnet device.
func (addr RGBaddress) Verify() error {
	if addr.R > 511 {
		return fmt.Errorf("rgb device red channel outside of dmx frame (channel=%v)", addr.R)
	}
	if addr.G > 511 {
		return fmt.Errorf("rgb device green channel outside of dmx frame (channel=%v)", addr.G)
	}
	if addr.B > 511 {
		return fmt.Errorf("rgb device blue channel outside of dmx frame (channel=%v)", addr.B)
	}

	if !(addr.R != addr.G && addr.G != addr.B && addr.R != addr.B) {
		return fmt.Errorf("channels should be different (r=%v, g=%v, b=%v)", addr.R, addr.G, addr.B)
	}

	if addr.Net > 127 {
		return fmt.Errorf("invalid artnet net (net=%v)", addr.Net)
	}
	if addr.SubNet > 15 {
		return fmt.Errorf("invalid artnet subnet (subnet=%v)", addr.SubNet)
	}
	if addr.Universe > 15 {
		return fmt.Errorf("invalid artnet universe (universe=%v)", addr.Universe)
	}

	return nil
}

// SendColorUpdate sends a color update to the given device with the given artnet controller.
func SendColorUpdate(controller *ArtNetController, device RGBaddress, color color.RGBA) error {
	var frame DmxFrame
	frame[device.R] = color.R
	frame[device.G] = color.G
	frame[device.B] = color.B

	return controller.SendDmx(frame, device.Net, device.SubNet)
}
