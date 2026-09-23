package utils

import (
	"context"

	corei18n "oa.98ent.com/p9/common/i18n"
)

func TGPlatformGame(ctx context.Context, nameKey string) string {
	name := corei18n.TG(ctx, corei18n.CodePlatform, "game", nameKey)
	// 如果翻译结果与nameKey相同，说明没有翻译，尝试使用中文进行翻译
	if name == nameKey {
		ctx = corei18n.WithLang(ctx, "zh-CN")
		name = corei18n.TG(ctx, corei18n.CodePlatform, "game", nameKey)
	}
	return name
}

func TGPlatformBase(ctx context.Context, nameKey string) string {
	name := corei18n.TG(ctx, corei18n.CodePlatform, "base", nameKey)
	return name
}
