package i18n

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	corei18n "oa.98ent.com/p9/core/common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateI18nLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateI18nLogic {
	return &CreateI18nLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateI18nLogic) CreateI18n(req *types.CreateI18nReq) (resp *types.I18nInfo, err error) {
	out, err := l.svcCtx.Core.CreateI18N(l.ctx, convert.CreateI18nReq(req))
	if err != nil {
		return nil, err
	}
	corei18n.Invalidate(req.I18nCode, "", req.Lang)
	corei18n.Invalidate(req.I18nCode, req.I18nGroup, req.Lang)
	return convert.I18nInfo(out), nil
}
