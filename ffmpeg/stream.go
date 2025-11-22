package ffmpeg

type Stream struct {
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	CodecType string `json:"codec_type"`
}

type VideoData struct {
	Streams []Stream `json:"streams"`
}

func (v VideoData) IsPortrait() bool {
	for _, stream := range v.Streams {
		if stream.CodecType == "video" {
			return stream.Height > stream.Width
		}
	}
	return false
}
