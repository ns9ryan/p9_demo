package operator_admin

import (
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"
)

// toOperatorAdminInfo 转换分站管理员信息
func toOperatorAdminInfo(data *adminpb.AdminInfo) types.OperatorAdminInfo {
	return types.OperatorAdminInfo{
		Id:           data.Id,           // 管理员ID
		OperatorId:   data.OperatorId,   // 分站ID
		OperatorName: data.OperatorName, // 分站名称
		Username:     data.Username,     // 账号
		DisplayName:  data.DisplayName,  // 显示名称
		Status:       data.Status,       // 状态: 1启用, 2停用
		CreatedAt:    data.CreatedAt,    // 创建时间, Unix毫秒时间戳
		UpdatedAt:    data.UpdatedAt,    // 更新时间, Unix毫秒时间戳
	}
}
