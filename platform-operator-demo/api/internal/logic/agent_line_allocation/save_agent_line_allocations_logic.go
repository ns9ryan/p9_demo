// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package agent_line_allocation

import (
	"context"

	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/agentlineallocationpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type SaveAgentLineAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveAgentLineAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveAgentLineAllocationsLogic {
	return &SaveAgentLineAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveAgentLineAllocations 保存代理子线路分配
func (l *SaveAgentLineAllocationsLogic) SaveAgentLineAllocations(req *types.SaveAgentLineAllocationsRequest) (resp *types.SaveAgentLineAllocationsResponse, err error) {
	// 保存分站代理子线路分配
	_, err = l.svcCtx.AgentLineAllocationRpc.Save(
		l.ctx,
		&agentlineallocationpb.SaveAgentLineAllocationsRequest{
			OperatorId:     req.OperatorId,     // 分站ID
			AgentLineCodes: req.AgentLineCodes, // 当前分配的代理子线路编码
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回保存结果
	return &types.SaveAgentLineAllocationsResponse{}, nil
}
