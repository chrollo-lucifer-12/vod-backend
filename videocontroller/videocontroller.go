package videocontroller

import (
	"net/http"

	"github.com/chrollo-lucifer-12/vod/usecase"
	"github.com/gin-gonic/gin"
)

type VideoController struct {
	uc *usecase.VideoUsecase
}

func NewVideoController(uc *usecase.VideoUsecase) *VideoController {
	return &VideoController{
		uc: uc,
	}
}

func (v *VideoController) UploadVideo(c *gin.Context) {
	file, header, err := c.Request.FormFile("video")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Video file is required"})
		return
	}
	defer file.Close()

	filename := header.Filename

	if err := v.uc.ProcessAndSave(filename, file); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully", "filename": filename})
}
