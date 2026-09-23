package operatoradminservicelogic

import (
	"strings"

	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"
)

// trimOptionalString 整理可选字符串
func trimOptionalString(value *string) *string {
	if value == nil {
		return nil
	}

	return new(strings.TrimSpace(*value))
}

// toAdminInfo 转换分站管理员信息
func toAdminInfo(data *ent.OperatorAdmin) *adminpb.AdminInfo {
	var operatorName string
	if data.Edges.Operator != nil {
		operatorName = data.Edges.Operator.Name
	}

	return &adminpb.AdminInfo{
		Id:           data.ID,                    // 管理员ID
		OperatorId:   data.OperatorID,            // 分站ID
		OperatorName: operatorName,               // 分站名称
		Username:     data.Username,              // 账号
		DisplayName:  data.DisplayName,           // 显示名称
		Status:       data.Status,                // 状态: 1启用, 2停用
		CreatedAt:    data.CreatedAt.UnixMilli(), // 创建时间, Unix毫秒时间戳
		UpdatedAt:    data.UpdatedAt.UnixMilli(), // 更新时间, Unix毫秒时间戳
	}
}

