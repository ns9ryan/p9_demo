package regionallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorregionallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/regionallocationpb"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// List 获取经营地区分配列表
func (l *ListLogic) List(in *regionallocationpb.ListRegionAllocationsRequest) (*regionallocationpb.ListRegionAllocationsResponse, error) {
	// 分站ID必须大于0
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 确认分站存在
	_, err := l.svcCtx.DB.Operator.Get(l.ctx, in.OperatorId)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 获取分站当前全部经营地区分配
	results, err := l.svcCtx.DB.OperatorRegionAllocation.
		Query().
		Where(operatorregionallocation.OperatorIDEQ(in.OperatorId)).
		Order(
			operatorregionallocation.ByCreatedAt(sql.OrderDesc()), // 按分配时间倒序
			operatorregionallocation.ByID(sql.OrderDesc()),        // 分配时间相同时按ID倒序
		).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换经营地区分配列表
	list := make([]*regionallocationpb.RegionAllocationInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toRegionAllocationInfo(result))
	}

	// 返回经营地区分配列表
	return &regionallocationpb.ListRegionAllocationsResponse{
		List: list, // 经营地区分配列表
	}, nil
}
