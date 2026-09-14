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

type ListOperatorsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOperatorsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOperatorsLogic {
	return &ListOperatorsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListOperators 获取分站管理列表
func (l *ListOperatorsLogic) ListOperators(req *types.ListOperatorsRequest) (resp *types.ListOperatorsResponse, err error) {
	// 获取分站管理列表
	result, err := l.svcCtx.OperatorRpc.List(
		l.ctx,
		&operatorpb.ListOperatorsRequest{
			Page:           req.Page,           // 页码, 从1开始
			PageSize:       req.PageSize,       // 每页数量
			Keyword:        req.Keyword,        // 搜索关键字, 匹配分站编码或名称
			CreationStatus: req.CreationStatus, // 创建状态: 1草稿, 2已完成
			PublishStatus:  req.PublishStatus,  // 发布状态: 1未发布, 2发布中, 3已发布, 4发布失败
			Status:         req.Status,         // 分站状态: 1正常, 2暂停, 3关闭
		},
	)
	if err != nil {
		return nil, err
	}

	// 转换分站列表
	list := make([]types.OperatorInfo, 0, len(result.List))
	for _, item := range result.List {
		list = append(list, toOperatorInfo(item))
	}

	// 返回分站管理列表
	return &types.ListOperatorsResponse{
		Total: result.Total, // 数据总数
		List:  list,         // 分站列表
	}, nil
}
