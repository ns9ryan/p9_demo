package syncservicelogic

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/config"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	gs "oa.98ent.com/p9/platform-game/rpc/internal/synchro"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
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
func (l *SyncRunLogic) SyncRun(in *platform_game.SyncRunRequest) (*platform_game.SyncRunResp, error) {
	l.Infof("🚀 SyncRun 请求开始")

	// 检查数据库连接
	if l.svcCtx.DB == nil {
		l.Errorf("❌ 数据库连接不可用")
		return &platform_game.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	// 根据对象类型获取同步范围
	syncScope := l.objectTypeToSyncScope(in.ObjectType)

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

	createdCheckpoint, err := l.svcCtx.DAOManager.GameSyncCheckpoint.CreateGameSyncCheckpoint(l.ctx, checkpoint)
	if err != nil {
		l.Errorf("[RPC SyncRun] create checkpoint failed: %v", err)
		return &platform_game.SyncRunResp{
			Code:    constant.CodeInternalError,
			Message: "create checkpoint failed: " + err.Error(),
		}, nil
	}

	l.Infof("[RPC SyncRun] 检查点已创建: id=%d", createdCheckpoint.ID)

	// 第二步：立即返回检查点 ID，不阻塞
	resp := &platform_game.SyncRunResp{
		Code:         constant.CodeSuccess,
		Message:      "async sync started",
		CheckpointId: createdCheckpoint.ID,
	}
	newCtx := utils.CloneCtxWithTraceID(l.ctx)
	// 第三步：启动协程异步处理同步逻辑
	go l.doAsyncSync(newCtx, l.svcCtx.Config, createdCheckpoint.ID, in.ObjectType, in.SyncCols, syncScope, createdCheckpoint.ID)

	l.Infof("[RPC SyncRun] async sync started for checkpoint: id=%d", checkpoint.ID)
	l.Infof("🎉 SyncRun 请求完成")
	return resp, nil
}

// doAsyncSync 异步执行同步逻辑（在 RPC 侧）
func (l *SyncRunLogic) doAsyncSync(ctx context.Context, config config.Config, checkpointID int64, objectType string, syncCols []string, syncScope string, localCheckpointID int64) {
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
	syncService := gs.NewSyncServiceImpl(ctx, config, l.svcCtx.DAOManager)

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
	updates := map[string]interface{}{
		"ID":               checkpointID,
		"SyncStatus":       int64(0),
		"LastErrorMessage": errMsg,
		"UpdatedAt":        time.Now(),
	}
	if _, err := l.svcCtx.DAOManager.GameSyncCheckpoint.UpdateGameSyncCheckpoint(context.Background(), updates); err != nil {
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
