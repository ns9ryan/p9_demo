// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package agent_line_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/pkg/agentline"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/agentlineallocationpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAgentLineAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListAgentLineAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAgentLineAllocationsLogic {
	return &ListAgentLineAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ListAgentLineAllocations 获取代理子线路分配列表
func (l *ListAgentLineAllocationsLogic) ListAgentLineAllocations(req *types.ListAgentLineAllocationsRequest) (resp *types.ListAgentLineAllocationsResponse, err error) {
	// 获取分站当前代理子线路分配关系
	result, err := l.svcCtx.AgentLineAllocationRpc.List(
		l.ctx,
		&agentlineallocationpb.ListAgentLineAllocationsRequest{
			OperatorId: req.OperatorId, // 分站ID
		},
	)
	if err != nil {
		return nil, err
	}

	// 记录已分配代理子线路和分配时间
	allocatedAgentLines := make(map[string]*int64, len(result.List))
	for _, item := range result.List {
		allocatedAgentLines[item.AgentLineCode] = new(item.CreatedAt)
	}

	// 组合全部代理子线路和当前分配关系
	agentLineCodes := agentline.All()
	list := make([]types.AgentLineAllocationInfo, 0, len(agentLineCodes))
	for _, agentLineCode := range agentLineCodes {
		allocatedAt, allocated := allocatedAgentLines[agentLineCode]

		list = append(list, types.AgentLineAllocationInfo{
			AgentLineCode: agentLineCode, // 代理子线路编码
			Allocated:     allocated,     // 是否已分配
			AllocatedAt:   allocatedAt,   // 分配时间, 未分配时为空
		})
	}

	// 返回代理子线路分配列表
	return &types.ListAgentLineAllocationsResponse{
		List: list, // 代理子线路分配列表
	}, nil
}
