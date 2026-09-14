package middleware

import (
	"context"
	"net/http"
	"strings"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/response"
	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/common/xerr"

	"github.com/zeromicro/go-zero/rest"
)

// Client 客户端接口
type Client interface {
	CheckToken(ctx context.Context, accessToken string) (*ctxdata.Claims, error)
	Enforce(ctx context.Context, claims *ctxdata.Claims, path, method string) (bool, error)
}

// I18n 国际化
func I18n(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lang := i18n.ParseLang(r.Header.Get("X-Lang"))
		next(w, r.WithContext(i18n.WithLang(r.Context(), lang)))
	}
}

// ClientIP 把客户端 IP 写入上下文（含 gRPC outgoing metadata）。
func ClientIP(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r.WithContext(ctxdata.WithClientIP(r.Context(), utils.ClientIP(r))))
	}
}

// JWT 认证
func JWT(c Client) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			raw := stripBearer(r.Header.Get("Authorization"))
			ctx := r.Context()
			if ctxdata.ClientIPFromCtx(ctx) == "" {
				ctx = ctxdata.WithClientIP(ctx, utils.ClientIP(r))
			}
			claims, err := c.CheckToken(ctx, raw)
			if err != nil {
				response.FailCtx(ctx, w, err)
				return
			}
			if claims != nil && claims.TokenType == jwt.TokenPreview && previewWriteDenied(r.Method, r.URL.Path) {
				response.FailCtx(ctx, w, xerr.Forbidden(i18n.AuthPreviewReadOnly))
				return
			}
			ctx = ctxdata.WithClaims(ctx, claims)
			ctx = ctxdata.WithRawToken(ctx, raw)
			next(w, r.WithContext(ctx))
		}
	}
}

// Authority 权限控制
func Authority(c Client) rest.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			claims := ctxdata.ClaimsFromCtx(r.Context())
			ok, err := c.Enforce(r.Context(), claims, r.URL.Path, r.Method)
			if err != nil {
				response.FailCtx(r.Context(), w, err)
				return
			}
			if !ok {
				response.FailCtx(r.Context(), w, xerr.Forbidden(i18n.Forbidden))
				return
			}
			next(w, r)
		}
	}
}

func stripBearer(v string) string {
	v = strings.TrimSpace(v)
	if len(v) > 7 && strings.EqualFold(v[:7], "Bearer ") {
		return strings.TrimSpace(v[7:])
	}
	return v
}

// IsAdminWrite 判断是否为管理员写权限
func IsAdminWrite(method, path string) bool {
	switch strings.ToUpper(method) {
	case http.MethodGet, http.MethodHead, http.MethodOptions:
		return false
	}
	if !strings.EqualFold(method, http.MethodPost) {
		return true
	}
	if strings.HasSuffix(path, "/list") {
		return false
	}
	if strings.HasSuffix(path, "/authority/menu/role") || strings.HasSuffix(path, "/authority/api/role") {
		return false
	}
	return true
}

// previewWriteDenied 预览模式下是否允许写操作
func previewWriteDenied(method, path string) bool {
	// 预览模式下允许更新菜单
	if strings.HasSuffix(path, "/menu/update") {
		return false
	}

	return IsAdminWrite(method, path)
}
