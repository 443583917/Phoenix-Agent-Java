package common

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/response"
)

// FileUploadHandler handles file upload operations.
type FileUploadHandler struct {
}

// NewFileUploadHandler creates a new FileUploadHandler.
func NewFileUploadHandler() *FileUploadHandler {
	return &FileUploadHandler{}
}

// UploadAvatar handles avatar file upload.
// POST /upload/avatar
func (h *FileUploadHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}
	log.Printf("Uploaded avatar: %s (%d bytes)", file.Filename, file.Size)
	response.Success(c, gin.H{"url": "/api/upload/avatar/stub-filename.png"})
}

// GetFile returns file path information for a static file request.
// GET /upload/*filepath
func (h *FileUploadHandler) GetFile(c *gin.Context) {
	filepath := c.Param("filepath")
	response.Success(c, gin.H{"filepath": filepath})
}
