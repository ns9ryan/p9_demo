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

type UpdateI18nLangLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateI18nLangLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateI18nLangLogic {
	return &UpdateI18nLangLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateI18nLangLogic) UpdateI18NLang(in *core.UpdateI18NLangReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.UpdateI18nLang(l.ctx, service.UpdateI18nLangReq{
		ID:        in.GetId(),
		Lang:      in.Lang,
		Name:      in.Name,
		Disabled:  logic.ToInt16Ptr(in.Disabled),
		IsDefault: logic.ToInt16Ptr(in.IsDefault),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
