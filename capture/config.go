package capture

// CaptureConfig holds the configuration for the screen capture.
type CaptureConfig struct {
	ImageDataConfig
	// Monitor holds the monitor used for capture.
	Monitor int
}

// imageDataConfig holds the configuration for computing colors
type ImageDataConfig struct {
	// Spacing holds the averaging spacing.
	Spacing int
	// Threshold holds the averaged color threshold.
	Threshold int
}
