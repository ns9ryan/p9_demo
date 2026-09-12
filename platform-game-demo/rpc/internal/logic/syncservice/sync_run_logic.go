package syncservicelogic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	pkgsync "oa.98ent.com/p9/platform-game/pkg/sync"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	sync "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

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
func (l *SyncRunLogic) SyncRun(in *sync.SyncRunRequest) (*sync.SyncRunResp, error) {
	l.Infof("🚀 SyncRun 请求开始")
	l.Infof("   📋 ObjectType: %s", in.ObjectType)
	l.Infof("   ⚙️  SyncCols: %v", in.SyncCols)

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Errorf("❌ 数据库连接不可用")
		return &sync.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}
	l.Infof("✅ 数据库连接已确认")

	// 根据对象类型获取同步范围
	syncScope := l.objectTypeToSyncScope(in.ObjectType)
	l.Infof("📍 同步范围: %s", syncScope)

	// 从配置读取 grpcServerAddr
	grpcServerAddr := l.svcCtx.Config.GrpcServerAddr
	l.Infof("📍 gRPC 服务器地址配置: %s", grpcServerAddr)

	// if grpcServerAddr == "" {
	// 	l.Errorf("❌ GrpcServerAddr 未配置")
	// 	l.Infof("💡 请在配置文件中设置 GrpcServerAddr，例子如下")
	// 	l.Infof("   GrpcServerAddr: game-vendor-sync:19009")
	// 	return &sync.SyncRunResp{
	// 		Code:    constant.CodeInternalError,
	// 		Message: "GrpcServerAddr not configured",
	// 	}, nil
	// }

	// 第一步：创建同步检查点记录（RPC 侧创建）
	checkpoint := &ent.GameSyncCheckpoint{
		SyncScope:       syncScope,
		SyncStatus:      1, // 进行中
		CheckpointValue: "0",
		Progress:        0, // 初始进度为 0
		RemoteTotal:     0,
		LocalTotal:      0,
		CreatedCount:    0,
		UpdatedCount:    0,
		DeletedCount:    0,
		FailedCount:     0,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if err := l.svcCtx.DB.Create(checkpoint).Error; err != nil {
		l.Errorf("[RPC SyncRun] create checkpoint failed: %v", err)
		return &sync.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "create checkpoint failed: " + err.Error(),
		}, nil
	}

	l.Infof("[RPC SyncRun] checkpoint created: id=%d", checkpoint.Id)

	// 第二步：立即返回检查点 ID，不阻塞
	resp := &sync.SyncRunResp{
		Code:         constant.CodeSuccess,
		Message:      "async sync started",
		CheckpointId: checkpoint.Id,
	}

	// 第三步：启动协程异步处理同步逻辑
	go l.doAsyncSync(checkpoint.Id, in.ObjectType, in.SyncCols, grpcServerAddr, syncScope, checkpoint.Id)

	l.Infof("[RPC SyncRun] async sync started for checkpoint: id=%d", checkpoint.Id)
	l.Infof("🎉 SyncRun 请求完成")
	return resp, nil
}

// doAsyncSync 异步执行同步逻辑（在 RPC 侧）
func (l *SyncRunLogic) doAsyncSync(checkpointID int64, objectType string, syncCols []string, grpcServerAddr string, syncScope string, localCheckpointID int64) {
	// 创建新的上下文用于异步操作，不依赖请求上下文
	ctx := context.Background()
	logger := logx.WithContext(ctx)

	logger.Infof("[RPC AsyncSync] starting async sync for checkpoint: id=%d, type=%s", checkpointID, objectType)

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("[RPC AsyncSync] panic: %v", r)
			errMsg := fmt.Sprintf("panic: %v", r)
			l.updateCheckpointError(checkpointID, errMsg)
		}
	}()

	// 创建同步服务实例
	syncService := pkgsync.NewSyncServiceImpl(l.svcCtx.DB, grpcServerAddr)

	// 执行同步
	result, err := syncService.Run(ctx, objectType, syncCols, localCheckpointID)
	if err != nil {
		logger.Errorf("[RPC AsyncSync] sync run failed: %v", err)
		l.updateCheckpointError(checkpointID, fmt.Sprintf("sync run failed: %v", err))
		return
	}

	if result == nil {
		logger.Error("[RPC AsyncSync] sync result is nil")
		l.updateCheckpointError(checkpointID, "sync result is nil")
		return
	}

	logger.Infof("[RPC AsyncSync] checkpoint updated successfully: checkpoint_id=%d", checkpointID)
}

// updateCheckpointError 更新检查点错误状态
func (l *SyncRunLogic) updateCheckpointError(checkpointID int64, errMsg string) {
	checkpoint := &ent.GameSyncCheckpoint{
		SyncStatus:       0, // 失败
		LastErrorMessage: sql.NullString{String: errMsg, Valid: true},
		UpdatedAt:        time.Now(),
	}

	if err := l.svcCtx.DB.Model(&ent.GameSyncCheckpoint{}).Where("id = ?", checkpointID).Updates(checkpoint).Error; err != nil {
		logx.Errorf("[RPC AsyncSync] update checkpoint error failed: %v", err)
	}
}

// objectTypeToSyncScope 根据对象类型转换为同步范围
func (l *SyncRunLogic) objectTypeToSyncScope(objectType string) string {
	switch objectType {
	case "category":
		return constant.SyncScopeCategory
	case "provider":
		return constant.SyncScopeProvider
	case "channel":
		return constant.SyncScopeChannel
	case "game":
		return constant.SyncScopeGame
	case "currency":
		return constant.SyncScopeGameCurrency
	default:
		return objectType
	}
}
