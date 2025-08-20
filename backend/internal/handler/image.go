package handler

import (
	"fmt"
	"os"
	"path/filepath"

	"pea-blog-backend/pkg/logger"
	"pea-blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ImageHandler struct {
	logger *logger.Logger
}

func NewImageHandler(logger *logger.Logger) *ImageHandler {
	return &ImageHandler{logger: logger}
}

func (h *ImageHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.BadRequest(c, "Invalid file")
		return
	}

	extension := filepath.Ext(file.Filename)
	newFileName := uuid.New().String() + extension
	
	wd, err := os.Getwd()
	if err != nil {
		h.logger.Error("Failed to get working directory", "error", err)
		response.InternalServerError(c, "Failed to upload image")
		return
	}

	uploadPath := filepath.Join(wd, "..", "frontend", "public", "uploads")
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		h.logger.Error("Failed to create upload directory", "error", err)
		response.InternalServerError(c, "Failed to upload image")
		return
	}

	dst := filepath.Join(uploadPath, newFileName)
	if err := c.SaveUploadedFile(file, dst); err != nil {
		h.logger.Error("Failed to save uploaded file", "error", err)
		response.InternalServerError(c, "Failed to upload image")
		return
	}

	imageUrl := fmt.Sprintf("/uploads/%s", newFileName)
	response.Success(c, gin.H{"url": imageUrl})
}