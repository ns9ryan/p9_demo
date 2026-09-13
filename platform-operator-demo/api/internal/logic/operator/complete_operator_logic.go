// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CompleteOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteOperatorLogic {
	return &CompleteOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CompleteOperator 完成分站创建
func (l *CompleteOperatorLogic) CompleteOperator(req *types.CompleteOperatorRequest) (resp *types.CompleteOperatorResponse, err error) {
	// 完成分站创建
	_, err = l.svcCtx.OperatorRpc.Complete(
		l.ctx,
		&operatorpb.CompleteOperatorRequest{
			Id: req.Id, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回完成结果
	return &types.CompleteOperatorResponse{}, nil
}
