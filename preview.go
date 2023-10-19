package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"github.com/bauersimon/ScreenToArtNet/ambilight"
	"github.com/bauersimon/ScreenToArtNet/capture"
)

func preview() error {
	areas, _, _, err := ambilight.ReadConfig(*args.Config)
	if err != nil {
		return err
	}

	s := capture.NewScreen(
		areas,
		capture.NewCaptureConfig(args),
	)

	// capture screen
	err = s.Capture()
	if err != nil {
		return err
	}
	// save monitor image
	s.ImageData.Update()
	monitorImage, err := s.ImageData.GetImage()
	if err != nil {
		return err
	}
	err = saveArea("monitor.png", monitorImage)
	if err != nil {
		return err
	}

	// save area images
	for _, a := range s.Areas {
		a.ImageData.Update()
		areaImage, err := a.ImageData.GetImage()
		if err != nil {
			return err
		}
		err = saveArea(fmt.Sprintf("area_%s.png", a.Name), areaImage)
		if err != nil {
			return err
		}
	}

	return nil
}

func saveArea(filename string, area *image.RGBA) error {

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// create output dir
	dirPath := filepath.Join(cwd, "preview")
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return err
	}

	// create output file
	filePath := filepath.Join(dirPath, filename)
	outputFile, err := os.Create(filePath)
	if err != nil {
		return err
	}

	// write file
	err = png.Encode(outputFile, area)
	if err != nil {
		return err
	}

	return outputFile.Close()
}
