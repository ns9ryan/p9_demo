package logic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetGameChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameChannelLogic {
	return &GetGameChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取单个游戏渠道
func (l *GetGameChannelLogic) GetGameChannel(in *platformgame.GetGameChannelRequest) (*platformgame.GetGameChannelResp, error) {
	l.Infof("[RPC GetGameChannel] received req: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Error("[RPC GetGameChannel] database not available")
		return &platformgame.GetGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: "database not available",
		}, nil
	}

	var channel model.Channel
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Where("id = ? AND deleted_at IS NULL", in.Id).
		First(&channel).Error; err != nil {
		l.Errorf("[RPC GetGameChannel] query failed: %v", err)
		return &platformgame.GetGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	// 查询该渠道下的游戏数
	var gameCount int64
	if err := l.svcCtx.DB.WithContext(l.ctx).
		Model(&model.Game{}).
		Where("channel_id = ? AND deleted_at IS NULL", channel.ID).
		Count(&gameCount).Error; err != nil {
		l.Infof("[RPC GetGameChannel] failed to count games for channel %d: %v", channel.ID, err)
		gameCount = 0
	}

	var sourceNameI18n string
	if channel.SourceNameI18n != nil {
		sourceNameI18n = string(utils.MustMarshalJSON(channel.SourceNameI18n))
	}

	resp := &platformgame.GetGameChannelResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data: &platformgame.GameChannelResp{
			Id:             channel.ID,
			SourceId:       channel.SourceID,
			ChannelCode:    channel.ChannelCode,
			SourceNameI18N: sourceNameI18n,
			Status:         int32(channel.Status),
			GameCount:      gameCount,
			CreatedAt:      channel.CreatedAt.Unix(),
			UpdatedAt:      channel.UpdatedAt.Unix(),
		},
	}

	l.Infof("[RPC GetGameChannel] success: id=%d, game_count=%d", channel.ID, gameCount)
	return resp, nil
}
