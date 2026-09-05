package rpcerror

import (
	"errors"

	"google.golang.org/grpc/status"
)

// Error RPC调用错误
type Error struct {
	Method string // RPC方法
	err    error  // 原始错误
}

// Wrap 包装RPC调用错误
func Wrap(method string, err error) error {
	if err == nil {
		return nil
	}

	return &Error{
		Method: method,
		err:    err,
	}
}

// Error 返回错误信息
func (e *Error) Error() string {
	return e.err.Error()
}

// Unwrap 返回原始错误 - Go 错误链
func (e *Error) Unwrap() error {
	return e.err
}

// GRPCStatus 返回gRPC状态 - gRPC 状态
func (e *Error) GRPCStatus() *status.Status {
	return status.Convert(e.err)
}

// FromError 获取RPC调用错误
func FromError(err error) (*Error, bool) {
	var rpcErr *Error
	if !errors.As(err, &rpcErr) {
		return nil, false
	}

	return rpcErr, true
}
