package service

import (
	"context"
	"net"
	"strings"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/operator"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginReq struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	OperatorCode string `json:"operator_code"`
	ClientIP     string `json:"-"`
	UserAgent    string `json:"-"`
}

type TokenInfo struct {
	AccessToken   string `json:"access_token"`
	RefreshToken  string `json:"refresh_token"`
	Expire        int64  `json:"expire"`
	RefreshExpire int64  `json:"refresh_expire"`
}

type LoginResult struct {
	Token TokenInfo  `json:"token"`
	User  UserPublic `json:"user"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

type LogoutReq struct {
	RefreshToken string `json:"refresh_token"`
}

type UserPublic struct {
	ID           int64    `json:"id"`
	UserCode     string   `json:"user_code"`
	Username     string   `json:"username"`
	DisplayName  string   `json:"display_name"`
	OperatorID   *int64   `json:"operator_id,omitempty"`
	IsSuperAdmin bool     `json:"is_super_admin"`
	Status       int16    `json:"status"`
	RoleCodes    []string `json:"role_codes"`
	RoleNames    []string `json:"role_names"`
	HomePath     string   `json:"home_path,omitempty"`
	CreatedAt    int64    `json:"created_at"`
	LastLoginAt  *int64   `json:"last_login_at,omitempty"`
	Mobile       *string  `json:"mobile,omitempty"`
	Email        *string  `json:"email,omitempty"`
}

func toPublic(u *model.User, roles UserRoles) UserPublic {
	if roles.Codes == nil {
		roles.Codes = []string{}
	}
	if roles.Names == nil {
		roles.Names = []string{}
	}
	out := UserPublic{
		ID:           u.ID,
		UserCode:     u.UserCode,
		Username:     u.Username,
		DisplayName:  u.DisplayName,
		OperatorID:   u.OperatorID,
		IsSuperAdmin: u.IsSuperAdmin,
		Status:       u.Status,
		RoleCodes:    roles.Codes,
		RoleNames:    roles.Names,
		HomePath:     "/dashboard",
		CreatedAt:    u.CreatedAt.Unix(),
		Mobile:       u.Mobile,
		Email:        u.Email,
	}
	if u.LastLoginAt != nil {
		ts := u.LastLoginAt.Unix()
		out.LastLoginAt = &ts
	}
	return out
}

func (d *Deps) Login(ctx context.Context, req LoginReq) (*LoginResult, error) {
	req.Username = strings.TrimSpace(req.Username)
	res, u, reason, err := d.doLogin(ctx, req)
	d.writeLoginLog(ctx, req, u, err == nil, reason)
	return res, err
}

func (d *Deps) doLogin(ctx context.Context, req LoginReq) (*LoginResult, *model.User, string, error) {
	if req.Username == "" || req.Password == "" {
		return nil, nil, i18n.AuthUsernamePasswordRequired, xerr.BadRequest(i18n.AuthUsernamePasswordRequired)
	}
	u, err := d.loginUser(ctx, req)
	if err != nil {
		return nil, nil, xerr.AsError(err).Message, err
	}
	if u.Status != model.StatusNormal {
		return nil, u, i18n.AuthUserDisabled, xerr.Forbidden(i18n.AuthUserDisabled)
	}
	if !CheckPassword(u.PasswordHash, req.Password) {
		return nil, u, i18n.AuthPasswordIncorrect, xerr.BadRequest(i18n.AuthPasswordIncorrect)
	}
	roles, err := d.RolesOfUser(ctx, u.ID)
	if err != nil {
		return nil, u, xerr.AsError(err).Message, err
	}
	if len(roles.Codes) == 0 {
		return nil, u, i18n.AuthNoActiveRole, xerr.Forbidden(i18n.AuthNoActiveRole)
	}
	tok, err := d.SignTokenPair(ctx, u, roles.Codes)
	if err != nil {
		return nil, u, xerr.AsError(err).Message, err
	}
	d.touchLogin(ctx, u, req.ClientIP)
	return &LoginResult{Token: tok, User: toPublic(u, roles)}, u, "", nil
}

func (d *Deps) writeLoginLog(ctx context.Context, req LoginReq, u *model.User, ok bool, reason string) {
	if d == nil || d.Client == nil {
		return
	}
	if opID := d.loginLogOperatorID(ctx, req, u); opID != 0 {
		ctx = ctxdata.WithClaims(ctx, &ctxdata.Claims{OperatorID: opID})
	}
	username := req.Username
	if username == "" {
		username = "-"
	}
	username = clip(username, 64)
	ip := req.ClientIP
	if net.ParseIP(ip) == nil {
		ip = "0.0.0.0"
	}
	ua := clip(req.UserAgent, 1000)
	c := d.Client.LoginLog.Create().
		SetUsername(username).
		SetLoginIP(ip).
		SetLoginAt(time.Now())
	if ok {
		c.SetLoginResult(model.LoginResultSuccess)
	} else {
		c.SetLoginResult(model.LoginResultFail)
		if reason != "" {
			c.SetFailureReason(clip(reason, 255))
		}
	}
	if u != nil {
		c.SetUserID(u.ID)
	}
	if ua != "" {
		c.SetUserAgent(ua)
	}
	if err := c.Exec(ctx); err != nil {
		logx.Errorf("write login log: %v", err)
	}
}

func (d *Deps) loginLogOperatorID(ctx context.Context, req LoginReq, u *model.User) int64 {
	if u != nil && u.OperatorID != nil && *u.OperatorID != 0 {
		return *u.OperatorID
	}
	if d.Mode != ModeOn {
		return 0
	}
	code := strings.TrimSpace(req.OperatorCode)
	if code == "" {
		return 0
	}
	op, err := d.Client.Operator.Query().Where(operator.OperatorCodeEQ(code)).Only(ctx)
	if err != nil {
		return 0
	}
	return op.ID
}

func clip(s string, n int) string {
	if n <= 0 || len(s) <= n {
		return s
	}
	return s[:n]
}

func (d *Deps) loginUser(ctx context.Context, req LoginReq) (*model.User, error) {
	q := d.Client.User.Query().Where(user.DeletedAtIsNil(), user.UsernameEqualFold(req.Username))
	if d.Mode != ModeOn {
		row, err := q.Where(user.OperatorIDIsNil()).Only(ctx)
		if err != nil {
			return nil, xerr.Unauthorized(i18n.AuthInvalidCredentials)
		}
		return userFromEnt(row), nil
	}
	code := strings.TrimSpace(req.OperatorCode)
	if code == "" {
		return nil, xerr.BadRequest(i18n.AuthOperatorCodeRequired)
	}
	op, err := d.Client.Operator.Query().Where(operator.OperatorCodeEQ(code)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.Unauthorized(i18n.AuthInvalidCredentials)
		}
		return nil, err
	}
	if op.Status != model.StatusNormal {
		return nil, xerr.Forbidden(i18n.AuthOperatorDisabled)
	}
	row, err := q.Where(user.OperatorIDEQ(op.ID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.BadRequest(i18n.AuthPasswordIncorrect)
		}
		return nil, xerr.InternalServerError(i18n.InternalError)
	}
	return userFromEnt(row), nil
}

func (d *Deps) touchLogin(ctx context.Context, u *model.User, ip string) {
	upd := d.Client.User.UpdateOneID(u.ID).SetLastLoginAt(time.Now())
	if parsed := net.ParseIP(ip); parsed != nil {
		upd.SetLastLoginIP(parsed.String())
	}
	_ = upd.Exec(ctx)
}

func (d *Deps) Refresh(ctx context.Context, req RefreshReq) (*LoginResult, error) {
	raw := jwt.StripBearer(req.RefreshToken)
	c, err := jwt.ParseTyped(d.JWTRefreshSecret, raw, jwt.TokenRefresh)
	if err != nil {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if d.TokenBlacklisted(ctx, raw) {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	u, roles, err := d.sessionFromClaims(ctx, c)
	if err != nil {
		return nil, err
	}
	tok, err := d.SignTokenPair(ctx, u, roles.Codes)
	if err != nil {
		return nil, err
	}
	exp := int64(0)
	if c.ExpiresAt != nil {
		exp = c.ExpiresAt.Unix()
	}
	_ = d.BlacklistToken(ctx, raw, exp)
	return &LoginResult{Token: tok, User: toPublic(u, roles)}, nil
}
