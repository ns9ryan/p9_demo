package agentlineallocationservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/agentlineallocation"
)

// toAgentLineAllocationInfo 转换代理子线路分配信息
func toAgentLineAllocationInfo(data *ent.OperatorAgentLineAllocation) *agentlineallocation.AgentLineAllocationInfo {
	return &agentlineallocation.AgentLineAllocationInfo{
		Id:            data.ID,                    // 分配记录ID
		OperatorId:    data.OperatorID,            // 分站ID
		AgentLineCode: data.AgentLineCode,         // 代理子线路编码
		CreatedAt:     data.CreatedAt.UnixMilli(), // 分配时间
	}
}
