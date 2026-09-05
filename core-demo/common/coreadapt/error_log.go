package coreadapt

import (
	"context"

	"oa.98ent.com/p9/core/common/errorlog"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

// errorRecorder 错误记录器
type errorRecorder struct {
	cli coreclient.Core
}

// ErrorRecorder 创建错误记录器
func ErrorRecorder(cli coreclient.Core) errorlog.Recorder {
	return errorRecorder{cli: cli}
}

// RecordError 记录错误
func (e errorRecorder) RecordError(ctx context.Context, rec errorlog.Record) {
	req := &coreclient.CreateErrorLogReq{
		RequestMethod:  rec.RequestMethod,
		RequestPath:    rec.RequestPath,
		ServiceName:    rec.ServiceName,
		ResponseStatus: int32(rec.ResponseStatus),
		DurationMs:     int32(rec.DurationMS),
		ClientIp:       rec.ClientIP,
		UserId:         i64Ptr(rec.UserID),
		RequestQuery:   strPtr(rec.RequestQuery),
		RequestBody:    strPtr(rec.RequestBody),
		ResponseBody:   strPtr(rec.ResponseBody),
		Subject:        strPtr(rec.Subject),
		Detail:         strPtr(rec.Detail),
		UserAgent:      strPtr(rec.UserAgent),
	}
	if _, err := e.cli.CreateErrorLog(ctx, req); err != nil {
		logx.Errorf("create error log: %v", err)
	}
}
