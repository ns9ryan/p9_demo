package tracing

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"google.golang.org/grpc"
	"oa.98ent.com/p9/core/common/xerr"
)

func TestErrorNilDoesNotPanic(t *testing.T) {
	Error(context.Background(), nil)
}

func TestErrorAndDebugOnBackground(t *testing.T) {
	ctx := context.Background()
	Error(ctx, errors.New("boom"))
	Error(ctx, &xerr.Error{Status: 500, Message: "fail", Stack: "goroutine 1 [running]:\nmain.main()"})
	Debug(ctx, "debug msg", String("k", "v"), Int("n", 1))
	Attrs(ctx, Bool("ok", true), Int64("id", 2))
	HTTPError(ctx, http.StatusInternalServerError, "fail", "goroutine 1 [running]:\nmain.main()")
	HTTPError(ctx, http.StatusBadRequest, "skip", "")
}

func TestUnaryServerInterceptorRecordsError(t *testing.T) {
	boom := errors.New("boom")
	interceptor := UnaryServerInterceptor()
	resp, err := interceptor(context.Background(), nil, &grpc.UnaryServerInfo{FullMethod: "/x.Y/Z"}, func(context.Context, any) (any, error) {
		return nil, boom
	})
	if resp != nil || err != boom {
		t.Fatalf("resp=%v err=%v", resp, err)
	}
}

func TestUnaryServerInterceptorPassesSuccess(t *testing.T) {
	interceptor := UnaryServerInterceptor()
	resp, err := interceptor(context.Background(), "req", &grpc.UnaryServerInfo{FullMethod: "/x.Y/Z"}, func(context.Context, any) (any, error) {
		return "ok", nil
	})
	if err != nil || resp != "ok" {
		t.Fatalf("resp=%v err=%v", resp, err)
	}
}
