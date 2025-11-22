package ffmpeg

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type FFmpegService struct{}

func NewFFmpegService() *FFmpegService {
	return &FFmpegService{}
}

func (s *FFmpegService) GetVideoDetails(path string) (*VideoData, error) {
	data, err := ffmpeg_go.Probe(path)

	if err != nil {
		return nil, fmt.Errorf("failed to probe video: %w", err)
	}

	var result VideoData
	if err := json.Unmarshal([]byte(data), &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (s *FFmpegService) Transcode(input string, isPortrait bool) error {
	inputDir := strings.Split(input, ".")[0]
	os.Mkdir(inputDir, 0755)

	var wg sync.WaitGroup
	errCh := make(chan error, len(VideoQualities))

	for _, q := range VideoQualities {
		q := q
		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := s.TranscodeQuality(input, inputDir, q, isPortrait); err != nil {
				errCh <- err
			}
		}()
	}

	wg.Wait()

	if len(errCh) > 0 {
		var sb strings.Builder
		for e := range errCh {
			sb.WriteString(e.Error() + "\n")
		}
		return fmt.Errorf(sb.String())
	}

	return s.generateMasterPlaylist(input)

}

func (s *FFmpegService) TranscodeQuality(input, inputDir string, q VideoQuality, isPortrait bool) error {
	qualityDir := filepath.Join(inputDir, "normal_hls", q.Name)
	if err := os.MkdirAll(qualityDir, 0755); err != nil {
		return err
	}

	segmentPath := filepath.Join(qualityDir, "%03d.ts")
	playlistPath := filepath.Join(qualityDir, "index.m3u8")
	scaleFilter := q.ScaleHorizontally()
	if isPortrait {
		scaleFilter = q.ScaleVertically()
	}

	cmd := ffmpeg_go.Input(input).Output(playlistPath, s.getFFmpegArgs(q, segmentPath, []string{scaleFilter, q.LandScape()}))

	err := cmd.OverWriteOutput().WithOutput(nil, os.Stdout).Run()
	if err != nil {
		return err
	}

	return nil
}

func (s *FFmpegService) getFFmpegArgs(q VideoQuality, segmentPath string, filters []string) ffmpeg_go.KwArgs {
	return ffmpeg_go.KwArgs{
		"c:v":                  "h264",
		"profile:v":            "main",
		"crf":                  "20",
		"sc_threshold":         "0",
		"g":                    "48",
		"keyint_min":           "48",
		"b:v":                  q.Bitrate,
		"maxrate":              q.Maxrate,
		"bufsize":              q.Bufsize,
		"c:a":                  "aac",
		"ar":                   "48000",
		"b:a":                  "128k",
		"hls_list_size":        "0",
		"hls_time":             "6",
		"hls_playlist_type":    "vod",
		"start_number":         "1",
		"hls_segment_filename": segmentPath,
		"hls_flags":            "round_durations+split_by_time",
		"hls_allow_cache":      "1",
		"vf":                   filters[0],
		"s":                    filters[1],
	}
}

func (s *FFmpegService) generateMasterPlaylist(input string) error {
	inputDir := strings.TrimSuffix(input, filepath.Ext(input))
	masterFilePath := filepath.Join(inputDir, "master.m3u8")

	masterFile, err := os.Create(masterFilePath)
	if err != nil {
		return err
	}
	defer masterFile.Close()

	writer := bufio.NewWriter(masterFile)
	defer writer.Flush()

	if _, err := writer.WriteString("#EXTM3U\n"); err != nil {
		return err
	}

	for _, q := range VideoQualities {
		bandwidth := extractBandwidth(q.Bitrate)
		line := fmt.Sprintf("#EXT-X-STREAM-INF:BANDWIDTH=%d,RESOLUTION=%s\n%s/index.m3u8\n", bandwidth, q.LandScape(), q.Name)
		if _, err := writer.WriteString(line); err != nil {
			return err
		}
	}

	return nil
}

func extractBandwidth(bitrate string) int {
	if strings.HasSuffix(bitrate, "k") {
		bitrate = strings.TrimSuffix(bitrate, "k")
	}
	kbps, err := strconv.Atoi(bitrate)
	if err != nil {
		return 0
	}
	return kbps * 1000
}
