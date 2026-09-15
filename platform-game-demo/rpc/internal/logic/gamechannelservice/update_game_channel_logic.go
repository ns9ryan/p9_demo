package gamechannelservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGameChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameChannelLogic {
	return &UpdateGameChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏渠道
func (l *UpdateGameChannelLogic) UpdateGameChannel(in *platformgame.UpdateGameChannelRequest) (*platformgame.UpdateGameChannelResp, error) {
	l.Infof("[RPC UpdateGameChannel] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		return &platformgame.UpdateGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	// 查询当前渠道数据
	channelRecord, err := l.svcCtx.DAOManager.GameChannel.GetGameChannelByID(l.ctx, in.Id)
	if err != nil {
		l.Errorf("[RPC UpdateGameChannel] failed to get channel: %v", err)
		return &platformgame.UpdateGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := channelRecord.Status

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
	updatedChannel, err := l.svcCtx.DAOManager.GameChannel.UpdateGameChannel(l.ctx, in.Id, updates)
	if err != nil {
		return &platformgame.UpdateGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameChannel] sending force-quit event for channel id=%d", channelRecord.SourceID)
			// 查出所有游戏渠道下的游戏，并发送强制踢线事件
			gameIds, err := getGamesByChannelID(context.Background(), l.svcCtx, channelRecord.SourceID)
			if err != nil {
				l.Errorf("[RPC UpdateGameChannel] failed to get games by channel id=%d: %v", channelRecord.SourceID, err)
				return
			}
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGameChannel,
				ChanID:  channelRecord.SourceID,
				GameIDs: gameIds,
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGameChannel] sending force-quit event json: %s", utils.JSON(forceQuit))
			err = utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGameChannel] failed to send force-quit event for channel id=%d: %v", channelRecord.SourceID, err)
			}
		}()
	}

	return &platformgame.UpdateGameChannelResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.ChannelModelToProto(updatedChannel),
	}, nil
}

func getGamesByChannelID(ctx context.Context, svcCtx *svc.ServiceContext, channelID int64) ([]int64, error) {
	if svcCtx == nil || svcCtx.DAOManager == nil {
		return nil, nil
	}

	games, err := svcCtx.DAOManager.Game.GetAllGame(
		ctx,
		0,         // providerId
		0,         // categoryId
		channelID, // channelId
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
