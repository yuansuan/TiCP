package zone

import (
	"github.com/gin-gonic/gin"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/common"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/config"
	"github.com/yuansuan/ticp/iPaaS/project-root/internal/job/service/v1/zone"
)

// Handler ...
type Handler struct {
	cfg config.CustomT
}

// NewZoneHandler ...
func NewZoneHandler(cfg config.CustomT) *Handler {
	return &Handler{
		cfg: cfg,
	}
}

// List 获取所有区域
func (h *Handler) List(c *gin.Context) {
	// req not need
	resp := zone.List(h.cfg)
	// 在这里执行获取指定计算区域的逻辑
	common.SuccessResp(c, resp)
}
