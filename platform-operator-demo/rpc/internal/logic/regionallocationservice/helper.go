package regionallocationservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/regionallocation"
)

// toRegionAllocationInfo 转换经营地区分配信息
func toRegionAllocationInfo(data *ent.OperatorRegionAllocation) *regionallocation.RegionAllocationInfo {
	return &regionallocation.RegionAllocationInfo{
		Id:         data.ID,                    // 分配记录ID
		OperatorId: data.OperatorID,            // 分站ID
		RegionCode: data.RegionCode,            // 经营地区编码
		CreatedAt:  data.CreatedAt.UnixMilli(), // 分配时间
	}
}
