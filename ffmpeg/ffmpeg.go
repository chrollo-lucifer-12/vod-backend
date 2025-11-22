package ffmpeg

import "fmt"

type VideoQuality struct {
	Name    string
	Width   int
	Height  int
	Bitrate string
	Maxrate string
	Bufsize string
}

var VideoQualities = []VideoQuality{
	{"1080p", 1920, 1080, "4500k", "4700k", "6000k"},
	{"720p", 1280, 720, "2500k", "2675k", "3750k"},
	{"480p", 854, 480, "1000k", "1075k", "1500k"},
	{"360p", 640, 360, "600k", "650k", "900k"},
	{"240p", 426, 240, "400k", "450k", "600k"},
	{"144p", 256, 144, "250k", "275k", "400k"},
}

func (vq VideoQuality) ScaleHorizontally() string {
	return fmt.Sprintf("scale=w%d:h=%d:force_original_aspect_ratio=decrease", vq.Width, vq.Height)
}
