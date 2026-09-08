package i18n

import "context"

// languageContextKey 用于在上下文中保存语言
type languageContextKey struct{}

// WithLanguage 将语言写入上下文
func WithLanguage(ctx context.Context, language string) context.Context {
	return context.WithValue(ctx, languageContextKey{}, language)
}

// LanguageFromContext 从上下文获取语言
func LanguageFromContext(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	language, _ := ctx.Value(languageContextKey{}).(string)
	return language
}
