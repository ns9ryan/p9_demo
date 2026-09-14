package errorhandler

import (
	"context"
	"errors"
	"net/http"

	"oa.98ent.com/p9/platform-operator/pkg/api/response"
	"oa.98ent.com/p9/platform-operator/pkg/api/rpcerror"
	"oa.98ent.com/p9/platform-operator/pkg/api/validate"
	"oa.98ent.com/p9/platform-operator/pkg/i18n"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Handler API错误处理器
type Handler struct {
	trans *i18n.Translator // 翻译器
	debug bool             // 是否开启调试信息
}

// New 创建API错误处理器
func New(trans *i18n.Translator, debug bool) *Handler {
	return &Handler{
		trans: trans,
		debug: debug,
	}
}

// Handle 处理API错误
func (h *Handler) Handle(ctx context.Context, err error) (int, any) {
	// 参数校验错误已经由Validator完成多语言翻译
	var validationError *validate.Error
	if errors.As(err, &validationError) {
		return h.response(
			http.StatusBadRequest,
			validationError.Error(),
			h.debugInfo("api", "", err),
		)
	}

	source := "api"
	rpcMethod := ""

	// 判断错误是否来自RPC调用，并获取具体RPC方法
	if rpcErr, ok := rpcerror.FromError(err); ok {
		source = "rpc"
		rpcMethod = rpcErr.Method
	}

	// gRPC错误按照状态码转换并翻译
	if grpcStatus, ok := status.FromError(err); ok {
		return h.handleGRPC(
			ctx,
			err,
			grpcStatus,
			source,
			rpcMethod,
		)
	}

	// 普通错误主要来自JSON、path、form等请求解析失败
	return h.response(
		http.StatusBadRequest,
		h.trans.Trans(ctx, i18nkey.InvalidRequest),
		h.debugInfo("api", "", err),
	)
}

// handleGRPC 处理gRPC错误
func (h *Handler) handleGRPC(
	ctx context.Context,
	err error,
	grpcStatus *status.Status,
	source string,
	rpcMethod string,
) (int, any) {
	// 将gRPC状态码转换为对应的HTTP状态码
	httpStatus := grpcCodeToHTTPStatus(grpcStatus.Code())

	// 获取最终返回给前端的错误提示
	message := h.grpcMessage(ctx, grpcStatus)

	// HTTP 500及以上属于服务端错误，需要记录内部错误日志
	if httpStatus >= http.StatusInternalServerError {
		h.logInternal(ctx, source, rpcMethod, err)
	}

	return h.response(
		httpStatus,
		message,
		h.debugInfo(source, rpcMethod, err),
	)
}

// grpcMessage 获取gRPC错误提示
func (h *Handler) grpcMessage(
	ctx context.Context,
	grpcStatus *status.Status,
) string {
	switch grpcStatus.Code() {
	// 请求过多或服务资源配额耗尽
	case codes.ResourceExhausted:
		return h.trans.Trans(ctx, i18nkey.TooManyRequests)

	// RPC服务不可用，例如服务未启动或连接失败
	case codes.Unavailable:
		return h.trans.Trans(ctx, i18nkey.ServiceUnavailable)

	// RPC调用超过截止时间
	case codes.DeadlineExceeded:
		return h.trans.Trans(ctx, i18nkey.RequestTimeout)

	// RPC请求被取消
	case codes.Canceled:
		return h.trans.Trans(ctx, i18nkey.RequestTimeout)

	// RPC方法不存在或尚未实现
	case codes.Unimplemented:
		return h.trans.Trans(ctx, i18nkey.InternalError)

	// 未明确分类的gRPC错误
	case codes.Unknown:
		return h.trans.Trans(ctx, i18nkey.InternalError)
	}

	// RPC业务错误使用Message中的i18n Key进行翻译
	message := h.trans.Trans(ctx, grpcStatus.Message())

	// Internal错误翻译不到时，不向前端暴露原始错误信息
	if grpcStatus.Code() == codes.Internal &&
		message == grpcStatus.Message() {
		return h.trans.Trans(ctx, i18nkey.InternalError)
	}

	return message
}

// debugInfo 创建调试信息
func (h *Handler) debugInfo(
	source string,
	rpcMethod string,
	err error,
) *response.DebugInfo {
	if !h.debug {
		return nil
	}

	return &response.DebugInfo{
		Source: source,
		Rpc:    rpcMethod,
		Error:  err.Error(),
	}
}

// logInternal 记录内部错误
func (h *Handler) logInternal(
	ctx context.Context,
	source string,
	rpcMethod string,
	err error,
) {
	logger := logx.WithContext(ctx)

	// RPC错误额外记录具体RPC方法
	if rpcMethod != "" {
		logger.Errorw(
			"RPC调用失败",
			logx.Field("rpc", rpcMethod),
			logx.Field("error", err.Error()),
		)
		return
	}

	logger.Errorw(
		"接口处理失败",
		logx.Field("source", source),
		logx.Field("error", err.Error()),
	)
}

// response 创建统一错误响应
func (h *Handler) response(
	code int,
	message string,
	debug *response.DebugInfo,
) (int, any) {
	return code, &response.ErrorResponse{
		Code:  code,
		Msg:   message,
		Debug: debug,
	}
}

// grpcCodeToHTTPStatus 将gRPC状态码转换为HTTP状态码
func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	// 参数错误
	case codes.InvalidArgument,
		codes.FailedPrecondition,
		codes.OutOfRange:
		return http.StatusBadRequest

	// 未认证，例如未登录或Token无效
	case codes.Unauthenticated:
		return http.StatusUnauthorized

	// 已认证，但没有操作权限
	case codes.PermissionDenied:
		return http.StatusForbidden

	// 请求的资源不存在
	case codes.NotFound:
		return http.StatusNotFound

	// 请求创建的资源已经存在
	case codes.AlreadyExists:
		return http.StatusConflict

	// 请求过多或资源配额耗尽
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests

	// 请求被取消
	case codes.Canceled:
		return http.StatusRequestTimeout

	// 操作被中止，例如并发冲突
	case codes.Aborted:
		return http.StatusConflict

	// RPC方法尚未实现
	case codes.Unimplemented:
		return http.StatusNotImplemented

	// RPC服务当前不可用
	case codes.Unavailable:
		return http.StatusServiceUnavailable

	// RPC调用超过截止时间
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout

	// 其他错误统一按服务器内部错误处理
	default:
		return http.StatusInternalServerError
	}
}
