package operator

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/operator-base/rpc/ent"
	operatorent "oa.98ent.com/p9/operator-base/rpc/ent/operator"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// normalizeOperator 整理 operator 初始化基础数据
func normalizeOperator(req *InitializeRequest) {
	req.Code = strings.TrimSpace(req.Code)
	req.Name = strings.TrimSpace(req.Name)
	req.TimezoneCode = strings.TrimSpace(req.TimezoneCode)
	req.SettlementCurrencyCode = strings.TrimSpace(req.SettlementCurrencyCode)
}

// validateOperator 校验 operator 初始化基础数据
func validateOperator(req InitializeRequest) error {
	if req.Code == "" {
		return status.Error(codes.InvalidArgument, "operator code is required")
	}
	if req.Name == "" {
		return status.Error(codes.InvalidArgument, "operator name is required")
	}
	if req.TimezoneCode == "" {
		return status.Error(codes.InvalidArgument, "timezone code is required")
	}
	if req.SettlementCurrencyCode == "" {
		return status.Error(codes.InvalidArgument, "settlement currency code is required")
	}
	if req.Status < 1 || req.Status > 3 {
		return status.Error(codes.InvalidArgument, "operator status is invalid")
	}

	return nil
}

// getOperatorByCode 按业务编码查询 operator
func (s *Service) getOperatorByCode(ctx context.Context, code string) (*ent.Operator, bool, error) {
	data, err := s.db.Operator.
		Query().
		Where(operatorent.CodeEQ(code)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, false, nil
		}

		return nil, false, fmt.Errorf("查询 operator 失败: %w", err)
	}

	return data, true, nil
}

// checkExistingOperator 检查已存在 operator 是否符合初始化幂等要求
func checkExistingOperator(data *ent.Operator, req InitializeRequest) error {
	// 时区和结算货币首次初始化后不可修改
	if data.TimezoneCode != req.TimezoneCode ||
		data.SettlementCurrencyCode != req.SettlementCurrencyCode {
		return status.Error(codes.FailedPrecondition, "operator initialization data conflicts with existing operator")
	}

	// name、status 和资源可能已经在初始化后发生变化, 不使用旧初始化数据覆盖
	return nil
}

// createOperator 创建 operator
func createOperator(ctx context.Context, tx *ent.Tx, req InitializeRequest) (*ent.Operator, error) {
	data, err := tx.Operator.
		Create().
		SetCode(req.Code).                                     // operator 全局唯一业务编码
		SetName(req.Name).                                     // operator 名称
		SetTimezoneCode(req.TimezoneCode).                     // IANA 时区编码
		SetSettlementCurrencyCode(req.SettlementCurrencyCode). // 结算货币编码
		SetStatus(req.Status).                                 // operator 状态: 1正常, 2暂停, 3关闭
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("创建 operator 失败: %w", err)
	}

	return data, nil
}
