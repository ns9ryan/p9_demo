package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/pkg/sync"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncAllLogic {
	return &SyncAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 全量同步（同步所有对象类型）
func (l *SyncAllLogic) SyncAll(in *platformgame.SyncAllRequest) (*platformgame.SyncRunResp, error) {
	l.Infof("🚀 SyncAll 请求开始")
	l.Infof("   ⚙️  AutoApply: %v", in.AutoApply)

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Errorf("❌ 数据库连接不可用")
		return &platformgame.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}
	l.Infof("✅ 数据库连接已确认")

	// 从配置读�?grpcServerAddr
	grpcServerAddr := l.svcCtx.Config.GrpcServerAddr
	l.Infof("📍 gRPC 服务器地址配置: %s", grpcServerAddr)

	if grpcServerAddr == "" {
		l.Errorf("❌ GrpcServerAddr 未配置")
		return &platformgame.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "GrpcServerAddr not configured",
		}, nil
	}

	l.Infof("🔄 创建同步服务实例...")
	// 创建同步服务实例
	syncService := sync.NewSyncServiceImpl(l.svcCtx.DB, grpcServerAddr)

	l.Infof("🔗 正在连接 game-vendor-sync 服务...")
	l.Infof("   目标地址: %s", grpcServerAddr)

	// 执行全量同步（注意：当前 SyncAll 方法不支�?autoApply 参数�?
	result, err := syncService.SyncAll(l.ctx)
	if err != nil {
		l.Errorf("�?[SyncAll] 全量同步失败: %v", err)
		return &platformgame.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "Sync all failed: " + err.Error(),
		}, nil
	}

	if result == nil {
		l.Errorf("�?[SyncAll] 结果为空")
		return &platformgame.SyncRunResp{}, nil
	}

	l.Infof("✅ 全量同步执行成功")

	// 转换结果�?protobuf 消息
	// result �?*vendors.SyncRunResp，其 Apply 字段包含 Created, Updated, Deleted, Failed, Skipped
	resp := &platformgame.SyncRunResp{}

	if result.Apply != nil {
		total := result.Apply.Created + result.Apply.Updated + result.Apply.Deleted
		success := result.Apply.Created + result.Apply.Updated + result.Apply.Deleted
		failed := result.Apply.Failed

		l.Infof("�?[SyncAll] 同步结果:")
		l.Infof("   - 总计: %d", total)
		l.Infof("   - 成功: %d", success)
		l.Infof("   - 失败: %d", failed)

		resp.Total = total
		resp.Success = success
		resp.Failed = failed
	}

	l.Infof("🎉 SyncAll 请求完成")
	return resp, nil
}
