package videocontroller

import (
	"net/http"

	"github.com/chrollo-lucifer-12/vod/queue"
	"github.com/gin-gonic/gin"
)

type VideoController struct {
	q *queue.Queue
}

func NewVideoController(q *queue.Queue) *VideoController {
	return &VideoController{
		q: q,
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

	v.q.AddTask(file, header)

	c.JSON(http.StatusOK, gin.H{"message": "Video uploaded successfully", "filename": filename})
}
