package grpcerror

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Internal 创建内部错误
func Internal(message string) error {
	// Internal = 13 服务端内部错误
	return status.Error(codes.Internal, message)
}

// InvalidArgument 创建参数错误
func InvalidArgument(message string) error {
	// InvalidArgument = 3 请求参数无效
	return status.Error(codes.InvalidArgument, message)
}

// NotFound 创建资源不存在错误
func NotFound(message string) error {
	// NotFound = 5 请求的资源不存在
	return status.Error(codes.NotFound, message)
}
