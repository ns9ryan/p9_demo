package operator

import (
	"context"
	"fmt"
	"strings"

	"oa.98ent.com/p9/operator-base/rpc/ent"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// normalizeLanguageCodes 整理 operator 初始化语言编码
func normalizeLanguageCodes(languageCodes []string) {
	for i := range languageCodes {
		languageCodes[i] = strings.TrimSpace(languageCodes[i])
	}
}

// validateLanguageCodes 校验 operator 初始化语言编码
func validateLanguageCodes(languageCodes []string) error {
	// 初始化时至少需要一种语言
	if len(languageCodes) == 0 {
		return status.Error(codes.InvalidArgument, "language_codes is required")
	}

	seen := make(map[string]struct{}, len(languageCodes))

	for _, languageCode := range languageCodes {
		// 校验语言编码
		if languageCode == "" {
			return status.Error(codes.InvalidArgument, "language_codes contains empty value")
		}

		// 同一个 operator 不能分配重复语言
		if _, exists := seen[languageCode]; exists {
			return status.Error(codes.InvalidArgument, "language_codes contains duplicate value")
		}
		seen[languageCode] = struct{}{}
	}

	return nil
}

// createLanguages 创建 operator 当前有效语言
func createLanguages(ctx context.Context, tx *ent.Tx, operatorID int64, languageCodes []string) error {
	builders := make([]*ent.OperatorLanguageCreate, 0, len(languageCodes))

	for _, languageCode := range languageCodes {
		builders = append(builders,
			tx.OperatorLanguage.
				Create().
				SetOperatorID(operatorID).     // operator 本地主键
				SetLanguageCode(languageCode), // 语言编码
		)
	}

	if err := tx.OperatorLanguage.CreateBulk(builders...).Exec(ctx); err != nil {
		return fmt.Errorf("创建 operator 语言失败: %w", err)
	}

	return nil
}
