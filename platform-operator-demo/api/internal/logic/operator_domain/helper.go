package operator_domain

import (
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/domainpb"
)

// toOperatorDomainInfo 转换分站域名信息
func toOperatorDomainInfo(data *domainpb.DomainInfo) types.OperatorDomainInfo {
	return types.OperatorDomainInfo{
		Id:         data.Id,         // 域名ID
		OperatorId: data.OperatorId, // 分站ID
		DomainName: data.DomainName, // 域名, 不包含协议和端口
		DomainType: data.DomainType, // 域名类型: 1分站后台, 2代理后台, 3会员H5
		Status:     data.Status,     // 域名状态: 1启用, 2停用
		Remark:     data.Remark,     // 总网内部备注
		CreatedAt:  data.CreatedAt,  // 创建时间, Unix毫秒时间戳
		UpdatedAt:  data.UpdatedAt,  // 更新时间, Unix毫秒时间戳
	}
}
