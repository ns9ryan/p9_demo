package i18n

import (
	"context"
	"strings"
)

const (
	LangZH = "zh-CN"
	LangHK = "zh-HK"
	LangEN = "en-US"
)

type langKey struct{}

// WithLang 设置语言到上下文
func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langKey{}, lang)
}

// Lang 获取语言，默认返回英文
func Lang(ctx context.Context, defaultLang *string) string {
	if ctx == nil {
		if defaultLang != nil {
			return *defaultLang
		}
		return LangEN
	}
	if v, ok := ctx.Value(langKey{}).(string); ok && v != "" {
		return v
	}
	if defaultLang != nil {
		return *defaultLang
	}
	return LangEN
}

// ParseLang 解析语言
func ParseLang(header string) string {
	if header == "" {
		return LangEN
	}
	first := strings.TrimSpace(strings.Split(header, ",")[0])
	first = strings.TrimSpace(strings.Split(first, ";")[0])
	if first == "" {
		return LangZH
	}
	lower := strings.ToLower(first)
	if strings.HasPrefix(lower, "en") {
		return LangEN
	}
	if strings.HasPrefix(lower, "zh") {
		if strings.HasPrefix(lower, "zh-hk") {
			return LangHK
		}
		return LangZH
	}
	return first
}
