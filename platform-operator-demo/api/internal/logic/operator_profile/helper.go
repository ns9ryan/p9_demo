package operator_profile

import (
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/profilepb"
)

// toOperatorProfileInfo 转换分站档案信息
func toOperatorProfileInfo(data *profilepb.OperatorProfileInfo) *types.OperatorProfileInfo {
	if data == nil {
		return nil
	}

	return &types.OperatorProfileInfo{
		Id:           data.Id,           // 档案ID
		OperatorId:   data.OperatorId,   // 分站ID
		CompanyName:  data.CompanyName,  // 公司名称
		ContactName:  data.ContactName,  // 主要联系人名称
		ContactEmail: data.ContactEmail, // 主要联系人邮箱
		Remark:       data.Remark,       // 总网内部档案备注
		CreatedAt:    data.CreatedAt,    // 创建时间, Unix毫秒时间戳
		UpdatedAt:    data.UpdatedAt,    // 更新时间, Unix毫秒时间戳
	}
}
