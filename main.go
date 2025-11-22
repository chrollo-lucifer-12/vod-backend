package main

import (
	"github.com/chrollo-lucifer-12/vod/ffmpeg"
	"github.com/chrollo-lucifer-12/vod/queue"
	"github.com/chrollo-lucifer-12/vod/usecase"
	"github.com/chrollo-lucifer-12/vod/videocontroller"
	"github.com/gin-gonic/gin"
)

func main() {
	ffmpeg_servive := ffmpeg.NewFFmpegService()
	uc := usecase.NewVideoUsecase(ffmpeg_servive)
	q := queue.NewQueue()
	vc := videocontroller.NewVideoController(q)

	go func() {
		for task := range q.ConsumeTasks() {
			file := task.File
			header := task.Header
			filename := header.Filename
			if err := uc.ProcessAndSave(filename, file); err != nil {
				continue
			}
		}
	}()

	r := gin.Default()
	r.POST("/upload", vc.UploadVideo)
	r.Run(":8000")
}
