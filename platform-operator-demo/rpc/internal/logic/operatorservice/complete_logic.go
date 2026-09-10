package operatorservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/operator"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteLogic {
	return &CompleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 完成分站创建
func (l *CompleteLogic) Complete(in *operator.CompleteOperatorRequest) (*operator.CompleteOperatorResponse, error) {
	// todo: add your logic here and delete this line

	return &operator.CompleteOperatorResponse{}, nil
}
