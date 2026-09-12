package languageallocationservicelogic

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorlanguageallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/languageallocation"
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

// List 获取语言分配列表
func (l *ListLogic) List(in *languageallocation.ListLanguageAllocationsRequest) (*languageallocation.ListLanguageAllocationsResponse, error) {
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

	// 创建语言分配查询
	query := l.svcCtx.DB.OperatorLanguageAllocation.
		Query().
		Where(operatorlanguageallocation.OperatorIDEQ(in.OperatorId))

	// 按语言编码筛选
	if in.LanguageCode != nil {
		languageCode := strings.TrimSpace(*in.LanguageCode)
		if languageCode != "" {
			query = query.Where(operatorlanguageallocation.LanguageCodeEQ(languageCode))
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

	// 获取当前页语言分配数据
	results, err := query.
		Order(
			operatorlanguageallocation.ByCreatedAt(sql.OrderDesc()), // 按分配时间倒序
			operatorlanguageallocation.ByID(sql.OrderDesc()),        // 分配时间相同时按ID倒序
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换语言分配列表
	list := make([]*languageallocation.LanguageAllocationInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toLanguageAllocationInfo(result))
	}

	// 返回语言分配列表
	return &languageallocation.ListLanguageAllocationsResponse{
		Total: int64(total), // 数据总数
		List:  list,         // 语言分配列表
	}, nil
}
