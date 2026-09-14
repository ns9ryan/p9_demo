package agentlineallocationservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/agentlineallocationpb"
)

// toAgentLineAllocationInfo 转换代理子线路分配信息
func toAgentLineAllocationInfo(data *ent.OperatorAgentLineAllocation) *agentlineallocationpb.AgentLineAllocationInfo {
	return &agentlineallocationpb.AgentLineAllocationInfo{
		Id:            data.ID,                    // 分配记录ID
		OperatorId:    data.OperatorID,            // 分站ID
		AgentLineCode: data.AgentLineCode,         // 代理子线路编码
		CreatedAt:     data.CreatedAt.UnixMilli(), // 分配时间, Unix毫秒时间戳
	}
}
