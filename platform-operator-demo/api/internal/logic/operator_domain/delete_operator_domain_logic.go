// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_domain

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/domainpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteOperatorDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteOperatorDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteOperatorDomainLogic {
	return &DeleteOperatorDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// DeleteOperatorDomain 删除分站域名
func (l *DeleteOperatorDomainLogic) DeleteOperatorDomain(req *types.DeleteOperatorDomainRequest) (resp *types.DeleteOperatorDomainResponse, err error) {
	// 删除分站域名
	_, err = l.svcCtx.OperatorDomainRpc.Delete(
		l.ctx,
		&domainpb.DeleteDomainRequest{
			Id: req.Id, // 域名ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回删除结果
	return &types.DeleteOperatorDomainResponse{}, nil
}
