package basicresourceallocationservicelogic

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operator"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatoragentlineallocation"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorlanguageallocation"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorregionallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/basicresourceallocation"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

// allocationCount 分站资源分配数量
type allocationCount struct {
	OperatorID int64 `json:"operator_id"` // 分站ID
	Count      int64 `json:"count"`       // 分配数量
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// List 获取基础资源分配列表
func (l *ListLogic) List(in *basicresourceallocation.ListBasicResourceAllocationsRequest) (*basicresourceallocation.ListBasicResourceAllocationsResponse, error) {
	// 校验分页参数
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 创建分站查询
	query := l.svcCtx.DB.Operator.Query()

	// 按关键字筛选
	if in.Keyword != nil {
		keyword := strings.TrimSpace(*in.Keyword)
		if keyword != "" {
			query = query.Where(
				operator.Or(
					operator.CodeContainsFold(keyword), // 匹配分站编码
					operator.NameContainsFold(keyword), // 匹配分站名称
				),
			)
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

	// 获取当前页分站数据
	operators, err := query.
		Order(
			operator.ByCreatedAt(sql.OrderDesc()), // 按创建时间倒序
			operator.ByID(sql.OrderDesc()),        // 创建时间相同时按ID倒序
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 当前页没有数据时直接返回
	if len(operators) == 0 {
		return &basicresourceallocation.ListBasicResourceAllocationsResponse{
			Total: int64(total),                                             // 数据总数
			List:  []*basicresourceallocation.BasicResourceAllocationInfo{}, // 基础资源分配列表
		}, nil
	}

	// 整理当前页分站ID
	operatorIDs := make([]int64, 0, len(operators))
	for _, item := range operators {
		operatorIDs = append(operatorIDs, item.ID)
	}

	var languageCountMap map[int64]int64  // 语言分配数量
	var regionCountMap map[int64]int64    // 经营地区分配数量
	var agentLineCountMap map[int64]int64 // 代理子线路分配数量

	// 语言分配数量
	{
		// 查询语言分配数量
		var counts []allocationCount
		err = l.svcCtx.DB.OperatorLanguageAllocation.
			Query().
			Where(operatorlanguageallocation.OperatorIDIn(operatorIDs...)).
			GroupBy(operatorlanguageallocation.FieldOperatorID).
			Aggregate(ent.As(ent.Count(), "count")).
			Scan(l.ctx, &counts)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}

		// 整理语言分配数量
		languageCountMap = make(map[int64]int64, len(counts))
		for _, item := range counts {
			languageCountMap[item.OperatorID] = item.Count
		}
	}

	// 经营地区分配数量
	{
		// 查询经营地区分配数量
		var counts []allocationCount
		err = l.svcCtx.DB.OperatorRegionAllocation.
			Query().
			Where(operatorregionallocation.OperatorIDIn(operatorIDs...)).
			GroupBy(operatorregionallocation.FieldOperatorID).
			Aggregate(ent.As(ent.Count(), "count")).
			Scan(l.ctx, &counts)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}

		// 整理经营地区分配数量
		regionCountMap = make(map[int64]int64, len(counts))
		for _, item := range counts {
			regionCountMap[item.OperatorID] = item.Count
		}
	}

	// 代理子线路分配数量
	{
		// 查询代理子线路分配数量
		var counts []allocationCount
		err = l.svcCtx.DB.OperatorAgentLineAllocation.
			Query().
			Where(operatoragentlineallocation.OperatorIDIn(operatorIDs...)).
			GroupBy(operatoragentlineallocation.FieldOperatorID).
			Aggregate(ent.As(ent.Count(), "count")).
			Scan(l.ctx, &counts)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}

		// 整理代理子线路分配数量
		agentLineCountMap = make(map[int64]int64, len(counts))
		for _, item := range counts {
			agentLineCountMap[item.OperatorID] = item.Count
		}
	}

	// 组装基础资源分配列表
	list := make([]*basicresourceallocation.BasicResourceAllocationInfo, 0, len(operators))
	for _, item := range operators {
		list = append(list, &basicresourceallocation.BasicResourceAllocationInfo{
			OperatorId:     item.ID,                    // 分站ID
			OperatorCode:   item.Code,                  // 分站业务编码
			OperatorName:   item.Name,                  // 分站名称
			LanguageCount:  languageCountMap[item.ID],  // 已分配语言数量
			RegionCount:    regionCountMap[item.ID],    // 已分配经营地区数量
			AgentLineCount: agentLineCountMap[item.ID], // 已分配代理子线路数量
		})
	}

	// 返回基础资源分配列表
	return &basicresourceallocation.ListBasicResourceAllocationsResponse{
		Total: int64(total), // 数据总数
		List:  list,         // 基础资源分配列表
	}, nil
}
