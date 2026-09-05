package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"oa.98ent.com/p9/core/common/errorlog"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"

	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Ok 返回成功响应
func Ok(w http.ResponseWriter, data any) {
	OkCtx(context.Background(), w, data)
}

// OkCtx 返回成功响应带上下文
func OkCtx(ctx context.Context, w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	_ = json.NewEncoder(w).Encode(map[string]any{"code": 0, "msg": i18n.T(ctx, i18n.Success), "data": data})
}

// Fail 返回失败响应
func Fail(w http.ResponseWriter, err error) {
	FailCtx(context.Background(), w, err)
}

// FailCtx 返回失败响应带上下文
func FailCtx(ctx context.Context, w http.ResponseWriter, err error) {
	e := FromError(err)
	noteInternal(ctx, e)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(e.Status)
	_ = json.NewEncoder(w).Encode(map[string]any{"code": e.Status, "msg": localize(ctx, e)})
}

func localize(ctx context.Context, e *xerr.Error) string {
	if e == nil {
		return i18n.T(ctx, i18n.InternalError)
	}
	if len(e.Params) > 0 {
		return i18n.Tf(ctx, e.Message, e.Params)
	}
	return i18n.T(ctx, e.Message)
}

// FromError 从错误转换为xerr.Error
func FromError(err error) *xerr.Error {
	if err == nil {
		return xerr.InternalServerError(i18n.InternalError)
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

// SetupHTTPX 设置HTTPX的OK和Error处理函数
func SetupHTTPX() {
	httpx.SetOkHandler(func(ctx context.Context, data any) any {
		return map[string]any{"code": 0, "msg": i18n.T(ctx, i18n.Success), "data": data}
	})
	httpx.SetErrorHandlerCtx(func(ctx context.Context, err error) (int, any) {
		e := FromError(err)
		noteInternal(ctx, e)
		return e.Status, map[string]any{"code": e.Status, "msg": localize(ctx, e)}
	})
}

// noteInternal 记录内部错误
func noteInternal(ctx context.Context, e *xerr.Error) {
	if e == nil || !errorlog.ShouldCollect(e.Status) {
		return
	}
	errorlog.Note(ctx, xerr.Subject(e), e.Stack)
}
