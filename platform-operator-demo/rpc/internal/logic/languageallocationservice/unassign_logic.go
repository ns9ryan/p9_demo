package languageallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/languageallocation"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnassignLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnassignLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnassignLogic {
	return &UnassignLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 取消语言分配
func (l *UnassignLogic) Unassign(in *languageallocation.UnassignLanguageRequest) (*languageallocation.UnassignLanguageResponse, error) {
	// todo: add your logic here and delete this line

	return &languageallocation.UnassignLanguageResponse{}, nil
}
