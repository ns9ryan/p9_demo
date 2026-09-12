package operatorservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/operator"
)

// toOperatorInfo 转换分站信息
func toOperatorInfo(data *ent.Operator) *operator.OperatorInfo {
	var publishedAt *int64
	if data.PublishedAt != nil {
		value := data.PublishedAt.UnixMilli()
		publishedAt = &value
	}

	return &operator.OperatorInfo{
		Id:                     data.ID,                     // 分站ID
		Code:                   data.Code,                   // 分站业务编码
		Name:                   data.Name,                   // 分站名称
		TimezoneCode:           data.TimezoneCode,           // 时区编码
		SettlementCurrencyCode: data.SettlementCurrencyCode, // 结算币种编码
		CreationStatus:         data.CreationStatus,         // 创建状态: 1草稿, 2已完成
		PublishStatus:          data.PublishStatus,          // 发布状态: 1未发布, 2发布中, 3已发布, 4发布失败
		Status:                 data.Status,                 // 分站状态: 1正常, 2暂停, 3关闭
		Remark:                 data.Remark,                 // 内部备注
		PublishedAt:            publishedAt,                 // 首次发布成功时间
		CreatedAt:              data.CreatedAt.UnixMilli(),  // 创建时间
		UpdatedAt:              data.UpdatedAt.UnixMilli(),  // 更新时间
	}
}
