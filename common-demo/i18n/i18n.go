package i18n

import (
	"context"
	"net/http"
	"strings"
)

const (
	LangZH = "zh-CN"
	LangHK = "zh-HK"
	LangEN = "en-US"
)

// langKey 语言上下文键
type langKey struct{}

// I18nLangMiddleware 国际化语言中间件
type I18nLangMiddleware struct{ defaultLang string }

func NewI18nLangMiddleware(defaultLang string) *I18nLangMiddleware {
	return &I18nLangMiddleware{defaultLang: defaultLang}
}

func (m *I18nLangMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := ParseLangWithDefault(r.Header.Get("X-Lang"), m.defaultLang)
		next(w, r.WithContext(WithLang(r.Context(), lang)))
	}
}

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

// ParseLangWithDefault 解析语言，默认返回默认语言
func ParseLangWithDefault(header string, defaultLang string) string {
	if header == "" {
		return defaultLang
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
