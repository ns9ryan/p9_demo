package service

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/jwt"
	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/operator"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"
)

type UpdateOperatorReq struct {
	TimezoneCode           *string `json:"timezone_code"`
	SettlementCurrencyCode *string `json:"settlement_currency_code"`
}

type IssuePreviewTokenReq struct {
	OperatorCode string
}

type PreviewToken struct {
	AccessToken  string
	Expire       int64
	OperatorCode string
	HomePath     string
}

func (d *Deps) GetOperatorSelf(ctx context.Context, claims *ctxdata.Claims) (*model.Operator, error) {
	if d.Mode != ModeOn {
		return nil, xerr.BadRequest(i18n.OperatorAPIOnlyOnMode)
	}
	if claims == nil || claims.OperatorID == 0 {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	op, err := d.Client.Operator.Get(ctx, claims.OperatorID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.NotFound(i18n.OperatorNotFound)
		}
		return nil, err
	}
	return operatorFromEnt(op), nil
}

func (d *Deps) UpdateOperatorSelf(ctx context.Context, claims *ctxdata.Claims, req UpdateOperatorReq) error {
	op, err := d.GetOperatorSelf(ctx, claims)
	if err != nil {
		return err
	}
	upd := d.Client.Operator.UpdateOneID(op.ID)
	if req.TimezoneCode != nil {
		upd.SetTimezoneCode(strings.TrimSpace(*req.TimezoneCode))
	}
	if req.SettlementCurrencyCode != nil {
		upd.SetSettlementCurrencyCode(strings.TrimSpace(*req.SettlementCurrencyCode))
	}
	return upd.Exec(ctx)
}

func (d *Deps) IssuePreviewToken(ctx context.Context, req IssuePreviewTokenReq) (*PreviewToken, error) {
	if d.Mode != ModeOn {
		return nil, xerr.BadRequest(i18n.OperatorAPIOnlyOnMode)
	}
	code := strings.TrimSpace(req.OperatorCode)
	if code == "" {
		return nil, xerr.BadRequest(i18n.AuthOperatorCodeRequired)
	}
	op, err := d.Client.Operator.Query().
		Where(operator.OperatorCodeEQ(code)).
		Only(ctxdata.SkipTenant(ctx))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.NotFound(i18n.OperatorNotFound)
		}
		return nil, err
	}
	if op.Status != model.StatusNormal {
		return nil, xerr.Forbidden(i18n.AuthOperatorDisabled)
	}
	admin, err := d.Client.User.Query().
		Where(
			user.OperatorIDEQ(op.ID),
			user.IsSuperAdminEQ(true),
			user.DeletedAtIsNil(),
			user.StatusEQ(model.StatusNormal),
		).
		Only(ctxdata.SkipTenant(ctx))
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.NotFound(i18n.UserNotFound)
		}
		return nil, err
	}
	tok, exp, err := jwt.Sign(d.JWTSecret, d.JWTExpire, jwt.Claims{
		UserID:       admin.ID,
		UserCode:     admin.UserCode,
		Username:     admin.Username,
		OperatorID:   op.ID,
		OperatorCode: op.OperatorCode,
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
		OperatorCode: op.OperatorCode,
		HomePath:     "/dashboard",
	}, nil
}
