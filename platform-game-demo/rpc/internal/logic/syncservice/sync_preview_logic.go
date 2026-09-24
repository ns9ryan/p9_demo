package syncservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	gs "oa.98ent.com/p9/platform-game/rpc/internal/synchro"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncPreviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncPreviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncPreviewLogic {
	return &SyncPreviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SyncPreviewLogic) SyncPreview(in *platform_game.SyncPreviewRequest) (*platform_game.SyncPreviewResp, error) {
	l.Infof("🚀 ========== SyncPreview 请求开始 ==========")
	// 检查入参合法性
	if in == nil {
		l.Errorf("❌ 请求对象为空")
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Request is nil",
		}, nil
	}

	// 检查数据库连接
	if l.svcCtx == nil {
		l.Errorf("❌ ServiceContext 为空")
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "ServiceContext not available",
		}, nil
	}

	if l.svcCtx.DB == nil {
		l.Errorf("❌ 数据库连接不可用")
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	// 检查 DAOManager 是否初始化
	if l.svcCtx.DAOManager == nil {
		l.Errorf("❌ DAO 管理器未初始化")
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not initialized",
		}, nil
	}

	// 创建同步服务实例
	syncService := gs.NewSyncServiceImpl(l.ctx, l.svcCtx.Config, l.svcCtx.DAOManager)

	// 检查 syncService 是否正确初始化
	if syncService == nil {
		l.Errorf("❌ 同步服务实例创建失败（返回 nil）")
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Failed to create sync service",
		}, nil
	}
	l.Infof("✅ 同步服务实例创建成功")

	// 执行预检查（result 是 *vendors.SyncPreviewResp，直接返回，不是包装在另一个对象中）
	result, err := syncService.Preview(l.ctx, in)
	if err != nil {
		l.Errorf("❌ [SyncPreview] 预检查失败: %v", err)
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Preview failed: " + err.Error(),
		}, nil
	}

	if result == nil {
		l.Errorf("❌ [SyncPreview] 结果为空")
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Empty result from sync service",
		}, nil
	}

	l.Infof("✅ 预检查执行成功")
	result.Code = constant.CodeSuccess
	result.Message = "ok"
	l.Infof("🚀 ========== SyncPreview 请求完成 ==========")

	return result, nil
}
