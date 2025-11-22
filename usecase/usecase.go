package usecase

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/chrollo-lucifer-12/vod/ffmpeg"
)

type VideoUsecase struct {
	ffmpeg *ffmpeg.FFmpegService
}

func NewVideoUsecase(ffmpeg *ffmpeg.FFmpegService) *VideoUsecase {
	return &VideoUsecase{ffmpeg: ffmpeg}
}
func (uc *VideoUsecase) ProcessAndSave(filename string, r io.Reader) error {
	videoPath := filepath.Join("videos", filename)
	file, err := os.Create(videoPath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	videoDetails, err := uc.ffmpeg.GetVideoDetails(videoPath)
	if err != nil {
		return err
	}

	if videoDetails == nil {
		return fmt.Errorf("no video details available")
	}

	fmt.Println("video details", videoDetails)

	isPortrait := videoDetails.IsPortrait()
	if err := uc.ffmpeg.Transcode(videoPath, isPortrait); err != nil {
		return err
	}

	return nil
}
