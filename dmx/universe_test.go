package dmx

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUniverseVerify(t *testing.T) {
	type testCase struct {
		Name string

		Universe *Universe
		Error    error
	}

	validate := func(t *testing.T, tc *testCase) {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.Universe.Verify()
			assert.Equal(t, tc.Error, err)
		})
	}

	validate(t, &testCase{
		Name: "Invalid Net",

		Universe: &Universe{
			Net: 128,
		},
		Error: errors.New("invalid ArtNet net (net=128)"),
	})
	validate(t, &testCase{
		Name: "Invalid SubNet",

		Universe: &Universe{
			SubNet: 16,
		},
		Error: errors.New("invalid ArtNet subnet (subnet=16)"),
	})
}
