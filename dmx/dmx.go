package dmx

import (
	"fmt"
	"image/color"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
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

// Controller holds the network data of the artnet controller.
type Controller struct {
	node *net.UDPAddr
	gate *net.UDPConn
}

// SendColorUpdate sends a color update to the given device with the given artnet controller.
func SendColorUpdate(controller *Controller, device RGBaddress, color color.RGBA) error {
	dmxFrame := [512]byte{0x00, 0x00, 0x00, 0x00, 0x00}
	dmxFrame[device.R] = color.R
	dmxFrame[device.G] = color.G
	dmxFrame[device.B] = color.B

	// TODO How to target the universe.
	p := &packet.ArtDMXPacket{
		SubUni: device.SubNet,
		Net:    device.Universe,
		Data:   dmxFrame,
	}

	b, err := p.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = controller.gate.WriteTo(b, controller.node)
	if err != nil {
		return err
	}

	return nil
}

// NewArtNetController registers a new artnet node on this device to control the given artnet subnet.
func NewArtNetController(src string, dst string) (*Controller, error) {
	nodeAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dst, packet.ArtNetPort))
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", src, packet.ArtNetPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}

	return &Controller{
		node: nodeAddr,
		gate: conn,
	}, nil
}
