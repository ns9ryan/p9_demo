package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReorderI18nLangLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReorderI18nLangLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReorderI18nLangLogic {
	return &ReorderI18nLangLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReorderI18nLangLogic) ReorderI18NLang(in *core.ReorderI18NLangReq) (*core.Empty, error) {
	if err := l.svcCtx.Deps.ReorderI18nLang(l.ctx, in.GetId(), in.GetTargetId()); err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
