package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/pkg/sync"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type SyncRunLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSyncRunLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SyncRunLogic {
	return &SyncRunLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 同步执行（执行同步操作）
func (l *SyncRunLogic) SyncRun(in *platformgame.SyncRunRequest) (*platformgame.SyncRunResp, error) {
	l.Infof("🚀 SyncRun 请求开始")
	l.Infof("   📋 ObjectType: %s", in.ObjectType)
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
		l.Infof("💡 请在配置文件中设置 GrpcServerAddr，例子如下")
		l.Infof("   GrpcServerAddr: game-vendor-sync:19009")
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
	l.Infof("   预期服务: VendorGameService")

	// 执行同步（注意：当前 Run 方法只接�?objectType 参数，不支持 autoApply 参数�?
	// 如果需要支�?autoApply，需要修�?pkg/sync 中的接口
	result, err := syncService.Run(l.ctx, in.ObjectType)
	if err != nil {
		l.Errorf("❌ [SyncRun] 同步执行失败")
		l.Errorf("   错误详情: %v", err)
		l.Errorf("   错误类型: %T", err)
		l.Errorf("   💡 排查步骤:")
		l.Errorf("      1. 检�?game-vendor-sync 服务是否正在运行")
		l.Errorf("      2. 确认 GrpcServerAddr 配置正确: %s", grpcServerAddr)
		l.Errorf("      3. 查看 game-vendor-sync 的启动日志")
		l.Errorf("      4. 尝试手动测试连接: telnet %s", grpcServerAddr)
		return &platformgame.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "Sync run failed: " + err.Error(),
		}, nil
	}

	if result == nil {
		l.Errorf("❌ [SyncRun] 结果为空")
		return &platformgame.SyncRunResp{}, nil
	}

	l.Infof("✅ 同步执行成功")
	l.Infof("   结果统计:")
	l.Infof("     - Created: %d", result.Apply.Created)
	l.Infof("     - Updated: %d", result.Apply.Updated)
	l.Infof("     - Deleted: %d", result.Apply.Deleted)
	l.Infof("     - Failed: %d", result.Apply.Failed)
	l.Infof("     - Skipped: %d", result.Apply.Skipped)

	// 转换结果�?protobuf 消息
	// result �?*vendors.SyncRunResp，其 Apply 字段包含 Created, Updated, Deleted, Failed, Skipped
	resp := &platformgame.SyncRunResp{}

	if result.Apply != nil {
		total := result.Apply.Created + result.Apply.Updated + result.Apply.Deleted
		success := result.Apply.Created + result.Apply.Updated + result.Apply.Deleted
		failed := result.Apply.Failed

		l.Infof("[SyncRun] 转换结果�?protobuf 消息")
		l.Infof("   - 总计: %d", total)
		l.Infof("   - 成功: %d", success)
		l.Infof("   - 失败: %d", failed)

		resp.Total = total
		resp.Success = success
		resp.Failed = failed
		resp.Created = result.Apply.Created
		resp.Updated = result.Apply.Updated
		resp.Deleted = result.Apply.Deleted
		resp.Skipped = result.Apply.Skipped
	}

	l.Infof("✅ [SyncRun] 最终响应 %+v", resp)
	l.Infof("🎉 SyncRun 请求完成")
	return resp, nil
}
