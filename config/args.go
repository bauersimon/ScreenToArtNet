package config

import (
	"flag"
	"os"
)

type Args struct {
	Mode      *string
	Src       *string
	Dst       *string
	Pause     *int
	Screen    *int
	Spacing   *int
	Threshold *int
	Config    *string
}

func Parse() Args {
	var args = Args{
		flag.String("mode", "run", "tool mode {run|preview}"),
		flag.String("src", "", "artnet source"),
		flag.String("dst", "", "artnet destination"),
		flag.Int("pause", 0, "pause time in ms"),
		flag.Int("screen", 0, "screen identifier"),
		flag.Int("spacing", 1, "spacing of pixels for averaging"),
		flag.Int("threshold", 0, "threshold of color (0<255)"),
		flag.String("config", "config.json", "config file"),
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
