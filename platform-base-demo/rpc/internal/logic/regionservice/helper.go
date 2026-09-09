package regionservicelogic

import (
	"oa.98ent.com/p9/platform-base/rpc/ent"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/region"
)

// toRegionInfo 转换国家地区信息
func toRegionInfo(data *ent.Region) *region.RegionInfo {
	return &region.RegionInfo{
		Id:          data.ID,          // 国家地区ID
		Code:        data.Code,        // 国家地区编码
		CallingCode: data.CallingCode, // 国际电话区号
		NameKey:     data.NameKey,     // 名称翻译Key
		Status:      data.Status,      // 状态: 1启用, 2停用
		SortNo:      data.SortNo,      // 排序值
	}
}
