package regionallocationservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/regionallocationpb"
)

// toRegionAllocationInfo 转换经营地区分配信息
func toRegionAllocationInfo(data *ent.OperatorRegionAllocation) *regionallocationpb.RegionAllocationInfo {
	return &regionallocationpb.RegionAllocationInfo{
		Id:         data.ID,                    // 分配记录ID
		OperatorId: data.OperatorID,            // 分站ID
		RegionCode: data.RegionCode,            // 国家地区编码
		CreatedAt:  data.CreatedAt.UnixMilli(), // 分配时间, Unix毫秒时间戳
	}
}
