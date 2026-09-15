package gameservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGameLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameLogic {
	return &UpdateGameLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏
func (l *UpdateGameLogic) UpdateGame(in *platformgame.UpdateGameRequest) (*platformgame.UpdateGameResp, error) {
	l.Infof("[RPC UpdateGame] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		return &platformgame.UpdateGameResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	// 查询当前游戏数据
	gameRecord, err := l.svcCtx.DAOManager.Game.GetGameByID(l.ctx, in.Id)
	if err != nil {
		l.Errorf("[RPC UpdateGame] failed to get game: %v", err)
		return &platformgame.UpdateGameResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := gameRecord.Status

	// 准备更新数据
	updates := make(map[string]interface{})
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.SortNo > 0 {
		updates["sort_no"] = int64(in.SortNo)
	}
	if in.Status > 0 {
		updates["status"] = int64(in.Status)
	}
	if in.ImageUrl != "" {
		updates["image_url"] = in.ImageUrl
	}
	if in.SupportsEmbed != 0 {
		updates["supports_embed"] = in.SupportsEmbed == 1
	}
	if in.SupportsRedirect != 0 {
		updates["supports_redirect"] = in.SupportsRedirect == 1
	}
	updates["updated_at"] = time.Now()

	// 执行更新
	_, err = l.svcCtx.DAOManager.Game.UpdateGame(l.ctx, in.Id, updates)
	if err != nil {
		return &platformgame.UpdateGameResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGame] sending force-quit event for game id=%d", gameRecord.SourceID)
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGame,
				GameIDs: []int64{gameRecord.SourceID},
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGame] sending force-quit event json: %s", utils.JSON(forceQuit))
			err := utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGame] failed to send force-quit event for game id=%d: %v", gameRecord.SourceID, err)
			}
		}()
	}

	return &platformgame.UpdateGameResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
	}, nil
}
