// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package language_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveLanguageAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveLanguageAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLanguageAllocationsLogic {
	return &SaveLanguageAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveLanguageAllocationsLogic) SaveLanguageAllocations(req *types.SaveLanguageAllocationsRequest) (resp *types.SaveLanguageAllocationsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
