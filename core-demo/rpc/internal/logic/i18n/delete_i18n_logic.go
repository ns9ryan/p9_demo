package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteI18nLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteI18nLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteI18nLogic {
	return &DeleteI18nLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteI18nLogic) DeleteI18N(in *core.IDsReq) (*core.Empty, error) {
	if err := l.svcCtx.Deps.DeleteI18ns(l.ctx, utils.ResolveIDs(in.Id, in.Ids)); err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
