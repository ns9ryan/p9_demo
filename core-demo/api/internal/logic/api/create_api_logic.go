package api

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/convert"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/i18n"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateApiLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateApiLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateApiLogic {
	return &CreateApiLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateApiLogic) CreateApi(req *types.CreateApiReq) (resp *types.ApiInfo, err error) {
	out, err := l.svcCtx.Core.CreateApi(l.ctx, convert.CreateApiReq(req))
	if err != nil {
		return nil, err
	}
	return convert.ApiInfo(l.ctx, i18n.CodeByPartnerMode(l.svcCtx.Config.PartnerMode), out), nil
}
