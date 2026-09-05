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
	"oa.98ent.com/p9/core/rpc/ent/adminactionlog"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateAdminActionLogReq struct {
	UserID         int64
	RequestMethod  string
	RequestPath    string
	RequestQuery   string
	RequestBody    string
	ActionResult   int16
	ResponseStatus int
	ResponseBody   string
	DurationMS     int
	ClientIP       string
	UserAgent      string
}

type AdminActionLogListReq struct {
	PageReq
	UserID        int64
	Username      string
	RequestMethod string
	RequestPath   string
	ActionResult  int16
	CreatedAtFrom int64
	CreatedAtTo   int64
}

func (d *Deps) CreateAdminActionLog(ctx context.Context, req CreateAdminActionLogReq) {
	if d == nil || d.Client == nil || req.UserID == 0 {
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
	ip := req.ClientIP
	if net.ParseIP(ip) == nil {
		ip = "0.0.0.0"
	}
	result := req.ActionResult
	if result != model.ActionResultSuccess {
		result = model.ActionResultFail
	}
	c := d.Client.AdminActionLog.Create().
		SetUserID(req.UserID).
		SetRequestMethod(method).
		SetRequestPath(path).
		SetActionResult(result).
		SetResponseStatus(req.ResponseStatus).
		SetDurationMs(req.DurationMS).
		SetClientIP(ip).
		SetCreatedAt(time.Now())
	if s := strings.TrimSpace(req.RequestQuery); s != "" {
		c.SetRequestQuery(s)
	}
	if req.RequestBody != "" {
		c.SetRequestBody(req.RequestBody)
	}
	if req.ResponseBody != "" {
		c.SetResponseBody(req.ResponseBody)
	}
	if ua := clip(req.UserAgent, 1000); ua != "" {
		c.SetUserAgent(ua)
	}
	if err := c.Exec(ctx); err != nil {
		logx.Errorf("write admin action log: %v", err)
	}
}

func (d *Deps) ListAdminActionLogs(ctx context.Context, claims *ctxdata.Claims, req AdminActionLogListReq) ([]model.AdminActionLog, int64, error) {
	if claims == nil {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	if d.Mode == ModeOn && claims.OperatorID == 0 {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	ctx = ctxdata.WithClaims(ctx, claims)
	q := d.Client.AdminActionLog.Query().WithUser()
	if req.UserID != 0 {
		q.Where(adminactionlog.UserIDEQ(req.UserID))
	}
	if s := strings.TrimSpace(req.Username); s != "" {
		q.Where(adminactionlog.HasUserWith(user.UsernameContains(s)))
	}
	if s := strings.TrimSpace(req.RequestMethod); s != "" {
		q.Where(adminactionlog.RequestMethodEQ(strings.ToUpper(s)))
	}
	if s := strings.TrimSpace(req.RequestPath); s != "" {
		q.Where(adminactionlog.RequestPathContains(s))
	}
	if req.ActionResult != 0 {
		q.Where(adminactionlog.ActionResultEQ(req.ActionResult))
	}
	if req.CreatedAtFrom > 0 {
		q.Where(adminactionlog.CreatedAtGTE(time.Unix(req.CreatedAtFrom, 0)))
	}
	if req.CreatedAtTo > 0 {
		q.Where(adminactionlog.CreatedAtLTE(time.Unix(req.CreatedAtTo, 0)))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalize(20)
	list, err := q.Order(ent.Desc(adminactionlog.FieldCreatedAt), ent.Desc(adminactionlog.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return adminActionLogsFromEnt(list), int64(total), nil
}
