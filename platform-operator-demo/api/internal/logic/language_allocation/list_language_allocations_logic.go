// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package language_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListLanguageAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListLanguageAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLanguageAllocationsLogic {
	return &ListLanguageAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListLanguageAllocationsLogic) ListLanguageAllocations(req *types.ListLanguageAllocationsRequest) (resp *types.ListLanguageAllocationsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
