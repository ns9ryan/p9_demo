package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"oa.98ent.com/p9/common/errorlog"
	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/common/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// SetupHTTPX 设置HTTPX的OK和Error处理函数
func SetupHTTPX(trans *i18n.Translator, code string, isDebug bool) {
	// 成功处理函数
	httpx.SetOkHandler(func(ctx context.Context, data any) any {
		return map[string]any{"code": 0, "msg": trans.T(ctx, "common.success"), "data": data}
	})
	// 错误处理函数
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		// 转换为xerr.Error
		e := FromError(err)
		// 记录内部错误
		noteInternal(ctx, e)
		// 基础响应结构
		result := map[string]any{"code": e.Status, "msg": localize(ctx, trans, code, e)}
		// 是否返回debug信息
		if isDebug {
			result["debug"] = map[string]any{
				"error": err.Error(),     // 内部错误信息
				"stack": xerr.StackOf(e), // 错误堆栈
				"cause": xerr.Subject(e), // 错误原因
			}
		}
		return e.Status, result
	})
}

// OkCtx 返回成功响应带上下文
func OkCtx(ctx context.Context, w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "msg": "ok", "data": data})
}

// FailCtx 返回失败响应带上下文
func FailCtx(ctx context.Context, w http.ResponseWriter, err error) {
	httpx.ErrorCtx(ctx, w, err)
}

func localize(ctx context.Context, trans *i18n.Translator, code string, e *xerr.Error) string {
	if e == nil || trans == nil {
		return trans.T(ctx, "common.internalError")
	}
	if len(e.Params) > 0 {
		return trans.Tf(ctx, e.Message, e.Params)
	}
	return i18n.TG(ctx, code, i18n.GroupError, trans.T(ctx, e.Message))
}

// FromError 从错误转换为xerr.Error
func FromError(err error) *xerr.Error {
	if err == nil {
		return xerr.InternalServerError("common.internalError")
	}
	var e *xerr.Error
	if errors.As(err, &e) {
		return e
	}
	if re := asRequestError(err); re != nil {
		return re
	}
	if st, ok := status.FromError(err); ok {
		e := xerr.Err(grpcToHTTP(st.Code()), st.Message())
		if sub, stack := xerr.FromStatus(st); sub != "" || stack != "" {
			if sub != "" {
				e.Cause = errors.New(sub)
			}
			e.Stack = stack
		}
		return e
	}
	return xerr.AsError(err)
}

// grpcToHTTP 将gRPC的错误码转换为HTTP状态码
func grpcToHTTP(code codes.Code) int {
	switch code {
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.NotFound:
		return http.StatusNotFound
	default:
		if int(code) >= 400 && int(code) < 600 {
			return int(code)
		}
		return http.StatusInternalServerError
	}
}

// noteInternal 记录内部错误
func noteInternal(ctx context.Context, e *xerr.Error) {
	if e == nil || !errorlog.ShouldCollect(e.Status) {
		return
	}
	errorlog.Note(ctx, xerr.Subject(e), e.Stack)
}
