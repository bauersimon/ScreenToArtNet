package capture

import "github.com/bauersimon/ScreenToArtNet/config"

// CaptureConfig holds the configuration for the screen capture.
type CaptureConfig struct {
	imageDataConfig
	// Monitor holds the monitor used for capture.
	Monitor int
}

// imageDataConfig holds the configuration for computing colors
type imageDataConfig struct {
	// Spacing holds the averaging spacing.
	Spacing int
	// Threshold holds the averaged color threshold.
	Threshold int
}

// TODO: Do the same for other configs
func NewCaptureConfig(args config.Args) CaptureConfig {
	return CaptureConfig{
		Monitor: *args.Screen,
		imageDataConfig: imageDataConfig{
			Spacing:   *args.Spacing,
			Threshold: *args.Threshold,
		},
	}
}
