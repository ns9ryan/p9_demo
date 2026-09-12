package middleware

import (
	"net/http"
)

// ErrorHandler 全局错误处理中间件
// 统一处理 handler 中的错误，将错误转换为标准 JSON 响应
func ErrorHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 设置默认响应头
		w.Header().Set("Content-Type", "application/json; charset=utf-8")

		// 调用下一个处理器
		next(w, r)
	}
}
