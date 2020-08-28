package dmx

import (
	"fmt"
	"net"

	"github.com/jsimonetti/go-artnet/packet"
)

// DMXFrame holds a single 512 byte DMX frame.
type DMXFrame [512]byte

// ArtNetController holds the networking data for ArtNet communication.
type ArtNetController struct {
	// node holds the targeted ArtNet node.
	node *net.UDPAddr
	// gate holds the local gateway to use for ArtNet communication.
	gate *net.UDPConn
}

// SendDMX sends the DMX frame to the given net and sub-net.
func (c *ArtNetController) SendDMX(frame DMXFrame, net uint8, sub uint8) error {
	pack := &packet.ArtDMXPacket{
		Net:    net,
		SubUni: sub,
		Data:   frame,
	}

	binary, err := pack.MarshalBinary()
	if err != nil {
		return err
	}

	_, err = c.gate.WriteTo(binary, c.node)
	if err != nil {
		return err
	}

	return nil
}

// NewArtNetController registers all ArtNet communication on this device to control the given ArtNet target node.
func NewArtNetController(srcIP string, dstIP string) (*ArtNetController, error) {
	nodeAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", dstIP, packet.ArtNetPort))
	if err != nil {
		return nil, err
	}

	localAddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", srcIP, packet.ArtNetPort))
	if err != nil {
		return nil, err
	}

	conn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}

	return &ArtNetController{
		node: nodeAddr,
		gate: conn,
	}, nil
}
