# Welcome to the ScreenToArtNet wiki!

This tool utilizes [kbinani/screenshot](https://github.com/kbinani/screenshot) to capture your screen, computes average color values and sends them to an ArtNet node via [jsimonetti/go-artnet](https://github.com/jsimonetti/go-artnet). This makes it possible to set up an ambilight system similar to "Philips Hue" (if you have the spare parts lying around of course).

## Installation:
Firstly, install the Golang language and compilers from https://go.dev/
Then, once installed, download and extract the repository (or clone it using git) and open a new terminal or command prompt window in the downloaded directory. Ensure Go has been installed correctly by typing `go` into the terminal and checking the response. If you get a `command not found` error. Try rebooting your PC. 

To create a Windows Executable file, simply type: <br />
`go build`

You should then see a `ScreenToArtNet.exe` file appear in the directory.


## Using the program:
ScreenToArtNet is a command-line tool. This means that it does not have a GUI (graphical user interface). Firstly, open `config.json` and add all your ArtNet DMX devices. **NOTE: Currently ScreenToArtNet can only output on Universe 1.** An example of the config file is provided. 

Next, in the terminal run:<br />
`ScreenToArtNet.exe -config <path-to-config-file> -dst <artnet-node-ip-address>`

This will start the program and send colour values to your configured DMX devices. 


## Configuration Options:
`-config` - Config file path. Default "config.json"<br />
`-dst` - Destination IP address. If you have multiple devices you wish to send ArtNet to, use the broadcast address `255.255.255.255`.<br />
`-mode` - Tool mode. Default is `run` but enter `preview` to run the tool without ArtNet output. <br />
`-pause` - Interval time to take between processing the screen snapshot (in milliseconds). <br />
`-screen` - If you have multiple screens, use this to identify which screen should be captured. <br />
`-spacing` - Spacing of colour pixels, used to increase update rate. Defaults to 1. <br />
`-src` - ArtNet source. <br />
`-threshold` - Colour threshold. (integer between 0 and 255). 
