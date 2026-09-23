package operator

import (
	"context"
	"fmt"

	"oa.98ent.com/p9/operator-base/rpc/ent"
)

// InitializeRequest operator 初始化请求
type InitializeRequest struct {
	Code                   string             // operator 全局唯一业务编码
	Name                   string             // operator 名称
	TimezoneCode           string             // IANA 时区编码
	SettlementCurrencyCode string             // 结算货币编码
	Status                 int64              // operator 状态: 1正常, 2暂停, 3关闭
	Domains                []InitializeDomain // 当前有效域名
	LanguageCodes          []string           // 当前有效语言编码
	RegionCodes            []string           // 当前有效经营地区编码
	AgentLineCodes         []string           // 当前有效代理子线路编码
}

// Initialize 初始化 operator
func (s *Service) Initialize(ctx context.Context, req InitializeRequest) error {
	// 整理初始化数据
	normalizeInitializeRequest(&req)

	// 校验初始化数据
	if err := validateInitializeRequest(req); err != nil {
		return err
	}

	// 检查 operator 是否已经初始化
	existing, exists, err := s.getOperatorByCode(ctx, req.Code)
	if err != nil {
		return err
	}
	if exists {
		return checkExistingOperator(existing, req)
	}

	// 开启初始化事务
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("开启初始化事务失败: %w", err)
	}

	// 创建 operator 基础数据
	operatorData, err := createOperator(ctx, tx, req)
	if err != nil {
		_ = tx.Rollback()

		// 并发初始化相同 operator 时, 由 code 唯一约束完成最终幂等保护
		if ent.IsConstraintError(err) {
			existing, exists, queryErr := s.getOperatorByCode(ctx, req.Code)
			if queryErr != nil {
				return queryErr
			}
			if exists {
				return checkExistingOperator(existing, req)
			}
		}

		return err
	}

	// 创建当前有效域名
	if err = createDomains(ctx, tx, operatorData.ID, req.Domains); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 创建当前有效语言
	if err = createLanguages(ctx, tx, operatorData.ID, req.LanguageCodes); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 创建当前有效经营地区
	if err = createRegions(ctx, tx, operatorData.ID, req.RegionCodes); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 创建当前有效代理子线路
	if err = createAgentLines(ctx, tx, operatorData.ID, req.AgentLineCodes); err != nil {
		_ = tx.Rollback()
		return err
	}

	// 提交初始化事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交初始化事务失败: %w", err)
	}

	return nil
}

// normalizeInitializeRequest 整理 operator 初始化数据
func normalizeInitializeRequest(req *InitializeRequest) {
	normalizeOperator(req)
	normalizeDomains(req.Domains)
	normalizeLanguageCodes(req.LanguageCodes)
	normalizeRegionCodes(req.RegionCodes)
	normalizeAgentLineCodes(req.AgentLineCodes)
}

// validateInitializeRequest 校验 operator 初始化数据
func validateInitializeRequest(req InitializeRequest) error {
	// 校验 operator 基础数据
	if err := validateOperator(req); err != nil {
		return err
	}

	// 校验域名
	if err := validateDomains(req.Domains); err != nil {
		return err
	}

	// 校验语言
	if err := validateLanguageCodes(req.LanguageCodes); err != nil {
		return err
	}

	// 校验经营地区
	if err := validateRegionCodes(req.RegionCodes); err != nil {
		return err
	}

	// 校验代理子线路
	if err := validateAgentLineCodes(req.AgentLineCodes); err != nil {
		return err
	}

	return nil
}
