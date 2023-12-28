package main

import (
	"flag"
	"fmt"
)

type Args struct {
	Mode      *string
	Src       *string
	Dst       *string
	Screen    *int
	Spacing   *int
	Threshold *int
	Config    *string
	Fps       *uint
	Workers   *uint
}

func Parse() Args {
	var args = Args{
		flag.String("mode", "run", "tool mode {run|preview}"),
		flag.String("src", "", "artnet source"),
		flag.String("dst", "", "artnet destination"),
		flag.Int("screen", 0, "screen identifier"),
		flag.Int("spacing", 1, "spacing of pixels for averaging"),
		flag.Int("threshold", 0, "threshold of color (0<255)"),
		flag.String("config", "./config.json", "config file"),
		flag.Uint("fps", 40, "target frames per second output"),
		flag.Uint("workers", 1, "max number of worker threads to use"),
	}
	flag.Parse()
	return args
}

func (a *Args) Validate() bool {
	valid := false
	if *a.Mode == "run" {
		if *a.Dst == "" {
			fmt.Println("dst must be set to artnet node")
		} else if *a.Spacing < 1 {
			fmt.Println("spacing must be >= 1")
		} else if *a.Fps < 1 {
			fmt.Println("fps must be >= 1")
		} else if *a.Workers < 1 {
			fmt.Println("workers must be >= 1")
		} else {
			valid = true
		}
	} else {
		valid = true
	}

	if !valid {
		flag.PrintDefaults()
	}

	return valid
}
