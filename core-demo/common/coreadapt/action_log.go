package coreadapt

import (
	"context"

	"oa.98ent.com/p9/core/common/middleware"
	"oa.98ent.com/p9/core/rpc/coreclient"

	"github.com/zeromicro/go-zero/core/logx"
)

// actionRecorder 动作记录器
type actionRecorder struct {
	cli coreclient.Core
}

// ActionRecorder 创建动作记录器
func ActionRecorder(cli coreclient.Core) middleware.ActionRecorder {
	return actionRecorder{cli: cli}
}

// RecordAction 记录动作
func (a actionRecorder) RecordAction(ctx context.Context, rec middleware.ActionRecord) {
	req := &coreclient.CreateAdminActionLogReq{
		UserId:         rec.UserID,
		RequestMethod:  rec.RequestMethod,
		RequestPath:    rec.RequestPath,
		ActionResult:   int32(rec.ActionResult),
		ResponseStatus: int32(rec.ResponseStatus),
		DurationMs:     int32(rec.DurationMS),
		ClientIp:       rec.ClientIP,
		RequestQuery:   strPtr(rec.RequestQuery),
		RequestBody:    strPtr(rec.RequestBody),
		ResponseBody:   strPtr(rec.ResponseBody),
		UserAgent:      strPtr(rec.UserAgent),
	}
	if _, err := a.cli.CreateAdminActionLog(ctx, req); err != nil {
		logx.Errorf("create admin action log: %v", err)
	}
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func i64Ptr(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}
