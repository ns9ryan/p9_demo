package tracing

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"oa.98ent.com/p9/common/xerr"
)

// Error 把错误写入当前 span。Jaeger 中该 span 会标红，Logs 的 exception 事件带堆栈。
func Error(ctx context.Context, err error) {
	if err == nil {
		return
	}
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	span.RecordError(err, trace.WithStackTrace(true))
	span.SetStatus(codes.Error, err.Error())
	stack := xerr.StackOf(err)
	if stack == "" {
		_, stack = xerr.FromGRPC(err)
	}
	if stack != "" {
		span.SetAttributes(attribute.String("exception.stacktrace", stack))
		// Debug(ctx, "exception", String("stacktrace", stack))
	}
}

// HTTPError 把 HTTP 5xx / panic 写入当前 span。业务 handler 不必再逐个调用 Error。
func HTTPError(ctx context.Context, status int, subject, stack string) {
	if status < http.StatusInternalServerError {
		return
	}
	if subject == "" {
		subject = http.StatusText(status)
	}
	Error(ctx, &xerr.Error{Status: status, Message: subject, Stack: stack})
}

// UnaryServerInterceptor RPC handler 返回 error 时自动写入当前 span。
// 业务侧只需在需要额外 Tags 时再调 Error / Attrs。
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			Error(ctx, err)
		}
		return resp, err
	}
}

// Debug 在当前 span 记录一条调试事件。Jaeger 中显示在 Logs。
func Debug(ctx context.Context, msg string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	if len(attrs) == 0 {
		span.AddEvent(msg)
		return
	}
	span.AddEvent(msg, trace.WithAttributes(attrs...))
}

// Attrs 给当前 span 打 Tags。
func Attrs(ctx context.Context, attrs ...attribute.KeyValue) {
	if len(attrs) == 0 {
		return
	}
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}
	span.SetAttributes(attrs...)
}

// String 构造字符串属性
func String(key, value string) attribute.KeyValue {
	return attribute.String(key, value)
}

// Int 构造整型属性
func Int(key string, value int) attribute.KeyValue {
	return attribute.Int(key, value)
}

// Int64 构造 int64 属性
func Int64(key string, value int64) attribute.KeyValue {
	return attribute.Int64(key, value)
}

// Bool 构造布尔属性
func Bool(key string, value bool) attribute.KeyValue {
	return attribute.Bool(key, value)
}
