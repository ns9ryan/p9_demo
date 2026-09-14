package languageallocationservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorlanguageallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/languageallocationpb"

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

// List 获取语言分配列表
func (l *ListLogic) List(in *languageallocationpb.ListLanguageAllocationsRequest) (*languageallocationpb.ListLanguageAllocationsResponse, error) {
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

	// 获取分站当前全部语言分配
	results, err := l.svcCtx.DB.OperatorLanguageAllocation.
		Query().
		Where(operatorlanguageallocation.OperatorIDEQ(in.OperatorId)).
		Order(
			operatorlanguageallocation.ByCreatedAt(sql.OrderDesc()), // 按分配时间倒序
			operatorlanguageallocation.ByID(sql.OrderDesc()),        // 分配时间相同时按ID倒序
		).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换语言分配列表
	list := make([]*languageallocationpb.LanguageAllocationInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toLanguageAllocationInfo(result))
	}

	// 返回语言分配列表
	return &languageallocationpb.ListLanguageAllocationsResponse{
		List: list, // 语言分配列表
	}, nil
}
