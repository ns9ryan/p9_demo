package grpc_error

import (
	"fmt"
	"net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ErrorResponse API错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// ConvertGrpcError 将gRPC错误转换为HTTP错误
func ConvertGrpcError(err error) (int, string) {
	if err == nil {
		return 1, "success"
	}

	st, ok := status.FromError(err)
	if !ok {
		// 不是gRPC错误，返回通用错误
		return 500, fmt.Sprintf("internal error: %v", err)
	}

	switch st.Code() {
	case codes.OK:
		return 1, st.Message()
	case codes.Canceled:
		return 408, "request canceled"
	case codes.Unknown:
		return 500, "unknown error"
	case codes.InvalidArgument:
		return 400, st.Message()
	case codes.DeadlineExceeded:
		return 408, "request timeout"
	case codes.NotFound:
		return 404, st.Message()
	case codes.AlreadyExists:
		return 409, st.Message()
	case codes.PermissionDenied:
		return 403, st.Message()
	case codes.ResourceExhausted:
		return 429, "resource exhausted"
	case codes.FailedPrecondition:
		return 400, st.Message()
	case codes.Aborted:
		return 409, st.Message()
	case codes.OutOfRange:
		return 400, st.Message()
	case codes.Unimplemented:
		return 501, "method not implemented"
	case codes.Internal:
		return 500, st.Message()
	case codes.Unavailable:
		return 503, "service unavailable"
	case codes.DataLoss:
		return 500, "data loss"
	case codes.Unauthenticated:
		return 401, st.Message()
	default:
		return 500, st.Message()
	}
}

// GetHTTPStatusCode 获取HTTP状态码对应的gRPC code
func GetHTTPStatusCode(code int) int {
	switch code {
	case 400:
		return http.StatusBadRequest
	case 401:
		return http.StatusUnauthorized
	case 403:
		return http.StatusForbidden
	case 404:
		return http.StatusNotFound
	case 408:
		return http.StatusRequestTimeout
	case 409:
		return http.StatusConflict
	case 429:
		return http.StatusTooManyRequests
	case 500:
		return http.StatusInternalServerError
	case 501:
		return http.StatusNotImplemented
	case 503:
		return http.StatusServiceUnavailable
	case 1:
		return http.StatusOK
	default:
		return http.StatusInternalServerError
	}
}
