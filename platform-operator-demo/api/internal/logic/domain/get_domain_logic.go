// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package domain

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetDomainLogic {
	return &GetDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetDomainLogic) GetDomain(req *types.GetDomainRequest) (resp *types.GetDomainResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
