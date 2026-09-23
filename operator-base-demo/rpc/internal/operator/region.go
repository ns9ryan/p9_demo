package operator

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/operator-base/rpc/ent"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// normalizeRegionCodes 整理 operator 初始化经营地区编码
func normalizeRegionCodes(regionCodes []string) {
	for i := range regionCodes {
		regionCodes[i] = strings.ToUpper(strings.TrimSpace(regionCodes[i]))
	}
}

// validateRegionCodes 校验 operator 初始化经营地区编码
func validateRegionCodes(regionCodes []string) error {
	// 初始化时至少需要一个经营地区
	if len(regionCodes) == 0 {
		return status.Error(codes.InvalidArgument, "region_codes is required")
	}

	seen := make(map[string]struct{}, len(regionCodes))

	for _, regionCode := range regionCodes {
		// 校验经营地区编码
		if regionCode == "" {
			return status.Error(codes.InvalidArgument, "region_codes contains empty value")
		}

		// 同一个 operator 不能分配重复经营地区
		if _, exists := seen[regionCode]; exists {
			return status.Error(codes.InvalidArgument, "region_codes contains duplicate value")
		}
		seen[regionCode] = struct{}{}
	}

	return nil
}

// createRegions 创建 operator 当前有效经营地区
func createRegions(ctx context.Context, tx *ent.Tx, operatorID int64, regionCodes []string) error {
	builders := make([]*ent.OperatorRegionCreate, 0, len(regionCodes))

	for _, regionCode := range regionCodes {
		builders = append(builders,
			tx.OperatorRegion.
				Create().
				SetOperatorID(operatorID). // operator 本地主键
				SetRegionCode(regionCode), // 经营地区编码
		)
	}

	if err := tx.OperatorRegion.CreateBulk(builders...).Exec(ctx); err != nil {
		return fmt.Errorf("创建 operator 经营地区失败: %w", err)
	}

	return nil
}
