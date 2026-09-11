package service

import (
	"context"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/model"
)

func (d *Deps) Logout(ctx context.Context, refreshToken string) error {
	claims := ctxdata.ClaimsFromCtx(ctx)
	if claims == nil {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	raw, exp, err := d.parseOwnRefresh(claims, refreshToken)
	if err != nil {
		return err
	}
	if err := d.BlacklistToken(ctx, ctxdata.RawTokenFromCtx(ctx), claims.ExpiresAt); err != nil {
		return err
	}
	return d.BlacklistToken(ctx, raw, exp)
}

func (d *Deps) parseOwnRefresh(claims *ctxdata.Claims, refreshToken string) (string, int64, error) {
	raw := jwt.StripBearer(refreshToken)
	if raw == "" {
		return "", 0, xerr.BadRequest(i18n.AuthRefreshTokenRequired)
	}
	c, err := jwt.ParseTyped(d.JWTRefreshSecret, raw, jwt.TokenRefresh)
	if err != nil || c.UserID != claims.UserID || c.Salt != claims.Salt {
		return "", 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	exp := int64(0)
	if c.ExpiresAt != nil {
		exp = c.ExpiresAt.Unix()
	}
	return raw, exp, nil
}

func (d *Deps) LogoutAll(ctx context.Context) error {
	claims := ctxdata.ClaimsFromCtx(ctx)
	if claims == nil {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	return d.RotateSalt(ctx, claims.UserID)
}

func (d *Deps) CurrentUser(ctx context.Context) (UserPublic, error) {
	claims := ctxdata.ClaimsFromCtx(ctx)
	if claims == nil {
		return UserPublic{}, xerr.Unauthorized(i18n.Unauthorized)
	}
	u, err := d.ActiveUserByID(ctx, claims.UserID)
	if err != nil {
		return UserPublic{}, xerr.Unauthorized(i18n.Unauthorized)
	}
	roles, err := d.RolesOfUser(ctx, u.ID)
	if err != nil {
		return UserPublic{}, err
	}
	return toPublic(u, roles), nil
}

func (d *Deps) CheckToken(ctx context.Context, raw string) (*ctxdata.Claims, error) {
	raw = jwt.StripBearer(raw)
	claims, err := jwt.Parse(d.JWTSecret, raw)
	if err != nil {
		if jwt.IsExpired(err) {
			return nil, xerr.TokenExpired(i18n.TokenExpired)
		}
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if claims.TokenType == jwt.TokenRefresh {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if d.TokenBlacklisted(ctx, raw) {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	// if claims.TokenType == jwt.TokenPreview {
	// 	return d.checkPreviewToken(ctx, claims)
	// }
	u, err := d.ActiveUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if u.Status != model.StatusNormal || u.Salt != claims.Salt {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if err := d.checkTokenTenant(ctx, u, claims); err != nil {
		return nil, err
	}
	codes, err := d.RoleCodesOfUser(ctx, u.ID)
	if err != nil {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	exp := int64(0)
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Unix()
	}
	return &ctxdata.Claims{
		UserID:       u.ID,
		UserCode:     u.UserCode,
		Username:     u.Username,
		OperatorID:   claims.OperatorID,
		OperatorCode: d.operatorCodeOf(ctx, claims),
		RoleCodes:    codes,
		Salt:         u.Salt,
		ExpiresAt:    exp,
		TokenType:    claims.TokenType,
		IsPlatform:   claims.IsPlatform,
	}, nil
}

func (d *Deps) checkPreviewToken(ctx context.Context, claims *jwt.Claims) (*ctxdata.Claims, error) {
	if claims.OperatorID == 0 {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	op, err := d.Client.Operator.Get(ctxdata.SkipTenant(ctx), claims.OperatorID)
	if err != nil {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if op.Status != model.StatusNormal {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if claims.UserID == 0 {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	u, err := d.ActiveUserByID(ctxdata.SkipTenant(ctx), claims.UserID)
	if err != nil {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if u.Status != model.StatusNormal || u.Salt != claims.Salt {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	codes := claims.RoleCodes
	if codes == nil {
		codes = []string{}
	}
	exp := int64(0)
	if claims.ExpiresAt != nil {
		exp = claims.ExpiresAt.Unix()
	}
	code := claims.OperatorCode
	if code == "" {
		code = op.OperatorCode
	}
	return &ctxdata.Claims{
		UserID:       u.ID,
		UserCode:     u.UserCode,
		Username:     u.Username,
		OperatorID:   op.ID,
		OperatorCode: code,
		RoleCodes:    codes,
		Salt:         u.Salt,
		ExpiresAt:    exp,
		IsPlatform:   true,
		TokenType:    jwt.TokenPreview,
	}, nil
}

func (d *Deps) Enforce(ctx context.Context, claims *ctxdata.Claims, path, method string) (bool, error) {
	if claims == nil || len(claims.RoleCodes) == 0 {
		return false, nil
	}
	dom := d.DomainFromClaims(claims)
	reqs := make([][]interface{}, 0, len(claims.RoleCodes))
	for _, code := range claims.RoleCodes {
		reqs = append(reqs, []interface{}{code, dom, path, method})
	}
	okList, err := d.Enforcer.BatchEnforce(reqs)
	if err != nil {
		return false, err
	}
	for _, ok := range okList {
		if ok {
			return true, nil
		}
	}
	return false, nil
}

func (d *Deps) sessionFromClaims(ctx context.Context, c *jwt.Claims) (*model.User, UserRoles, error) {
	u, err := d.ActiveUserByID(ctx, c.UserID)
	if err != nil {
		return nil, UserRoles{}, xerr.Unauthorized(i18n.Unauthorized)
	}
	if u.Status != model.StatusNormal || u.Salt != c.Salt {
		return nil, UserRoles{}, xerr.Unauthorized(i18n.Unauthorized)
	}
	if err := d.checkTokenTenant(ctx, u, c); err != nil {
		return nil, UserRoles{}, err
	}
	roles, err := d.RolesOfUser(ctx, u.ID)
	if err != nil {
		return nil, UserRoles{}, err
	}
	if len(roles.Codes) == 0 {
		return nil, UserRoles{}, xerr.Forbidden(i18n.AuthNoActiveRole)
	}
	return u, roles, nil
}

func (d *Deps) checkTokenTenant(ctx context.Context, u *model.User, c *jwt.Claims) error {
	if d.Mode != ModeOn {
		if u.OperatorID != nil || c.OperatorID != 0 {
			return xerr.Unauthorized(i18n.Unauthorized)
		}
		return nil
	}
	if u.OperatorID == nil || *u.OperatorID != c.OperatorID || c.OperatorID == 0 {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	op, err := d.Client.Operator.Get(ctx, c.OperatorID)
	if err != nil {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	if op.Status != model.StatusNormal {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	return nil
}

func (d *Deps) operatorCodeOf(ctx context.Context, c *jwt.Claims) string {
	if c.OperatorCode != "" {
		return c.OperatorCode
	}
	if c.OperatorID == 0 {
		return ""
	}
	op, err := d.Client.Operator.Get(ctx, c.OperatorID)
	if err != nil {
		return ""
	}
	return op.OperatorCode
}

func PublicUser(u *model.User, roles UserRoles) UserPublic {
	return toPublic(u, roles)
}

func PublicUsers(list []model.User) []UserPublic {
	out := make([]UserPublic, 0, len(list))
	for i := range list {
		out = append(out, toPublic(&list[i], emptyUserRoles()))
	}
	return out
}
