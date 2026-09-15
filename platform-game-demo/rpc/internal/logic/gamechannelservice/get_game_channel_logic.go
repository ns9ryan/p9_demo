package gamechannelservicelogic

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
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

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Error("[RPC GetGameChannel] DAO Manager not available")
		return &platformgame.GetGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	channel, err := l.svcCtx.DAOManager.GameChannel.GetGameChannelByID(l.ctx, in.Id)
	if err != nil {
		l.Errorf("[RPC GetGameChannel] query failed: %v", err)
		return &platformgame.GetGameChannelResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}

	resp := &platformgame.GetGameChannelResp{
		Code:    constant.CodeSuccess,
		Message: "success",
		Data:    logic.ChannelModelToProto(channel),
	}

	l.Infof("[RPC GetGameChannel] success: id=%d", channel.ID)
	return resp, nil
}
