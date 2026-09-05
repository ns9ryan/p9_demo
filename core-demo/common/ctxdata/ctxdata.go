package ctxdata

import (
	"context"
	"strconv"
	"strings"

	"google.golang.org/grpc/metadata"
)

type ctxKey int

const (
	claimsKey ctxKey = iota + 1
	rawTokenKey
)

type skipTenantKey struct{}
type skipSoftDeleteKey struct{}

type claimsBox struct {
	c *Claims
}

const (
	headerUserID       = "x-user-id"
	headerUserCode     = "x-user-code"
	headerUsername     = "x-username"
	headerOperatorID   = "x-operator-id"
	headerOperatorCode = "x-operator-code"
	headerRoleCodes    = "x-role-codes"
	headerSalt         = "x-salt"
	headerExpiresAt    = "x-expires-at"
	headerRawToken     = "x-raw-token"
	headerIsPlatform   = "x-is-platform"
	headerTokenType    = "x-token-type"
)

// Claims 用户身份信息
type Claims struct {
	UserID       int64
	UserCode     string
	Username     string
	OperatorID   int64
	OperatorCode string
	RoleCodes    []string
	Salt         string
	ExpiresAt    int64
	IsPlatform   bool
	TokenType    string
}

// WithClaimsHolder 在外层放入可写的身份盒子，内层 WithClaims 会回写，外层仍能读到。
func WithClaimsHolder(ctx context.Context) context.Context {
	if boxFrom(ctx) != nil {
		return ctx
	}
	return context.WithValue(ctx, claimsKey, &claimsBox{})
}

// WithClaims 添加用户身份信息到上下文
func WithClaims(ctx context.Context, c *Claims) context.Context {
	if c == nil {
		return ctx
	}
	if box := boxFrom(ctx); box != nil {
		box.c = c
	} else {
		ctx = context.WithValue(ctx, claimsKey, &claimsBox{c: c})
	}
	return metadata.AppendToOutgoingContext(ctx,
		headerUserID, strconv.FormatInt(c.UserID, 10),
		headerUserCode, c.UserCode,
		headerUsername, c.Username,
		headerOperatorID, strconv.FormatInt(c.OperatorID, 10),
		headerOperatorCode, c.OperatorCode,
		headerRoleCodes, strings.Join(c.RoleCodes, ","),
		headerSalt, c.Salt,
		headerExpiresAt, strconv.FormatInt(c.ExpiresAt, 10),
		headerIsPlatform, formatBool(c.IsPlatform),
		headerTokenType, c.TokenType,
	)
}

// ClaimsFromCtx 从上下文中获取用户身份信息
func ClaimsFromCtx(ctx context.Context) *Claims {
	if box := boxFrom(ctx); box != nil && box.c != nil {
		return box.c
	}
	return claimsFromIncomingMD(ctx)
}

func boxFrom(ctx context.Context) *claimsBox {
	box, _ := ctx.Value(claimsKey).(*claimsBox)
	return box
}

// OperatorIDFromCtx 从上下文中获取操作员ID
func OperatorIDFromCtx(ctx context.Context) int64 {
	c := ClaimsFromCtx(ctx)
	if c == nil {
		return 0
	}
	return c.OperatorID
}

// OperatorCodeFromCtx 从上下文中获取操作员代码
func OperatorCodeFromCtx(ctx context.Context) string {
	c := ClaimsFromCtx(ctx)
	if c == nil {
		return ""
	}
	return c.OperatorCode
}

// WithRawToken 添加原始令牌到上下文
func WithRawToken(ctx context.Context, token string) context.Context {
	if token == "" {
		return ctx
	}
	ctx = context.WithValue(ctx, rawTokenKey, token)
	return metadata.AppendToOutgoingContext(ctx, headerRawToken, token)
}

// RawTokenFromCtx 从上下文中获取原始令牌
func RawTokenFromCtx(ctx context.Context) string {
	if s, _ := ctx.Value(rawTokenKey).(string); s != "" {
		return s
	}
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}
	return first(md, headerRawToken)
}

// SkipTenant 跳过租户
func SkipTenant(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipTenantKey{}, true)
}

// TenantSkipped 是否跳过租户
func TenantSkipped(ctx context.Context) bool {
	ok, _ := ctx.Value(skipTenantKey{}).(bool)
	return ok
}

// SkipSoftDelete 跳过软删除
func SkipSoftDelete(ctx context.Context) context.Context {
	return context.WithValue(ctx, skipSoftDeleteKey{}, true)
}

// SoftDeleteSkipped 是否跳过软删除
func SoftDeleteSkipped(ctx context.Context) bool {
	ok, _ := ctx.Value(skipSoftDeleteKey{}).(bool)
	return ok
}

// claimsFromIncomingMD 从入站元数据中获取用户身份信息
func claimsFromIncomingMD(ctx context.Context) *Claims {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	userID := parseInt(first(md, headerUserID))
	username := first(md, headerUsername)
	operatorID := parseInt(first(md, headerOperatorID))
	isPlatform := parseBool(first(md, headerIsPlatform))
	if userID == 0 && username == "" && operatorID == 0 && !isPlatform {
		return nil
	}
	codes := []string{}
	if raw := first(md, headerRoleCodes); raw != "" {
		codes = strings.Split(raw, ",")
	}
	return &Claims{
		UserID:       userID,
		UserCode:     first(md, headerUserCode),
		Username:     username,
		OperatorID:   operatorID,
		OperatorCode: first(md, headerOperatorCode),
		RoleCodes:    codes,
		Salt:         first(md, headerSalt),
		ExpiresAt:    parseInt(first(md, headerExpiresAt)),
		IsPlatform:   isPlatform,
		TokenType:    first(md, headerTokenType),
	}
}

// first 获取第一个值
func first(md metadata.MD, key string) string {
	vals := md.Get(key)
	if len(vals) == 0 {
		return ""
	}
	return vals[0]
}

// parseInt 解析整数
func parseInt(s string) int64 {
	if s == "" {
		return 0
	}
	n, _ := strconv.ParseInt(s, 10, 64)
	return n
}

// formatBool 格式化布尔值
func formatBool(v bool) string {
	if v {
		return "1"
	}
	return "0"
}

// parseBool 解析布尔值
func parseBool(s string) bool {
	return s == "1" || strings.EqualFold(s, "true")
}
