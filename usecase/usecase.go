package usecase

import (
	"fmt"
	"io"
	"os"

	"github.com/chrollo-lucifer-12/vod/ffmpeg"
)

type VideoUsecase struct {
	ffmpeg *ffmpeg.FFmpegService
}

func NewVideoUsecase(ffmpeg *ffmpeg.FFmpegService) *VideoUsecase {
	return &VideoUsecase{ffmpeg: ffmpeg}
}

func (uc *VideoUsecase) ProcessAndSave(filename string, r io.Reader) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, r)
	if err != nil {
		return err
	}

	videoDetails, err := uc.ffmpeg.GetVideoDetails(filename)
	if err != nil {
		return err
	}

	if videoDetails == nil {
		return fmt.Errorf("no video details available")
	}

	isPortrait := videoDetails.IsPortrait()
	if err := uc.ffmpeg.Transcode(filename, isPortrait); err != nil {
		return err
	}

	return nil
}
