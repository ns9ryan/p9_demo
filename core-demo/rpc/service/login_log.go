package service

import (
	"context"
	"strings"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/loginlog"
	"oa.98ent.com/p9/core/rpc/model"
)

type LoginLogListReq struct {
	PageReq
	Username    string
	LoginResult int16
	UserID      int64
	LoginAtFrom int64
	LoginAtTo   int64
}

func (d *Deps) ListLoginLogs(ctx context.Context, claims *ctxdata.Claims, req LoginLogListReq) ([]model.LoginLog, int64, error) {
	if claims == nil {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	if d.Mode == ModeOn && claims.OperatorID == 0 {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	ctx = ctxdata.WithClaims(ctx, claims)
	q := d.Client.LoginLog.Query()
	if s := strings.TrimSpace(req.Username); s != "" {
		q.Where(loginlog.UsernameContains(s))
	}
	if req.LoginResult != 0 {
		q.Where(loginlog.LoginResultEQ(req.LoginResult))
	}
	if req.UserID != 0 {
		q.Where(loginlog.UserIDEQ(req.UserID))
	}
	if req.LoginAtFrom > 0 {
		q.Where(loginlog.LoginAtGTE(time.Unix(req.LoginAtFrom, 0)))
	}
	if req.LoginAtTo > 0 {
		q.Where(loginlog.LoginAtLTE(time.Unix(req.LoginAtTo, 0)))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalize(20)
	list, err := q.Order(ent.Desc(loginlog.FieldLoginAt), ent.Desc(loginlog.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return loginLogsFromEnt(list), int64(total), nil
}
