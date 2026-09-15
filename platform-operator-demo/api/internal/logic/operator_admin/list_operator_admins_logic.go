// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator_admin

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOperatorAdminsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperatorAdminsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperatorAdminsLogic {
	return &ListOperatorAdminsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListOperatorAdmins 获取分站管理员列表
func (l *ListOperatorAdminsLogic) ListOperatorAdmins(req *types.ListOperatorAdminsRequest) (resp *types.ListOperatorAdminsResponse, err error) {
	result, err := l.svcCtx.OperatorAdminRpc.List(
		l.ctx,
		&adminpb.ListAdminsRequest{
			Page:       req.Page,       // 页码, 从1开始
			PageSize:   req.PageSize,   // 每页数量
			OperatorId: req.OperatorId, // 分站ID
			Keyword:    req.Keyword,    // 搜索关键字, 匹配账号或显示名称
			Status:     req.Status,     // 状态: 1启用, 2停用
		},
	)
	if err != nil {
		return nil, err
	}

	list := make([]types.OperatorAdminInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, toOperatorAdminInfo(item))
	}

	return &types.ListOperatorAdminsResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 管理员列表
	}, nil
}
