// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package game_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveGameAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveGameAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveGameAllocationsLogic {
	return &SaveGameAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SaveGameAllocationsLogic) SaveGameAllocations(req *types.SaveGameAllocationsRequest) (resp *types.SaveGameAllocationsResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
