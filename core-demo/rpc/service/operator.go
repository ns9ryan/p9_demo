package service

import (
	"context"
	"strings"

	"oa.98ent.com/p9/common/ctxdata"
	"oa.98ent.com/p9/common/utils"
	"oa.98ent.com/p9/common/xerr"
	coreI18n "oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"
)

type IssuePreviewTokenReq struct {
	OperatorCode string
}

type PreviewToken struct {
	AccessToken  string
	Expire       int64
	OperatorCode string
	HomePath     string
}

func (d *Deps) IssuePreviewToken(ctx context.Context, req IssuePreviewTokenReq) (*PreviewToken, error) {
	if d.Mode != ModeOn {
		return nil, xerr.BadRequest(coreI18n.OperatorAPIOnlyOnMode)
	}
	code := strings.TrimSpace(req.OperatorCode)
	if code == "" {
		return nil, xerr.BadRequest(coreI18n.AuthOperatorCodeRequired)
	}
	admin, err := d.Client.User.Query().
		Where(
			user.OperatorCodeEQ(code),
			user.IsSuperAdminEQ(true),
			user.DeletedAtIsNil(),
			user.StatusEQ(model.StatusNormal),
		).
		Only(ctxdata.SkipTenant(ctx))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.NotFound(coreI18n.UserNotFound)
		}
		return nil, err
	}
	tok, exp, err := jwt.Sign(d.JWTSecret, d.JWTExpire, jwt.Claims{
		UserID:       admin.ID,
		UserCode:     admin.UserCode,
		Username:     admin.Username,
		OperatorCode: code,
		RoleCodes:    []string{RoleSuperAdmin},
		Salt:         admin.Salt,
		TokenType:    jwt.TokenPreview,
		IsPlatform:   true,
		ClientIP:     utils.NormalizeIP(ctxdata.ClientIPFromCtx(ctx)),
	})
	if err != nil {
		return nil, err
	}
	return &PreviewToken{
		AccessToken:  tok,
		Expire:       exp,
		OperatorCode: code,
		HomePath:     "/dashboard",
	}, nil
}
