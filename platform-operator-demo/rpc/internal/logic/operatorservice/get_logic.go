package operatorservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/operator"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLogic {
	return &GetLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站
func (l *GetLogic) Get(in *operator.GetOperatorRequest) (*operator.GetOperatorResponse, error) {
	// todo: add your logic here and delete this line

	return &operator.GetOperatorResponse{}, nil
}
