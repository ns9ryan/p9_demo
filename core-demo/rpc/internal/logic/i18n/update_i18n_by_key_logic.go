package i18n

import (
	"context"

	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/internal/svc"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateI18nByKeyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateI18nByKeyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateI18nByKeyLogic {
	return &UpdateI18nByKeyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateI18nByKeyLogic) UpdateI18NByKey(in *core.UpdateI18NByKeyReq) (*core.Empty, error) {
	err := l.svcCtx.Deps.UpdateI18nByKey(l.ctx, service.UpdateI18nByKeyReq{
		TransKey: in.GetTransKey(), Data: in.GetData(),
	})
	if err != nil {
		return nil, xerr.RpcErr(err)
	}
	return &core.Empty{}, nil
}
