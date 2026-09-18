package utils

import (
	"context"

	"github.com/zeromicro/go-zero/core/trace"
	otelTrace "go.opentelemetry.io/otel/trace"
)

func CloneCtxWithTraceID(ctx context.Context) context.Context {
	traceID := trace.TraceIDFromContext(ctx)
	spanID := trace.SpanIDFromContext(ctx)

	// 创建一个独立的 Background Context
	bgCtx := context.Background()

	// 【关键】将 TraceID 和 SpanID 重新注入到 bgCtx 中
	// 这需要借助 OpenTelemetry 的 API 来构造 SpanContext
	tid, _ := otelTrace.TraceIDFromHex(traceID)
	sid, _ := otelTrace.SpanIDFromHex(spanID)

	sc := otelTrace.NewSpanContext(otelTrace.SpanContextConfig{
		TraceID:    tid,
		SpanID:     sid,
		TraceFlags: otelTrace.FlagsSampled,
		Remote:     true, // 标记为远程传入，避免采样决策冲突
	})

	// 将 SpanContext 注入到 Context 中
	newCtx := otelTrace.ContextWithRemoteSpanContext(bgCtx, sc)

	return newCtx
}
