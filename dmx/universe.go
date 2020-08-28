package dmx

import "fmt"

// Universe holds a DMX universe.
type Universe struct {
	// Devices holds the devices of this universe.
	Devices []*Device

	// Net holds the ArtNet net, a group of 16 consecutive sub-nets or 256 consecutive universes.
	Net uint8
	// SubNet holds the ArtNet sub-net, a group of 16 consecutive universes.
	SubNet uint8
}

// Verify checks if the Universe is a valid DMX universe.
func (u *Universe) Verify() error {
	if u.Net > 127 {
		return fmt.Errorf("invalid ArtNet net (net=%v)", u.Net)
	}
	if u.SubNet > 15 {
		return fmt.Errorf("invalid ArtNet subnet (subnet=%v)", u.SubNet)
	}

	return nil
}

// SendColorUpdate sends a color update from the universe devices over the given controller.
func (u *Universe) SendColorUpdate(controller *ArtNetController) error {
	var frame DMXFrame

	for _, d := range u.Devices {
		d.UpdateFrame(&frame)
	}

	return controller.SendDMX(frame, u.Net, u.SubNet)
}
