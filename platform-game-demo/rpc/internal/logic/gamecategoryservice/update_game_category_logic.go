package gamecategoryservicelogic

import (
	"context"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type UpdateGameCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameCategoryLogic {
	return &UpdateGameCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏分类
func (l *UpdateGameCategoryLogic) UpdateGameCategory(in *platform_game.UpdateGameCategoryRequest) (*platform_game.UpdateGameCategoryResp, error) {
	l.Infof("[RPC UpdateGameCategory] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC UpdateGameCategory] DAO Manager not available")
		return &platform_game.UpdateGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	// 查询当前分类数据
	categoryRecord, err := l.svcCtx.DAOManager.GameCategory.GetGameCategoryByID(l.ctx, in.Id)
	if err != nil {
		l.Errorf("[RPC UpdateGameCategory] failed to get category: %v", err)
		return &platform_game.UpdateGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := categoryRecord.Status

	// 准备更新数据
	updates := make(map[string]interface{})
	if in.SortNo > 0 {
		updates["sort_no"] = int64(in.SortNo)
	}
	if in.Status > 0 {
		updates["status"] = int64(in.Status)
	}
	updates["updated_at"] = time.Now()

	// 执行更新
	_, err = l.svcCtx.DAOManager.GameCategory.UpdateGameCategory(l.ctx, in.Id, updates)
	if err != nil {
		return &platform_game.UpdateGameCategoryResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameCategory] sending force-quit event for category id=%d", categoryRecord.SourceID)
			// 查出所有游戏分类下的游戏，并发送强制踢线事件
			gameIds, err := getGamesByCategoryID(context.Background(), l.svcCtx, categoryRecord.SourceID)
			if err != nil {
				l.Errorf("[RPC UpdateGameCategory] failed to get games by category id=%d: %v", categoryRecord.SourceID, err)
				return
			}
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGameCategory,
				CatID:   categoryRecord.SourceID,
				GameIDs: gameIds,
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGameCategory] sending force-quit event json: %s", utils.JSON(forceQuit))
			err = utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGameCategory] failed to send force-quit event for category id=%d: %v", categoryRecord.SourceID, err)
			}
		}()
	}

	return &platform_game.UpdateGameCategoryResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
	}, nil
}

func getGamesByCategoryID(ctx context.Context, svcCtx *svc.ServiceContext, categoryID int64) ([]int64, error) {
	if svcCtx == nil || svcCtx.DAOManager == nil {
		return nil, nil
	}

	games, err := svcCtx.DAOManager.Game.GetAllGame(
		ctx,
		0,          // providerId
		categoryID, // categoryId
		0,          // channelId
	)
	if err != nil {
		return nil, err
	}

	ids := make([]int64, 0, len(games))
	for _, game := range games {
		if game != nil {
			ids = append(ids, game.SourceID)
		}
	}
	return ids, nil
}
