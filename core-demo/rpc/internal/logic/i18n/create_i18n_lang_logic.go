package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/logic"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateI18nLangLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateI18nLangLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateI18nLangLogic {
	return &CreateI18nLangLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateI18nLangLogic) CreateI18NLang(in *core.CreateI18NLangReq) (*core.I18NLangInfo, error) {
	row, err := l.svcCtx.Deps.CreateI18nLang(l.ctx, service.CreateI18nLangReq{
		Lang: in.GetLang(), Name: in.GetName(), Disabled: int16(in.GetDisabled()), IsDefault: int16(in.GetIsDefault()),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToI18nLangInfo(*row), nil
}
