package common

import (
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/phoenix-agent-go/infra/errcode"
	"github.com/phoenix-agent-go/infra/id"
	"github.com/phoenix-agent-go/infra/response"
)

const uploadDir = "./storage/upload"

type FileUploadHandler struct{}

func NewFileUploadHandler() *FileUploadHandler {
	_ = os.MkdirAll(uploadDir, 0o755)
	return &FileUploadHandler{}
}

func (h *FileUploadHandler) UploadAvatar(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Error(c, errcode.InvalidParams)
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".gif" && ext != ".webp" {
		response.ErrorWithMsg(c, errcode.InvalidParams, "仅支持 png/jpg/gif/webp 格式")
		return
	}
	if file.Size > 5*1024*1024 {
		response.ErrorWithMsg(c, errcode.InvalidParams, "文件大小不能超过 5MB")
		return
	}

	filename := fmt.Sprintf("avatar_%d%s", id.MustGenerateID(), ext)
	dst := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		response.Error(c, errcode.InternalError)
		return
	}

	url := "/api/upload/" + filename
	response.Success(c, gin.H{"url": url, "filename": filename})
}

func (h *FileUploadHandler) GetFile(c *gin.Context) {
	filePath := c.Param("filepath")
	filePath = strings.TrimPrefix(filePath, "/")

	fullPath := filepath.Join(uploadDir, filepath.Base(filePath))
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		response.Error(c, errcode.NotFound)
		return
	}

	ext := filepath.Ext(fullPath)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	c.Header("Content-Type", contentType)
	c.Header("Cache-Control", "public, max-age=86400")
	http.ServeFile(c.Writer, c.Request, fullPath)
}
