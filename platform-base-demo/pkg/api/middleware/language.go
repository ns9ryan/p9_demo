package middleware

import (
	"net/http"

	"oa.98ent.com/p9/platform-base/pkg/i18n"
)

// LanguageMiddleware 语言中间件
type LanguageMiddleware struct{}

// NewLanguageMiddleware 创建语言中间件
func NewLanguageMiddleware() *LanguageMiddleware {
	return &LanguageMiddleware{}
}

// Handle 处理请求语言
func (m *LanguageMiddleware) Handle(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 获取请求语言并写入上下文
		language := r.Header.Get("X-Lang")
		ctx := i18n.WithLanguage(r.Context(), language)

		next(w, r.WithContext(ctx))
	}
}
