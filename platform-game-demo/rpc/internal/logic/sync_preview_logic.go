package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/pkg/sync"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

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

// 同步预检查（查看差异，不修改数据）
func (l *SyncPreviewLogic) SyncPreview(in *platformgame.SyncPreviewRequest) (*platformgame.SyncPreviewResp, error) {
	l.Infof("🚀 SyncPreview 请求开始")
	l.Infof("   📋 ObjectType: %s", in.ObjectType)

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Errorf("❌ 数据库连接不可用")
		return &platformgame.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}
	l.Infof("✅ 数据库连接已确认")

	// 从配置读取 grpcServerAddr
	grpcServerAddr := l.svcCtx.Config.GrpcServerAddr
	l.Infof("📍 gRPC 服务器地址配置: %s", grpcServerAddr)

	if grpcServerAddr == "" {
		l.Errorf("❌ GrpcServerAddr 未配置")
		return &platformgame.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "GrpcServerAddr not configured",
		}, nil
	}

	l.Infof("🔄 创建同步服务实例...")
	// 创建同步服务实例
	syncService := sync.NewSyncServiceImpl(l.svcCtx.DB, grpcServerAddr)

	l.Infof("🔗 正在连接 game-vendor-sync 服务...")
	l.Infof("   目标地址: %s", grpcServerAddr)

	// 执行预检查（result 是 *vendors.SyncPreviewResp，直接返回，不是包装在另一个对象中）
	result, err := syncService.Preview(l.ctx, in.ObjectType)
	if err != nil {
		l.Errorf("❌ [SyncPreview] 预检查失败: %v", err)
		return &platformgame.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Preview failed: " + err.Error(),
		}, nil
	}

	if result == nil {
		l.Errorf("❌ [SyncPreview] 结果为空")
		return &platformgame.SyncPreviewResp{
			Code:    constant.CodeInternalError,
			Message: "Empty result from sync service",
		}, nil
	}

	l.Infof("✅ 预检查执行成功")

	// 转换结果为protobuf 消息
	// result 是 *vendors.SyncPreviewResp，包含 Stats 和 Diffs，需要转换为 platformgame 类型
	resp := &platformgame.SyncPreviewResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
	}

	// 转换 Stats
	if result.Stats != nil {
		resp.Stats = &platformgame.SyncStats{
			Total:     int32(result.Stats.RemoteTotal),
			Added:     int32(result.Stats.CreateTotal),
			Updated:   int32(result.Stats.UpdateTotal),
			Deleted:   int32(result.Stats.DeleteTotal),
			Conflicts: int32(result.Stats.ConflictTotal),
		}
		l.Infof("✅ [SyncPreview] 统计数据:")
		l.Infof("   - 总计: %d", result.Stats.RemoteTotal)
		l.Infof("   - 新增: %d", result.Stats.CreateTotal)
		l.Infof("   - 更新: %d", result.Stats.UpdateTotal)
		l.Infof("   - 删除: %d", result.Stats.DeleteTotal)
		l.Infof("   - 冲突: %d", result.Stats.ConflictTotal)
	}

	// 转换 Diffs
	l.Infof("🔍 [SyncPreview] 差异记录: %d 条", len(result.Diffs))
	if result.Diffs != nil && len(result.Diffs) > 0 {
		resp.Diffs = make([]*platformgame.SyncDiff, 0, len(result.Diffs))
		for _, diff := range result.Diffs {
			resp.Diffs = append(resp.Diffs, &platformgame.SyncDiff{
				Type:            diff.ObjectType,
				Action:          diff.Action,
				ObjectId:        fmt.Sprintf("%d", diff.ObjectId),
				Description:     diff.Details,
				ObjectNumericId: diff.ObjectId,
				ObjectCode:      diff.ObjectCode,
				RemoteId:        diff.RemoteId,
				RemoteCode:      diff.RemoteCode,
				Reason:          diff.Reason,
				ConflictType:    diff.ConflictType,
				Details:         diff.Details,
			})
		}
		l.Infof("✅ [SyncPreview] 已添加 %d 条差异记录到响应", len(resp.Diffs))
	} else {
		// 确保 diffs 是空数组而不�?nil
		resp.Diffs = make([]*platformgame.SyncDiff, 0)
		l.Infof("[SyncPreview] 没有差异记录，返回空数组")
	}

	l.Infof("✅ [SyncPreview] 最终响应已准备")
	l.Infof("🎉 SyncPreview 请求完成")
	return resp, nil
}
