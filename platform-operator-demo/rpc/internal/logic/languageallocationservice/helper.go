package languageallocationservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/languageallocationpb"
)

// toLanguageAllocationInfo 转换语言分配信息
func toLanguageAllocationInfo(data *ent.OperatorLanguageAllocation) *languageallocationpb.LanguageAllocationInfo {
	return &languageallocationpb.LanguageAllocationInfo{
		Id:           data.ID,                    // 分配记录ID
		OperatorId:   data.OperatorID,            // 分站ID
		LanguageCode: data.LanguageCode,          // 语言编码
		CreatedAt:    data.CreatedAt.UnixMilli(), // 分配时间, Unix毫秒时间戳
	}
}
