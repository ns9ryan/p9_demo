package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/constant"
	"oa.98ent.com/p9/platform-game/common/utils"
	"oa.98ent.com/p9/platform-game/pkg/game"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
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
func (l *GetGameProviderLogic) GetGameProvider(in *platformgame.GetGameProviderRequest) (*platformgame.GetGameProviderListResp, error) {
	l.Infof("[RPC GetGameProvider] received request: id=%d", in.Id)

	if l.svcCtx == nil {
		l.Error("[RPC GetGameProvider] ServiceContext is nil")
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "ServiceContext is nil",
			Data:    nil,
		}, nil
	}

	if l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameProvider] Database not available")
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
			Data:    nil,
		}, nil
	}

	l.Infof("[RPC GetGameProvider] �?Database is available")

	// Call pkg/game to get provider
	params := &game.GameProviderGetParams{
		ID: in.Id,
	}
	provider, err := game.GameProviderGetDB(l.ctx, l.svcCtx.DB, params)
	if err != nil {
		l.Errorf("[RPC GetGameProvider] failed to get provider: %v", err)
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeInternalError,
			Message: "failed to get provider: " + err.Error(),
			Data:    nil,
		}, nil
	}

	if provider == nil {
		l.Infof("[RPC GetGameProvider] provider not found: id=%d", in.Id)
		return &platformgame.GetGameProviderListResp{
			Code:    constant.CodeNotFound,
			Message: "provider not found",
			Data:    nil,
		}, nil
	}

	l.Infof("[RPC GetGameProvider] found provider: id=%d, code=%s", provider.ID, provider.ProviderCode)

	// Convert to proto message
	item := &platformgame.ProviderInfo{
		Id:             provider.ID,
		Code:           provider.ProviderCode,
		Name:           provider.ProviderCode,
		Status:         int32(provider.Status),
		NameI18N:       string(utils.MustMarshalJSON(provider.NameI18n)),
		SourceNameI18N: string(utils.MustMarshalJSON(provider.SourceNameI18n)),
	}

	return &platformgame.GetGameProviderListResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
		Data:    []*platformgame.ProviderInfo{item},
	}, nil
}
