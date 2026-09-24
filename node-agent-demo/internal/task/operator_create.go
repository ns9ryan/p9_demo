package task

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	operatorbasepb "oa.98ent.com/p9/operator-base/rpc/pb/operatorbaserpc/operatorpb"
	platformoperatorpb "oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"
)

// createOperatorParams 创建分站任务参数
type createOperatorParams struct {
	OperatorCode string `json:"operator_code"` // 分站全局唯一业务编码
}

// executeCreateOperator 执行创建分站任务
func (s *Service) executeCreateOperator(ctx context.Context, params json.RawMessage) (json.RawMessage, error) {
	// 解析任务参数
	var taskParams createOperatorParams
	if err := json.Unmarshal(params, &taskParams); err != nil {
		return nil, fmt.Errorf("解析创建分站任务参数失败: %w", err)
	}

	// 整理并校验分站业务编码
	operatorCode := strings.TrimSpace(taskParams.OperatorCode)
	if operatorCode == "" {
		return nil, fmt.Errorf("operator_code不能为空")
	}

	// 获取总网当前分站初始化数据
	initializationData, err := s.platformOperatorRpc.GetInitializationData(ctx, &platformoperatorpb.GetInitializationDataRequest{
		OperatorCode: operatorCode, // 分站全局唯一业务编码
	})
	if err != nil {
		return nil, fmt.Errorf("获取分站初始化数据失败: %w", err)
	}

	// 获取分站基础信息
	operatorInfo := initializationData.GetOperator()
	if operatorInfo == nil {
		return nil, fmt.Errorf("分站初始化基础信息不能为空")
	}

	// 转换域名初始化数据
	domains := make([]*operatorbasepb.OperatorDomainInitializationInfo, 0, len(initializationData.GetDomains()))
	for _, item := range initializationData.GetDomains() {
		domains = append(domains, &operatorbasepb.OperatorDomainInitializationInfo{
			DomainName: item.GetDomainName(), // 域名
			DomainType: item.GetDomainType(), // 域名类型
		})
	}

	// 初始化当前节点分站数据
	_, err = s.operatorBaseRpc.Initialize(ctx, &operatorbasepb.InitializeOperatorRequest{
		Operator: &operatorbasepb.OperatorInitializationInfo{
			Code:                   operatorInfo.GetCode(),                   // 分站全局唯一业务编码
			Name:                   operatorInfo.GetName(),                   // 分站名称
			TimezoneCode:           operatorInfo.GetTimezoneCode(),           // IANA时区编码
			SettlementCurrencyCode: operatorInfo.GetSettlementCurrencyCode(), // 结算币种编码
			Status:                 operatorInfo.GetStatus(),                 // 分站状态
		},
		Domains:        domains,                                // 当前有效域名
		LanguageCodes:  initializationData.GetLanguageCodes(),  // 当前有效语言编码
		RegionCodes:    initializationData.GetRegionCodes(),    // 当前有效经营地区编码
		AgentLineCodes: initializationData.GetAgentLineCodes(), // 当前有效代理子线路编码
	})
	if err != nil {
		return nil, fmt.Errorf("初始化当前节点分站数据失败: %w", err)
	}

	// 当前任务无需返回业务结果
	return nil, nil
}
