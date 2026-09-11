package syncservicelogic

import (
	"context"

	pkgsync "oa.98ent.com/p9/platform-game/pkg/sync"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
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
	l.Infof("🚀 SyncPreview 请求开始")
	l.Infof("   📋 ObjectType: %s", in.ObjectType)

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Errorf("❌ 数据库连接不可用")
		return &platform_game.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}
	l.Infof("✅ 数据库连接已确认")

	// 从配置读取 grpcServerAddr
	grpcServerAddr := l.svcCtx.Config.GrpcServerAddr
	l.Infof("📍 gRPC 服务器地址配置: %s", grpcServerAddr)

	// if grpcServerAddr == "" {
	// 	l.Errorf("❌ GrpcServerAddr 未配置")
	// 	return &platform_game.SyncPreviewResp{
	// 		Code:    constant.CodeInternalError,
	// 		Message: "GrpcServerAddr not configured",
	// 	}, nil
	// }

	l.Infof("🔄 创建同步服务实例...")
	// 创建同步服务实例
	syncService := pkgsync.NewSyncServiceImpl(l.svcCtx.DB, grpcServerAddr)

	l.Infof("🔗 正在连接 game-vendor-sync 服务...")
	l.Infof("   目标地址: %s", grpcServerAddr)

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

	return result, nil
}
