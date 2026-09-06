package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/response"
	"oa.98ent.com/p9/platform-game/pkg/sync"
)

// SyncHandler 同步相关的 API Handler
type SyncHandler struct {
	syncService sync.SyncService
}

// NewSyncHandler 创建 Sync Handler
func NewSyncHandler(db *gorm.DB, grpcServerAddr string) *SyncHandler {
	return &SyncHandler{
		syncService: sync.NewSyncServiceImpl(db, grpcServerAddr),
	}
}

// SyncCategoryPreviewReq 分类同步预检查请求
type SyncCategoryPreviewReq struct {
	ObjectType string `json:"object_type" binding:"required"`
}

// SyncCategoryPreviewResp 分类同步预检查响应
type SyncCategoryPreviewResp struct {
	Stats interface{} `json:"stats"`
	Diffs interface{} `json:"diffs"`
}

// PreviewCategory 分类同步预检查
func (h *SyncHandler) PreviewCategory(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := h.syncService.Preview(ctx, "category", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// SyncCategoryRunReq 分类同步执行请求
type SyncCategoryRunReq struct {
	AutoApply bool `json:"auto_apply" binding:"required"`
}

// RunCategory 分类同步执行
func (h *SyncHandler) RunCategory(c *gin.Context) {
	var req SyncCategoryRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误: "+err.Error()))
		return
	}

	ctx := c.Request.Context()
	result, err := h.syncService.Run(ctx, "category", nil, req.AutoApply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// PreviewProvider 厂商同步预检查
func (h *SyncHandler) PreviewProvider(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := h.syncService.Preview(ctx, "provider", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// RunProvider 厂商同步执行
func (h *SyncHandler) RunProvider(c *gin.Context) {
	var req SyncCategoryRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误: "+err.Error()))
		return
	}

	ctx := c.Request.Context()
	result, err := h.syncService.Run(ctx, "provider", nil, req.AutoApply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// PreviewChannel 渠道同步预检查
func (h *SyncHandler) PreviewChannel(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := h.syncService.Preview(ctx, "channel", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// RunChannel 渠道同步执行
func (h *SyncHandler) RunChannel(c *gin.Context) {
	var req SyncCategoryRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误: "+err.Error()))
		return
	}

	ctx := c.Request.Context()
	result, err := h.syncService.Run(ctx, "channel", nil, req.AutoApply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// PreviewGame 游戏同步预检查
func (h *SyncHandler) PreviewGame(c *gin.Context) {
	ctx := c.Request.Context()
	result, err := h.syncService.Preview(ctx, "game", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// RunGame 游戏同步执行
func (h *SyncHandler) RunGame(c *gin.Context) {
	var req SyncCategoryRunReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误: "+err.Error()))
		return
	}

	ctx := c.Request.Context()
	result, err := h.syncService.Run(ctx, "game", nil, req.AutoApply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}

// SyncAllReq 全量同步请求
type SyncAllReq struct {
	AutoApply bool `json:"auto_apply" binding:"required"`
}

// SyncAll 全量同步（按依赖关系）
func (h *SyncHandler) SyncAll(c *gin.Context) {
	var req SyncAllReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.Error(400, "参数错误: "+err.Error()))
		return
	}

	ctx := c.Request.Context()

	// 通过强制类型转换调用 SyncAll
	impl, ok := h.syncService.(*sync.SyncServiceImpl)
	if !ok {
		c.JSON(http.StatusInternalServerError, response.Error(500, "服务类型错误"))
		return
	}

	result, err := impl.SyncAll(ctx, req.AutoApply)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.Error(500, err.Error()))
		return
	}
	c.JSON(http.StatusOK, response.Success(result))
}
