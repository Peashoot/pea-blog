package handler

import (
	"pea-blog-backend/internal/frontend"
	"pea-blog-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type SystemHandler struct {
	buildService *frontend.BuildService
}

func NewSystemHandler(buildService *frontend.BuildService) *SystemHandler {
	return &SystemHandler{
		buildService: buildService,
	}
}

// RebuildFrontend 手动触发前端重建
func (h *SystemHandler) RebuildFrontend(c *gin.Context) {
	if err := h.buildService.ForceBuild(); err != nil {
		response.InternalServerError(c, "Frontend rebuild failed: "+err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "Frontend rebuilt successfully",
	})
}

// GetBuildStatus 获取构建状态信息
func (h *SystemHandler) GetBuildStatus(c *gin.Context) {
	// 这里可以扩展添加更多状态信息
	response.Success(c, gin.H{
		"auto_build_enabled": h.buildService != nil,
		"message":            "Build service is available",
	})
}