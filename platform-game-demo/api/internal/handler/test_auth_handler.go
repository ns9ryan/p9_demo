package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/platform-game/api/internal/logic"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
	"oa.98ent.com/p9/platform-game/api/internal/types"
	"oa.98ent.com/p9/platform-game/common/logger"
)

// TestAuthCheckHandler 测试鉴权检查接口处理器
// 用于验证从 core 服务获取的 token 是否有效
func TestAuthCheckHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewTestAuthCheckLogic(r.Context(), svcCtx)
		resp, err := l.TestAuthCheck(&types.TestAuthCheckReq{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		// 返回 JSON 响应
		data := map[string]interface{}{
			"code":    1,
			"message": "success",
			"data":    resp,
		}

		json.NewEncoder(w).Encode(data)
	}
}

// TestAuthWithTokenHandler 测试使用 token 的接口
// 这个接口需要有效的 token（通过 JWT 中间件验证）
func TestAuthWithTokenHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 手动从 Authorization 头中提取和验证 token
		// （因为这个路由在 routes.go 中没有被 JWT 中间件包装）
		if svcCtx.CoreAuth != nil {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				respondError(w, "missing authorization header", http.StatusUnauthorized)
				return
			}

			// 提取 Bearer token
			token := extractBearerToken(authHeader)
			if token == "" {
				respondError(w, "invalid authorization header format", http.StatusUnauthorized)
				return
			}

			// 验证 token
			logger.Infof("[TestAuthWithTokenHandler] 验证 Token: %s", token[:20]+"...")
			claims, err := svcCtx.CoreAuth.CheckToken(r.Context(), token)
			if err != nil {
				logger.Warnf("[TestAuthWithTokenHandler] Token 验证失败: %v", err)
				respondError(w, "token validation failed: "+err.Error(), http.StatusUnauthorized)
				return
			}

			// 将 claims 放入 context
			ctx := ctxdata.WithClaims(r.Context(), claims)
			ctx = ctxdata.WithRawToken(ctx, token)
			r = r.WithContext(ctx)
		}

		// 调用实际的 logic
		l := logic.NewTestAuthWithTokenLogic(r.Context(), svcCtx)
		resp, err := l.TestAuthWithToken()
		if err != nil {
			respondError(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		data := map[string]interface{}{
			"code":    1,
			"message": "success",
			"data":    resp,
		}

		json.NewEncoder(w).Encode(data)
	}
}

// extractBearerToken 从 Authorization 头中提取 Bearer token
func extractBearerToken(authHeader string) string {
	authHeader = strings.TrimSpace(authHeader)
	if len(authHeader) > 7 && strings.EqualFold(authHeader[:7], "Bearer ") {
		return strings.TrimSpace(authHeader[7:])
	}
	return ""
}

// respondError 返回错误响应
func respondError(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	data := map[string]interface{}{
		"code":    statusCode,
		"message": message,
		"data":    nil,
	}

	json.NewEncoder(w).Encode(data)
}
