package operatorprofileservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/profile"
)

// toOperatorProfileInfo 转换分站档案信息
func toOperatorProfileInfo(data *ent.OperatorProfile) *profile.OperatorProfileInfo {
	return &profile.OperatorProfileInfo{
		Id:           data.ID,                    // 档案ID
		OperatorId:   data.OperatorID,            // 分站ID
		CompanyName:  data.CompanyName,           // 公司名称
		ContactName:  data.ContactName,           // 主要联系人名称
		ContactEmail: data.ContactEmail,          // 主要联系人邮箱
		Remark:       data.Remark,                // 总网内部档案备注
		CreatedAt:    data.CreatedAt.UnixMilli(), // 创建时间
		UpdatedAt:    data.UpdatedAt.UnixMilli(), // 更新时间
	}
}
