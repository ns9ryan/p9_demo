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

type CreateI18nLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateI18nLogic {
	return &CreateI18nLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateI18nLogic) CreateI18N(in *core.CreateI18NReq) (*core.I18NInfo, error) {
	row, err := l.svcCtx.Deps.CreateI18n(l.ctx, service.CreateI18nReq{
		I18nGroup: in.GetI18NGroup(), TransKey: in.GetTransKey(), Lang: in.GetLang(), Value: in.GetValue(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return logic.ToI18nInfo(*row), nil
}
