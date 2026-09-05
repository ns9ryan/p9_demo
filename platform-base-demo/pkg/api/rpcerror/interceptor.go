package rpcerror

import (
	"context"
	"strings"

	"google.golang.org/grpc"
)

// UnaryClientInterceptor RPC客户端错误拦截器
func UnaryClientInterceptor(
	ctx context.Context,
	method string,
	req any,
	reply any,
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// 调用RPC方法
	err := invoker(ctx, method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}

	// 记录RPC方法并保留原始错误
	return Wrap(strings.TrimPrefix(method, "/"), err)
}
