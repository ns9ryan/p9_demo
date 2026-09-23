package operatorservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	operatorent "oa.98ent.com/p9/platform-operator/rpc/ent/operator"
	operatoragentlineallocationent "oa.98ent.com/p9/platform-operator/rpc/ent/operatoragentlineallocation"
	operatordomainent "oa.98ent.com/p9/platform-operator/rpc/ent/operatordomain"
	operatorlanguageallocationent "oa.98ent.com/p9/platform-operator/rpc/ent/operatorlanguageallocation"
	operatorregionallocationent "oa.98ent.com/p9/platform-operator/rpc/ent/operatorregionallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetInitializationDataLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetInitializationDataLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetInitializationDataLogic {
	return &GetInitializationDataLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetInitializationData 获取分站初始化数据
func (l *GetInitializationDataLogic) GetInitializationData(in *operatorpb.GetInitializationDataRequest) (*operatorpb.GetInitializationDataResponse, error) {
	// 整理分站业务编码
	operatorCode := strings.TrimSpace(in.GetOperatorCode())
	if operatorCode == "" {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 获取分站基础信息
	operatorData, err := l.svcCtx.DB.Operator.
		Query().
		Where(operatorent.CodeEQ(operatorCode)).
		Only(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 获取当前启用域名
	domainData, err := l.svcCtx.DB.OperatorDomain.
		Query().
		Where(
			operatordomainent.OperatorIDEQ(operatorData.ID),
			operatordomainent.StatusEQ(1),
		).
		Order(operatordomainent.ByDomainType()).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 初始化时三种域名必须全部存在
	if len(domainData) != 3 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 获取当前语言分配
	languageCodes, err := l.svcCtx.DB.OperatorLanguageAllocation.
		Query().
		Where(operatorlanguageallocationent.OperatorIDEQ(operatorData.ID)).
		Order(operatorlanguageallocationent.ByLanguageCode()).
		Select(operatorlanguageallocationent.FieldLanguageCode).
		Strings(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 初始化时至少需要一种语言
	if len(languageCodes) == 0 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 获取当前经营地区分配
	regionCodes, err := l.svcCtx.DB.OperatorRegionAllocation.
		Query().
		Where(operatorregionallocationent.OperatorIDEQ(operatorData.ID)).
		Order(operatorregionallocationent.ByRegionCode()).
		Select(operatorregionallocationent.FieldRegionCode).
		Strings(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 初始化时至少需要一个经营地区
	if len(regionCodes) == 0 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 获取当前代理子线路分配
	agentLineCodes, err := l.svcCtx.DB.OperatorAgentLineAllocation.
		Query().
		Where(operatoragentlineallocationent.OperatorIDEQ(operatorData.ID)).
		Order(operatoragentlineallocationent.ByAgentLineCode()).
		Select(operatoragentlineallocationent.FieldAgentLineCode).
		Strings(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 初始化时至少需要一条代理子线路
	if len(agentLineCodes) == 0 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 转换当前启用域名
	domains := make([]*operatorpb.DomainInitializationInfo, 0, len(domainData))
	for _, data := range domainData {
		domains = append(domains, &operatorpb.DomainInitializationInfo{
			DomainName: data.DomainName, // 域名
			DomainType: data.DomainType, // 域名类型
		})
	}

	// 返回分站初始化数据
	return &operatorpb.GetInitializationDataResponse{
		Operator: &operatorpb.InitializationInfo{
			Code:                   operatorData.Code,                   // 分站全局唯一业务编码
			Name:                   operatorData.Name,                   // 分站名称
			TimezoneCode:           operatorData.TimezoneCode,           // IANA 时区编码
			SettlementCurrencyCode: operatorData.SettlementCurrencyCode, // 结算币种编码
			Status:                 operatorData.Status,                 // 分站状态: 1正常, 2暂停, 3关闭
		},
		Domains:        domains,        // 当前启用域名
		LanguageCodes:  languageCodes,  // 当前分配的语言编码
		RegionCodes:    regionCodes,    // 当前分配的经营地区编码
		AgentLineCodes: agentLineCodes, // 当前分配的代理子线路编码
	}, nil
}
