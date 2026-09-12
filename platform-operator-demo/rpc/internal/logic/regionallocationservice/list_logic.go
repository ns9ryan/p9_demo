package regionallocationservicelogic

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorregionallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/regionallocation"
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
func (l *ListLogic) List(in *regionallocation.ListRegionAllocationsRequest) (*regionallocation.ListRegionAllocationsResponse, error) {
	// 校验分页参数
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

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

	// 创建经营地区分配查询
	query := l.svcCtx.DB.OperatorRegionAllocation.
		Query().
		Where(operatorregionallocation.OperatorIDEQ(in.OperatorId))

	// 按经营地区编码筛选
	if in.RegionCode != nil {
		regionCode := strings.ToUpper(strings.TrimSpace(*in.RegionCode))
		if regionCode != "" {
			query = query.Where(operatorregionallocation.RegionCodeEQ(regionCode))
		}
	}

	// 获取符合条件的数据总数
	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 计算分页偏移量
	offset := (in.Page - 1) * in.PageSize

	// 获取当前页经营地区分配数据
	results, err := query.
		Order(
			operatorregionallocation.ByCreatedAt(sql.OrderDesc()), // 按分配时间倒序
			operatorregionallocation.ByID(sql.OrderDesc()),        // 分配时间相同时按ID倒序
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换经营地区分配列表
	list := make([]*regionallocation.RegionAllocationInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toRegionAllocationInfo(result))
	}

	// 返回经营地区分配列表
	return &regionallocation.ListRegionAllocationsResponse{
		Total: int64(total), // 数据总数
		List:  list,         // 经营地区分配列表
	}, nil
}
