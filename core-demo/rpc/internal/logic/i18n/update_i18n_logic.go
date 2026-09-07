package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateI18nLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateI18nLogic {
	return &UpdateI18nLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateI18nLogic) UpdateI18N(in *core.UpdateI18NReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.UpdateI18n(l.ctx, service.UpdateI18nReq{
		ID: in.Id, I18nGroup: in.I18NGroup, TransKey: in.TransKey, Lang: in.Lang, Value: in.Value,
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
