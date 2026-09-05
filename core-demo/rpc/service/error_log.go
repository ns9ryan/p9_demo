package service

import (
	"context"
	"net"
	"strings"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/errorlog"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateErrorLogReq struct {
	UserID         int64
	OperatorID     int64
	RequestMethod  string
	RequestPath    string
	RequestQuery   string
	RequestBody    string
	ServiceName    string
	ResponseStatus int
	ResponseBody   string
	Subject        string
	Detail         string
	DurationMS     int
	ClientIP       string
	UserAgent      string
}

type ErrorLogListReq struct {
	PageReq
	UserID         int64
	RequestPath    string
	ServiceName    string
	ResponseStatus int
	CreatedAtFrom  int64
	CreatedAtTo    int64
}

func (d *Deps) CreateErrorLog(ctx context.Context, req CreateErrorLogReq) {
	if d == nil || d.Client == nil {
		return
	}
	method := clip(strings.ToUpper(strings.TrimSpace(req.RequestMethod)), 10)
	if method == "" {
		method = "POST"
	}
	path := clip(strings.TrimSpace(req.RequestPath), 500)
	if path == "" {
		path = "/"
	}
	name := clip(strings.TrimSpace(req.ServiceName), 100)
	if name == "" {
		name = "unknown"
	}
	ip := req.ClientIP
	if net.ParseIP(ip) == nil {
		ip = "0.0.0.0"
	}
	status := req.ResponseStatus
	if status <= 0 {
		status = 500
	}
	writeCtx := ctxdata.SkipTenant(ctx)
	c := d.Client.ErrorLog.Create().
		SetRequestMethod(method).
		SetRequestPath(path).
		SetServiceName(name).
		SetResponseStatus(status).
		SetDurationMs(req.DurationMS).
		SetClientIP(ip).
		SetCreatedAt(time.Now())
	if req.UserID != 0 {
		ok, err := d.Client.User.Query().Where(user.IDEQ(req.UserID)).Exist(writeCtx)
		if err == nil && ok {
			c.SetUserID(req.UserID)
		}
	}
	opID := req.OperatorID
	if opID == 0 {
		if claims := ctxdata.ClaimsFromCtx(ctx); claims != nil {
			opID = claims.OperatorID
		}
	}
	if opID != 0 {
		c.SetOperatorID(opID)
	}
	if s := strings.TrimSpace(req.RequestQuery); s != "" {
		c.SetRequestQuery(s)
	}
	if req.RequestBody != "" {
		c.SetRequestBody(req.RequestBody)
	}
	if req.ResponseBody != "" {
		c.SetResponseBody(req.ResponseBody)
	}
	if s := strings.TrimSpace(req.Subject); s != "" {
		c.SetSubject(s)
	}
	if req.Detail != "" {
		c.SetDetail(clip(req.Detail, 16*1024))
	}
	if ua := clip(req.UserAgent, 1000); ua != "" {
		c.SetUserAgent(ua)
	}
	if err := c.Exec(writeCtx); err != nil {
		logx.Errorf("write error log: %v", err)
	}
}

func (d *Deps) ListErrorLogs(ctx context.Context, claims *ctxdata.Claims, req ErrorLogListReq) ([]model.ErrorLog, int64, error) {
	if claims == nil {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	if d.Mode == ModeOn && claims.OperatorID == 0 {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	ctx = ctxdata.WithClaims(ctx, claims)
	q := d.Client.ErrorLog.Query().WithUser()
	if req.UserID != 0 {
		q.Where(errorlog.UserIDEQ(req.UserID))
	}
	if s := strings.TrimSpace(req.RequestPath); s != "" {
		q.Where(errorlog.RequestPathContains(s))
	}
	if s := strings.TrimSpace(req.ServiceName); s != "" {
		q.Where(errorlog.ServiceNameEQ(s))
	}
	if req.ResponseStatus != 0 {
		q.Where(errorlog.ResponseStatusEQ(req.ResponseStatus))
	}
	if req.CreatedAtFrom > 0 {
		q.Where(errorlog.CreatedAtGTE(time.Unix(req.CreatedAtFrom, 0)))
	}
	if req.CreatedAtTo > 0 {
		q.Where(errorlog.CreatedAtLTE(time.Unix(req.CreatedAtTo, 0)))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalize(20)
	list, err := q.Order(ent.Desc(errorlog.FieldCreatedAt), ent.Desc(errorlog.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return errorLogsFromEnt(list), int64(total), nil
}
