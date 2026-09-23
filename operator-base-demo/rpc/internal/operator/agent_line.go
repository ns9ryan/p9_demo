package operator

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/operator-base/rpc/ent"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// normalizeAgentLineCodes 整理 operator 初始化代理子线路编码
func normalizeAgentLineCodes(agentLineCodes []string) {
	for i := range agentLineCodes {
		agentLineCodes[i] = strings.TrimSpace(agentLineCodes[i])
	}
}

// validateAgentLineCodes 校验 operator 初始化代理子线路编码
func validateAgentLineCodes(agentLineCodes []string) error {
	// 初始化时至少需要一条代理子线路
	if len(agentLineCodes) == 0 {
		return status.Error(codes.InvalidArgument, "agent_line_codes is required")
	}

	seen := make(map[string]struct{}, len(agentLineCodes))

	for _, agentLineCode := range agentLineCodes {
		// 校验代理子线路编码
		if agentLineCode == "" {
			return status.Error(codes.InvalidArgument, "agent_line_codes contains empty value")
		}

		// 同一个 operator 不能分配重复代理子线路
		if _, exists := seen[agentLineCode]; exists {
			return status.Error(codes.InvalidArgument, "agent_line_codes contains duplicate value")
		}
		seen[agentLineCode] = struct{}{}
	}

	return nil
}

// createAgentLines 创建 operator 当前有效代理子线路
func createAgentLines(ctx context.Context, tx *ent.Tx, operatorID int64, agentLineCodes []string) error {
	builders := make([]*ent.OperatorAgentLineCreate, 0, len(agentLineCodes))

	for _, agentLineCode := range agentLineCodes {
		builders = append(builders,
			tx.OperatorAgentLine.
				Create().
				SetOperatorID(operatorID).       // operator 本地主键
				SetAgentLineCode(agentLineCode), // 代理子线路编码
		)
	}

	if err := tx.OperatorAgentLine.CreateBulk(builders...).Exec(ctx); err != nil {
		return fmt.Errorf("创建 operator 代理子线路失败: %w", err)
	}

	return nil
}
