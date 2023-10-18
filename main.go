package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bauersimon/ScreenToArtNet/config"
)

var args config.Args

func handleInterrupts() {
	// Make sure we clean everything up.
	abort := make(chan os.Signal)
	signal.Notify(abort, os.Interrupt, syscall.SIGTERM)
	go func() {
		s := <-abort
		fmt.Printf("\r%v received, stopping...\n", s)
		os.Exit(0)
	}()
}

func executeMode() error {

	var err error = nil
	switch *args.Mode {
	case "run":
		err = run()
	case "preview":
		err = preview()
	default:
		fmt.Printf("unknown mode: %s", *args.Mode)
	}

	return err
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("encountered error:\n%s\n", err.Error())
		os.Exit(-1)
	}
}

func main() {
	args = config.Parse()

	if !config.Validate() {
		os.Exit(-1)
	}

	handleInterrupts()

	handleError(executeMode())
}
