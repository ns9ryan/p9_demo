package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetI18nListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetI18nListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetI18nListLogic {
	return &GetI18nListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetI18nListLogic) GetI18nList(req *types.I18nListReq) (resp *types.I18nListResp, err error) {
	out, err := l.svcCtx.Core.GetI18NList(l.ctx, convert.I18nListReq(req))
	if err != nil {
		return nil, err
	}
	return convert.I18nList(out), nil
}
