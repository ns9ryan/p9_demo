// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOperatorLogic {
	return &DeleteOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteOperator 删除分站
func (l *DeleteOperatorLogic) DeleteOperator(req *types.DeleteOperatorRequest) (resp *types.DeleteOperatorResponse, err error) {
	// 获取当前分站
	current, err := l.svcCtx.OperatorRpc.Get(
		l.ctx,
		&operatorpb.GetOperatorRequest{
			Id: req.Id, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 仅未发布或发布失败的分站允许删除
	if current.Operator.PublishStatus != 1 && current.Operator.PublishStatus != 4 {
		return nil, grpcerror.InvalidArgument(i18nkey.ConstraintError)
	}

	// 按分站清理管理员账号，统一放 OperatorRpc.Delete 中处理
	// _, err = l.svcCtx.OperatorAdminRpc.DeleteByOperatorId(
	// 	l.ctx,
	// 	&adminpb.DeleteAdminsByOperatorIdRequest{
	// 		OperatorId: req.Id, // 分站ID
	// 	},
	// )
	// if err != nil {
	// 	return nil, err
	// }

	// TODO platform-game提供按分站清理游戏资源分配的RPC后在这里接入

	// 删除分站
	_, err = l.svcCtx.OperatorRpc.Delete(
		l.ctx,
		&operatorpb.DeleteOperatorRequest{
			Id: req.Id, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回删除结果
	return &types.DeleteOperatorResponse{}, nil
}
