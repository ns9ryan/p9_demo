package agentlineallocationservicelogic

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatoragentlineallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/agentlineallocation"
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

// List 获取代理子线路分配列表
func (l *ListLogic) List(in *agentlineallocation.ListAgentLineAllocationsRequest) (*agentlineallocation.ListAgentLineAllocationsResponse, error) {
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

	// 获取代理子线路分配
	results, err := l.svcCtx.DB.OperatorAgentLineAllocation.
		Query().
		Where(operatoragentlineallocation.OperatorIDEQ(in.OperatorId)).
		Order(
			operatoragentlineallocation.ByCreatedAt(sql.OrderDesc()), // 按分配时间倒序
			operatoragentlineallocation.ByID(sql.OrderDesc()),        // 分配时间相同时按ID倒序
		).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换代理子线路分配列表
	list := make([]*agentlineallocation.AgentLineAllocationInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toAgentLineAllocationInfo(result))
	}

	// 返回代理子线路分配列表
	return &agentlineallocation.ListAgentLineAllocationsResponse{
		List: list, // 代理子线路分配列表
	}, nil
}
