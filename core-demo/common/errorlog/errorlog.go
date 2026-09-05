package errorlog

import (
	"context"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/xerr"
)

const (
	MaxStack = 16 * 1024
	MaxBody  = 16 * 1024
)

type Record struct {
	UserID         int64
	RequestMethod  string
	RequestPath    string
	RequestQuery   string
	RequestBody    string
	ServiceName    string
	ResponseStatus int
	ResponseBody   string
	Subject        string
	Detail         string
	DurationMS     int
	ClientIP       string
	UserAgent      string
}

// Recorder 记录器接口
type Recorder interface {
	RecordError(ctx context.Context, rec Record)
}

// bag 错误日志上下文
type bag struct {
	Subject string
	Detail  string
}

// bagKey 错误日志上下文键
type bagKey struct{}

// WithBag 创建错误日志上下文
func WithBag(ctx context.Context) (context.Context, *bag) {
	b := &bag{}
	return context.WithValue(ctx, bagKey{}, b), b
}

// BagFromCtx 从上下文中获取错误日志上下文
func BagFromCtx(ctx context.Context) *bag {
	b, _ := ctx.Value(bagKey{}).(*bag)
	return b
}

// Note 记录错误日志
func Note(ctx context.Context, subject, detail string) {
	b := BagFromCtx(ctx)
	if b == nil {
		return
	}
	if subject != "" && b.Subject == "" {
		b.Subject = subject
	}
	if detail != "" && b.Detail == "" {
		b.Detail = xerr.ClipStack(detail)
	}
}

// ShouldCollect 是否收集错误日志
func ShouldCollect(status int) bool {
	return status >= 500
}

// Report 报告错误日志
func Report(parent context.Context, rec Recorder, record Record) {
	if rec == nil {
		return
	}
	claims := ctxdata.ClaimsFromCtx(parent)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		if claims != nil {
			ctx = ctxdata.WithClaims(ctx, claims)
		}
		rec.RecordError(ctx, record)
	}()
}
