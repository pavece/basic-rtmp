package config

type Rendition struct {
	Bitrate int
	Height  int
	Width   int
}

var Renditions []Rendition = []Rendition{{Bitrate: 1000, Height: 480, Width: 852}, {Bitrate: 5000, Height: 720, Width: 1280} /*, {Bitrate: 8000, Height: 1080, Width: 1920} */}
