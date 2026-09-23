package i18n

import (
	"context"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteI18nByKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteI18nByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteI18nByKeyLogic {
	return &DeleteI18nByKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteI18nByKeyLogic) DeleteI18NByKey(in *core.DeleteI18NByKeyReq) (*core.Empty, error) {
	if err := l.svcCtx.Deps.DeleteI18nsByKey(l.ctx, in.GetI18NCode(), in.GetI18NGroup(), in.GetTransKey()); err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
