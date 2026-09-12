package languageallocationservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/languageallocation"
)

// toLanguageAllocationInfo 转换语言分配信息
func toLanguageAllocationInfo(data *ent.OperatorLanguageAllocation) *languageallocation.LanguageAllocationInfo {
	return &languageallocation.LanguageAllocationInfo{
		Id:           data.ID,                    // 分配记录ID
		OperatorId:   data.OperatorID,            // 分站ID
		LanguageCode: data.LanguageCode,          // 语言编码
		CreatedAt:    data.CreatedAt.UnixMilli(), // 分配时间
	}
}
