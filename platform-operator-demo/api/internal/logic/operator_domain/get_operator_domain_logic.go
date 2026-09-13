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

type GetOperatorDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOperatorDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorDomainLogic {
	return &GetOperatorDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// GetOperatorDomain 获取分站域名
func (l *GetOperatorDomainLogic) GetOperatorDomain(req *types.GetOperatorDomainRequest) (resp *types.GetOperatorDomainResponse, err error) {
	// 获取分站域名
	result, err := l.svcCtx.OperatorDomainRpc.Get(
		l.ctx,
		&domainpb.GetDomainRequest{
			Id: req.Id, // 域名ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回分站域名信息
	return &types.GetOperatorDomainResponse{
		OperatorDomainInfo: toOperatorDomainInfo(result.Domain), // 分站域名信息
	}, nil
}
