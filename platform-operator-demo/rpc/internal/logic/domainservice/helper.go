package domainservicelogic

import (
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/domain"
)

// toDomainInfo 转换分站域名信息
func toDomainInfo(data *ent.OperatorDomain) *domain.DomainInfo {
	return &domain.DomainInfo{
		Id:         data.ID,                    // 域名ID
		OperatorId: data.OperatorID,            // 分站ID
		DomainName: data.DomainName,            // 域名
		DomainType: data.DomainType,            // 域名类型: 1分站后台, 2代理后台, 3会员H5
		Status:     data.Status,                // 域名状态: 1启用, 2停用
		Remark:     data.Remark,                // 总网内部备注
		CreatedAt:  data.CreatedAt.UnixMilli(), // 创建时间
		UpdatedAt:  data.UpdatedAt.UnixMilli(), // 更新时间
	}
}
