package task

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const (
	TargetOperatorBase = "operator-base"   // 任务目标: 分站基础服务
	TypeCreateOperator = "CREATE_OPERATOR" // 任务类型: 创建分站
)

// Execute 执行调度任务
func (s *Service) Execute(ctx context.Context, target string, taskType string, params json.RawMessage) (json.RawMessage, error) {
	// 整理任务信息
	target = strings.TrimSpace(target)
	taskType = strings.TrimSpace(taskType)

	// 根据目标服务执行任务
	switch target {
	case TargetOperatorBase:
		return s.executeOperatorBaseTask(ctx, taskType, params)

	default:
		return nil, fmt.Errorf("不支持的任务目标服务: %s", target)
	}
}

// executeOperatorBaseTask 执行分站基础服务任务
func (s *Service) executeOperatorBaseTask(ctx context.Context, taskType string, params json.RawMessage) (json.RawMessage, error) {
	// 根据任务类型执行任务
	switch taskType {
	case TypeCreateOperator:
		return s.executeCreateOperator(ctx, params)

	default:
		return nil, fmt.Errorf("不支持的operator-base任务类型: %s", taskType)
	}
}
