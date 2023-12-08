I don't have the equipment anymore to maintain this project. Please have a look at [https://github.com/tuffnerdstuff/ScreenToArtNet](tuffnerdstuff)'s fork, who took the time to address some of the remaining issues.

---

# ScreenToArtNet

This tool utilizes [kbinani/screenshot](https://github.com/kbinani/screenshot) to capture your screen, computes average color values and sends them to an ArtNet node via [jsimonetti/go-artnet](https://github.com/jsimonetti/go-artnet). This makes it possible to set up an ambilight system similar to "Philips Hue" (if you have the spare parts lying around of course).

# Disclaimer

- should work on Linux, Windows and MacOS (though I only tested it on Linux and Windows)
- not nearly finished, I still have some ideas and performance optimizations in mind