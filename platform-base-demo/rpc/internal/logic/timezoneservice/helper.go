package timezoneservicelogic

import (
	"oa.98ent.com/p9/platform-base/rpc/ent"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/timezone"
)

// toTimezoneInfo 转换时区信息
func toTimezoneInfo(data *ent.Timezone) *timezone.TimezoneInfo {
	return &timezone.TimezoneInfo{
		Id:       data.ID,       // 时区ID
		Code:     data.Code,     // IANA时区编码
		NameI18N: data.NameI18n, // 多语言名称
		Status:   data.Status,   // 状态: 1启用, 2停用
		SortNo:   data.SortNo,   // 排序值
	}
}
