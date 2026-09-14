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

type GetOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorLogic {
	return &GetOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetOperator 获取分站
func (l *GetOperatorLogic) GetOperator(req *types.GetOperatorRequest) (resp *types.GetOperatorResponse, err error) {
	// 获取分站
	result, err := l.svcCtx.OperatorRpc.Get(
		l.ctx,
		&operatorpb.GetOperatorRequest{
			Id: req.Id, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回分站信息
	return &types.GetOperatorResponse{
		OperatorInfo: toOperatorInfo(result.Operator),
	}, nil
}
