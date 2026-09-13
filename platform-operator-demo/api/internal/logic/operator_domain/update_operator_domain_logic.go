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

type UpdateOperatorDomainLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOperatorDomainLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOperatorDomainLogic {
	return &UpdateOperatorDomainLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateOperatorDomain 修改分站域名
func (l *UpdateOperatorDomainLogic) UpdateOperatorDomain(req *types.UpdateOperatorDomainRequest) (resp *types.UpdateOperatorDomainResponse, err error) {
	// 修改分站域名
	_, err = l.svcCtx.OperatorDomainRpc.Update(
		l.ctx,
		&domainpb.UpdateDomainRequest{
			Id:         req.Id,         // 域名ID
			DomainName: req.DomainName, // 域名, 不包含协议和端口
			DomainType: req.DomainType, // 域名类型: 1分站后台, 2代理后台, 3会员H5
			Status:     req.Status,     // 域名状态: 1启用, 2停用
			Remark:     req.Remark,     // 总网内部备注
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回修改结果
	return &types.UpdateOperatorDomainResponse{}, nil
}
