package dmx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeviceVerify(t *testing.T) {
	type testCase struct {
		Name string

		Device *Device
		Error  error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.Device.Verify()
			assert.Equal(t, tc.Error, err)
		})
	}

	t.Run("Color Channels", func(t *testing.T) {
		validate(t, &testCase{
			Name: "Red",

			Device: &Device{
				R: 512,
				G: 1,
				B: 2,
			},
			Error: errors.New("red channel outside of DMX range (channel=512)"),
		})
		validate(t, &testCase{
			Name: "Green",

			Device: &Device{
				R: 1,
				G: 512,
				B: 2,
			},
			Error: errors.New("green channel outside of DMX range (channel=512)"),
		})
		validate(t, &testCase{
			Name: "Blue",

			Device: &Device{
				R: 1,
				G: 2,
				B: 512,
			},
			Error: errors.New("blue channel outside of DMX range (channel=512)"),
		})
		validate(t, &testCase{
			Name: "Equal Channels",

			Device: &Device{
				R: 0,
				G: 0,
				B: 0,
			},
			Error: errors.New("color channels should be different (r=0, g=0, b=0)"),
		})
	})
	validate(t, &testCase{
		Name: "Static",

		Device: &Device{
			R: 1,
			G: 2,
			B: 3,

			Statics: map[uint16]uint8{
				512: 0,
			},
		},
		Error: errors.New("invalid static channel outside of DMX range (channel=512)"),
	})
}

func TestDeviceUpdateFrame(t *testing.T) {
	type testCase struct {
		Name string

		Device *Device
		Frame  DMXFrame
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			var frame DMXFrame
			tc.Device.UpdateFrame(&frame)
			assert.Equal(t, tc.Frame, frame)
		})
	}

	validate(t, &testCase{
		Name: "Colors",

		Device: &Device{
			R: 1,
			G: 2,
			B: 3,

			RValue: 1,
			GValue: 2,
			BValue: 3,
		},
		Frame: [512]byte{0, 1, 2, 3},
	})
	validate(t, &testCase{
		Name: "Statics",

		Device: &Device{
			R: 2,
			G: 3,
			B: 4,

			RValue: 0,
			GValue: 0,
			BValue: 0,

			Statics: map[uint16]uint8{
				1: 1,
			},
		},
		Frame: [512]byte{0, 1},
	})
}
