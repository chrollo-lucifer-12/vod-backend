package main

import (
	"github.com/chrollo-lucifer-12/vod/ffmpeg"
	"github.com/chrollo-lucifer-12/vod/usecase"
	"github.com/chrollo-lucifer-12/vod/videocontroller"
	"github.com/gin-gonic/gin"
)

func main() {
	ffmpeg_servive := ffmpeg.NewFFmpegService()
	uc := usecase.NewVideoUsecase(ffmpeg_servive)
	vc := videocontroller.NewVideoController(uc)

	r := gin.Default()
	r.POST("/upload", vc.UploadVideo)
	r.Run(":8000")
}
