package config

import (
	"flag"
	"os"
)

type Args struct {
	Mode      *string
	Src       *string
	Dst       *string
	Screen    *int
	Spacing   *int
	Threshold *int
	Config    *string
	Fps       *int
	Workers   *int
}

func Parse() Args {
	var args = Args{
		flag.String("mode", "run", "tool mode {run|preview}"),
		flag.String("src", "", "artnet source"),
		flag.String("dst", "", "artnet destination"),
		flag.Int("screen", 0, "screen identifier"),
		// TODO: use percent of area instead as areas could be of different size
		flag.Int("spacing", 1, "spacing of pixels for averaging"),
		flag.Int("threshold", 0, "threshold of color (0<255)"),
		flag.String("config", "./config.json", "config file"),
		flag.Int("fps", 40, "target frames per second output"),
		flag.Int("workers", 1, "max number of worker threads to use"),
	}
	flag.Parse()
	return args
}

func Validate() bool {
	if len(os.Args) == 1 {
		flag.PrintDefaults()
		return false
	}
	return true
}
