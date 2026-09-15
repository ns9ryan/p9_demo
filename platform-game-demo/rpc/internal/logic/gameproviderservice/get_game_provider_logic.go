package gameproviderservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type GetGameProviderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameProviderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameProviderLogic {
	return &GetGameProviderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取单个游戏供应商
func (l *GetGameProviderLogic) GetGameProvider(in *platformgame.GetGameProviderRequest) (*platformgame.GetGameProviderResp, error) {
	l.Infof("[RPC GetGameProvider] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetGameProvider] DAO Manager not available")
		return &platformgame.GetGameProviderResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
			Data:    nil,
		}, nil
	}

	provider, err := l.svcCtx.DAOManager.GameProvider.GetGameProviderByID(l.ctx, in.Id)
	if err != nil {
		l.Infof("[RPC GetGameProvider] provider not found: id=%d", in.Id)
		return &platformgame.GetGameProviderResp{
			Code:    constant.CodeNotFound,
			Message: "provider not found",
			Data:    nil,
		}, nil
	}

	item := logic.ProviderModelToProto(provider)

	return &platformgame.GetGameProviderResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    item,
	}, nil
}
