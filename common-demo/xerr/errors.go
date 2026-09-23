package xerr

import (
	"errors"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// maxStack 最大堆栈大小
const maxStack = 128 * 1024 // 128KB

// Error 错误
type Error struct {
	Status  int
	Message string
	Params  map[string]any
	Cause   error
	Stack   string
}

// Error 返回错误消息
func (e *Error) Error() string { return e.Message }

// Unwrap 返回错误原因
func (e *Error) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.Cause
}

// StatusTokenExpired 令牌过期状态码
const StatusTokenExpired = 498

// Err 创建错误
func Err(status int, msg string) *Error {
	return &Error{Status: status, Message: msg}
}

// ErrWith 创建带有参数的错误
func ErrWith(status int, msg string, params map[string]any) *Error {
	return &Error{Status: status, Message: msg, Params: params}
}

// BadRequest 创建400 Bad Request错误
func BadRequest(msg string) *Error { return Err(http.StatusBadRequest, msg) }

// BadRequestWith 创建带有参数的400 Bad Request错误
func BadRequestWith(msg string, params map[string]any) *Error {
	return ErrWith(http.StatusBadRequest, msg, params)
}

// 创建401 Unauthorized错误
func Unauthorized(msg string) *Error { return Err(http.StatusUnauthorized, msg) }

// 创建403 Forbidden错误
func Forbidden(msg string) *Error { return Err(http.StatusForbidden, msg) }

// 创建404 Not Found错误
func NotFound(msg string) *Error { return Err(http.StatusNotFound, msg) }

// 创建408 Token Expired错误
func TokenExpired(msg string) *Error { return Err(StatusTokenExpired, msg) }

// 创建409 Conflict错误
func Conflict(msg string) *Error { return Err(http.StatusConflict, msg) }

// 创建500 Internal Server Error错误
func InternalServerError(msg string) *Error { return Err(http.StatusInternalServerError, msg) }

// EntNotFound 创建 Ent 数据不存在错误
func EntNotFound(msg string, err error) *Error {
	return &Error{
		Status:  http.StatusNotFound,
		Message: msg,
		Cause:   err,
		Stack:   ClipStack(string(debug.Stack())),
	}
}

// EntConstraintError 创建 Ent 数据约束冲突错误
func EntConstraintError(msg string, err error) *Error {
	return &Error{
		Status:  http.StatusConflict,
		Message: msg,
		Cause:   err,
		Stack:   ClipStack(string(debug.Stack())),
	}
}

// EntValidationError 创建 Ent 数据校验失败错误
func EntValidationError(msg string, err error) *Error {
	return &Error{
		Status:  http.StatusBadRequest,
		Message: msg,
		Cause:   err,
		Stack:   ClipStack(string(debug.Stack())),
	}
}

// EntInternalServerError 创建 Ent 内部服务器错误
func EntInternalServerError(msg string, err error) *Error {
	return &Error{
		Status:  http.StatusInternalServerError,
		Message: msg,
		Cause:   err,
		Stack:   ClipStack(string(debug.Stack())),
	}
}

// AsError 将错误转换为xerr.Error
func AsError(err error) *Error {
	if err == nil {
		return Err(http.StatusInternalServerError, "common.internalError")
	}
	// 如果err已经是xerr.Error类型，则直接返回
	if ve, ok := errors.AsType[*Error](err); ok {
		return ve
	}

	logx.Errorf("internal error: %v", err)
	return &Error{
		Status:  http.StatusInternalServerError,
		Message: "common.internalError",
		Cause:   err,
		Stack:   ClipStack(string(debug.Stack())),
	}
}

// Subject 返回错误主题
func Subject(err error) string {
	if err == nil {
		return ""
	}
	var e *Error
	if errors.As(err, &e) && e != nil {
		if e.Cause != nil {
			return e.Cause.Error()
		}
		if e.Message != "" {
			return e.Message
		}
	}
	return err.Error()
}

// StackOf 返回错误堆栈
func StackOf(err error) string {
	var e *Error
	if errors.As(err, &e) && e != nil {
		return e.Stack
	}
	return ""
}

// ClipStack 裁剪堆栈
func ClipStack(s string) string {
	if len(s) <= maxStack {
		return s
	}
	return s[:maxStack]
}

// FromGRPC 从GRPC错误中获取主题和堆栈
func FromGRPC(err error) (subject, stack string) {
	st, ok := status.FromError(err)
	if !ok || st == nil {
		return "", ""
	}
	return FromStatus(st)
}

// FromStatus 从状态中获取主题和堆栈
func FromStatus(st *status.Status) (subject, stack string) {
	if st == nil {
		return "", ""
	}
	for _, d := range st.Details() {
		di, ok := d.(*errdetails.DebugInfo)
		if !ok || di == nil {
			continue
		}
		return di.GetDetail(), ClipStack(strings.Join(di.GetStackEntries(), "\n"))
	}
	return "", ""
}

// AttachGRPC 附加GRPC错误
func AttachGRPC(st *status.Status, e *Error) *status.Status {
	if st == nil || e == nil || e.Status < 500 {
		return st
	}
	di := &errdetails.DebugInfo{}
	if e.Cause != nil {
		di.Detail = e.Cause.Error()
	} else if e.Message != "" {
		di.Detail = e.Message
	}
	if e.Stack != "" {
		di.StackEntries = strings.Split(strings.TrimRight(e.Stack, "\n"), "\n")
	}
	if di.Detail == "" && len(di.StackEntries) == 0 {
		return st
	}
	with, err := st.WithDetails(di)
	if err != nil {
		return st
	}
	return with
}

// RpcErr 转换为GRPC错误
func RpcErr(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := status.FromError(err); ok {
		return err
	}

	e := AsError(err)
	var code codes.Code
	switch e.Status {
	case http.StatusBadRequest:
		code = codes.InvalidArgument
	case http.StatusConflict:
		code = codes.AlreadyExists
	case http.StatusUnauthorized:
		code = codes.Unauthenticated
	case http.StatusForbidden:
		code = codes.PermissionDenied
	case http.StatusNotFound:
		code = codes.NotFound
	default:
		if e.Status >= 400 && e.Status < 600 {
			code = codes.Code(e.Status)
		} else {
			code = codes.Internal
		}
	}
	st := status.New(code, e.Message)
	return AttachGRPC(st, e).Err()
}
