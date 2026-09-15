package middleware

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest"
	corei18n "oa.98ent.com/p9/core/common/i18n"
)

// LanguageMiddleware 语言中间件 - 从请求头读取语言并设置到context
func LanguageMiddleware() rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// 从请求头读取语言，支持多种格式：X-Lang、Accept-Language
			lang := r.Header.Get("X-Lang")
			if lang == "" {
				lang = r.Header.Get("Accept-Language")
			}

			// 解析并规范化语言代码
			lang = corei18n.ParseLang(lang)

			// 将语言设置到context中
			ctx := corei18n.WithLang(r.Context(), lang)

			// 继续处理请求
			next(w, r.WithContext(ctx))
		}
	}
}
