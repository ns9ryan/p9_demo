package logic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/kafka"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
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
func (l *UpdateGameChannelLogic) UpdateGameChannel(in *platformgame.UpdateGameChannelRequest) (*platformgame.GetGameChannelResp, error) {
	l.Infof("[RPC UpdateGameChannel] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		return &platformgame.GetGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	var channel model.Channel
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&channel).Error; err != nil {
		l.Errorf("[RPC UpdateGameChannel] failed: %v", err)
		return &platformgame.GetGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	updateData := make(map[string]interface{})
	updateData["updated_at"] = time.Now()

	if in.NameI18N != "" {
		updateData["name_i18n"] = utils.MustParseJSON([]byte(in.NameI18N))
	}
	if in.SortNo > 0 {
		updateData["sort_no"] = in.SortNo
	}
	if in.Status > 0 {
		updateData["status"] = in.Status
	}

	if err := l.svcCtx.DB.WithContext(l.ctx).Model(&channel).Updates(updateData).Error; err != nil {
		return &platformgame.GetGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	oldStatus := channel.Status

	l.svcCtx.DB.WithContext(l.ctx).Where("id = ?", in.Id).First(&channel)

	// 如果需要强制踢线且状态变更为停用，发�?kafka 事件
	if in.ForceLogout && in.Status == 2 && oldStatus != 2 {
		go func() {
			l.Infof("[RPC UpdateGameChannel] sending force-quit event for channel id=%d", channel.ID)
			if err := kafka.SendForceQuitEvent(l.ctx, channel.ID, "game_channel"); err != nil {
				l.Infof("[RPC UpdateGameChannel] send force-quit event failed: %v", err)
			}
		}()
	}

	return &platformgame.GetGameChannelResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data: &platformgame.GameChannelResp{
			Id:          channel.ID,
			SourceId:    channel.SourceID,
			ChannelCode: channel.ChannelCode,
			Status:      int32(channel.Status),
		},
	}, nil
}
