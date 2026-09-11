// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package domain

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteDomainLogic {
	return &DeleteDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteDomainLogic) DeleteDomain(req *types.DeleteDomainRequest) (resp *types.DeleteDomainResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
