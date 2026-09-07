package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetI18nDictLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetI18nDictLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetI18nDictLogic {
	return &GetI18nDictLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetI18nDictLogic) GetI18NDict(in *core.GetI18NDictReq) (*core.I18NDictResp, error) {
	items, err := l.svcCtx.Deps.GetI18nDict(l.ctx, in.GetI18NGroup(), in.GetLang())
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	if items == nil {
		items = map[string]string{}
	}
	return &core.I18NDictResp{Items: items}, nil
}
