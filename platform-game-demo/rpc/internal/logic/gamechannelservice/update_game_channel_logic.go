package gamechannelservicelogic

import (
	"context"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/rpc/ent"
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

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &platformgame.UpdateGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	channel := &ent.GameChannel{}
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(channel).Error; err != nil {
		l.Errorf("[RPC UpdateGameChannel] failed: %v", err)
		return &platformgame.UpdateGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	updateData := make(map[string]interface{})
	updateData["updated_at"] = time.Now()

	if in.NameI18N != "" {
		updateData["name_i18n"] = in.NameI18N
	}
	if in.SortNo > 0 {
		updateData["sort_no"] = in.SortNo
	}
	if in.Status > 0 {
		updateData["status"] = in.Status
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Model(channel).Updates(updateData).Error; err != nil {
		return &platformgame.UpdateGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := channel.Status

	l.svcCtx.DB.WithContext(l.ctx).Where("id = ?", in.Id).First(channel)

	// 如果需要强制踢线且状态变更为停用，发送 kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameChannel] sending force-quit event for channel id=%d", channel.SourceId)
			// 查出所有游戏渠道下的游戏，并发送强制踢线事件
			gameIds, err := getGamesByChannelID(context.Background(), l.svcCtx.DB, channel.SourceId)
			if err != nil {
				l.Errorf("[RPC UpdateGameChannel] failed to get games by channel id=%d: %v", channel.SourceId, err)
				return
			}
			forceQuit := utils.ForceQuitEvent{
				Basis:   utils.KafkaPayloadBasis{Stage: 0, Producer: "platform-game-rpc"},
				Scope:   constant.ScopeGameChannel,
				ChanID:  channel.SourceId,
				GameIDs: gameIds,
				QuitAt:  time.Now(),
			}
			l.Infof("[RPC UpdateGameChannel] sending force-quit event json: %s", utils.JSON(forceQuit))
			err = utils.SendForceQuitEvent(context.Background(), l.svcCtx.Config.Kafka.Brokers, &forceQuit)
			if err != nil {
				l.Errorf("[RPC UpdateGameChannel] failed to send force-quit event for channel id=%d: %v", channel.SourceId, err)
			}
		}()
	}

	return &platformgame.UpdateGameChannelResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    logic.ChannelModelToProto(channel),
	}, nil
}

func getGamesByChannelID(ctx context.Context, db *gorm.DB, channelID int64) ([]int64, error) {
	var games []ent.Game
	err := db.WithContext(ctx).
		Where("channel_id = ? AND deleted_at IS NULL", channelID).
		Find(&games).Error
	if err != nil {
		return nil, err
	}
	ids := make([]int64, len(games))
	for i, game := range games {
		ids[i] = game.SourceId.Int64
	}
	return ids, nil
}
