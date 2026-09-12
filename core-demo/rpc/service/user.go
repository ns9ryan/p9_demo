package service

import (
	"context"
	"strings"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/utils"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/role"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateUserReq struct {
	Username    string  `json:"username"`
	Password    string  `json:"password"`
	DisplayName string  `json:"display_name"`
	Mobile      *string `json:"mobile"`
	Email       *string `json:"email"`
	Status      int16   `json:"status"`
	RoleIDs     []int64 `json:"role_ids"`
}

type UpdateUserReq struct {
	ID          int64   `json:"id"`
	DisplayName *string `json:"display_name"`
	Mobile      *string `json:"mobile"`
	Email       *string `json:"email"`
	Status      *int16  `json:"status"`
}

type UpdateUserIpWhitelistReq struct {
	ID                 int64    `json:"id"`
	IPWhitelistEnabled int16    `json:"ip_whitelist_enabled"`
	IPWhitelist        []string `json:"ip_whitelist"`
}

type IDReq struct {
	ID int64 `json:"id"`
}

type IDsReq struct {
	IDs []int64 `json:"ids"`
}

const maxPageSize = 100

type PageReq struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

type UserListReq struct {
	PageReq
	Username    string
	Mobile      string
	Email       string
	DisplayName string
	RoleIDs     []int64
}

func (p *PageReq) normalize(defaultSize int) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = defaultSize
	}
	if p.PageSize > maxPageSize {
		p.PageSize = maxPageSize
	}
}

type BindRolesReq struct {
	UserID  int64   `json:"user_id"`
	RoleIDs []int64 `json:"role_ids"`
}

type PasswordReq struct {
	UserID      int64  `json:"user_id"`
	NewPassword string `json:"password"`
}

type SelfPasswordReq struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"password"`
}

func (d *Deps) CreateUser(ctx context.Context, claims *ctxdata.Claims, req CreateUserReq) (*model.User, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.Username == "" || req.Password == "" || req.DisplayName == "" {
		return nil, xerr.BadRequest(i18n.UserCreateFieldsRequired)
	}
	if req.Status == 0 {
		req.Status = model.StatusNormal
	}
	if d.Mode == ModeOn {
		if claims == nil || claims.OperatorID == 0 {
			return nil, xerr.Unauthorized(i18n.Unauthorized)
		}
	}
	hash, err := HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	salt, err := RandomSalt()
	if err != nil {
		return nil, err
	}
	row, err := d.Client.User.Create().
		SetUserCode(NewUserCode()).
		SetUsername(req.Username).
		SetPasswordHash(hash).
		SetSalt(salt).
		SetDisplayName(req.DisplayName).
		SetNillableMobile(req.Mobile).
		SetNillableEmail(req.Email).
		SetStatus(req.Status).
		SetIsSuperAdmin(false).
		Save(ctx)
	if err != nil {
		logx.Errorw(i18n.UserCreateFailed, logx.Field("error", err))
		return nil, xerr.BadRequest(i18n.UserCreateFailed)
	}
	u := userFromEnt(row)
	if len(req.RoleIDs) > 0 {
		if err := d.bindRoles(ctx, u, req.RoleIDs); err != nil {
			return nil, err
		}
	}
	return u, nil
}

func (d *Deps) UpdateUser(ctx context.Context, claims *ctxdata.Claims, req UpdateUserReq) error {
	u, err := d.mustTenantUser(ctx, claims, req.ID)
	if err != nil {
		return err
	}
	upd := d.Client.User.UpdateOneID(u.ID)
	if req.DisplayName != nil {
		upd.SetDisplayName(strings.TrimSpace(*req.DisplayName))
	}
	if req.Mobile != nil {
		upd.SetMobile(*req.Mobile)
	}
	if req.Email != nil {
		upd.SetEmail(*req.Email)
	}
	if req.Status != nil {
		upd.SetStatus(*req.Status)
		if *req.Status == model.StatusDisabled {
			if err := d.RotateSalt(ctx, u.ID); err != nil {
				return err
			}
		}
	}
	return upd.Exec(ctx)
}

func (d *Deps) UpdateUserIpWhitelist(ctx context.Context, claims *ctxdata.Claims, req UpdateUserIpWhitelistReq) error {
	if req.IPWhitelistEnabled != 0 && req.IPWhitelistEnabled != 1 {
		return xerr.BadRequest(i18n.InvalidParam)
	}
	list, ok := utils.NormalizeIPWhitelist(req.IPWhitelist)
	if !ok {
		return xerr.BadRequest(i18n.UserInvalidIpWhitelist)
	}
	if req.IPWhitelistEnabled == 1 && len(list) == 0 {
		return xerr.BadRequest(i18n.UserInvalidIpWhitelist)
	}
	u, err := d.mustTenantUser(ctx, claims, req.ID)
	if err != nil {
		return err
	}
	return d.Client.User.UpdateOneID(u.ID).
		SetIPWhitelistEnabled(req.IPWhitelistEnabled).
		SetIPWhitelist(list).
		Exec(ctx)
}

func (d *Deps) DeleteUsers(ctx context.Context, claims *ctxdata.Claims, ids []int64) error {
	if claims == nil {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	for _, id := range ids {
		if id == claims.UserID {
			return xerr.BadRequest(i18n.UserCannotDeleteSelf)
		}
		u, err := d.mustTenantUser(ctx, claims, id)
		if err != nil {
			return err
		}
		if u.IsSuperAdmin {
			return xerr.Forbidden(i18n.UserCannotDeleteRoot)
		}
		if err := d.Client.User.UpdateOneID(u.ID).SetDeletedAt(time.Now()).Exec(ctx); err != nil {
			return err
		}
		_ = d.RotateSalt(ctx, u.ID)
	}
	return nil
}

func (d *Deps) GetUser(ctx context.Context, claims *ctxdata.Claims, id int64) (*model.User, UserRoles, error) {
	u, err := d.mustTenantUser(ctx, claims, id)
	if err != nil {
		return nil, UserRoles{}, err
	}
	roles, err := d.RolesOfUser(ctx, u.ID)
	return u, roles, err
}

func (d *Deps) ListUsers(ctx context.Context, claims *ctxdata.Claims, req UserListReq) ([]model.User, int64, error) {
	if claims == nil {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	q := d.Client.User.Query().Where(user.DeletedAtIsNil())
	if s := strings.TrimSpace(req.Username); s != "" {
		q.Where(user.UsernameContains(s))
	}
	if s := strings.TrimSpace(req.Mobile); s != "" {
		q.Where(user.MobileContains(s))
	}
	if s := strings.TrimSpace(req.Email); s != "" {
		q.Where(user.EmailContains(s))
	}
	if s := strings.TrimSpace(req.DisplayName); s != "" {
		q.Where(user.DisplayNameContains(s))
	}
	if len(req.RoleIDs) > 0 {
		q.Where(user.HasRolesWith(role.IDIn(req.RoleIDs...), role.DeletedAtIsNil()))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalize(20)
	list, err := q.Order(ent.Desc(user.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	return usersFromEnt(list), int64(total), err
}

func (d *Deps) BindUserRoles(ctx context.Context, claims *ctxdata.Claims, req BindRolesReq) error {
	u, err := d.mustTenantUser(ctx, claims, req.UserID)
	if err != nil {
		return err
	}
	return d.bindRoles(ctx, u, req.RoleIDs)
}

func (d *Deps) bindRoles(ctx context.Context, u *model.User, roleIDs []int64) error {
	if len(roleIDs) == 0 {
		return d.Client.User.UpdateOneID(u.ID).ClearRoles().Exec(ctx)
	}
	roles, err := d.Client.Role.Query().Where(role.IDIn(roleIDs...), role.DeletedAtIsNil()).All(ctx)
	if err != nil {
		return err
	}
	if len(roles) != len(roleIDs) {
		return xerr.BadRequest(i18n.RoleNotFound)
	}
	return d.Client.User.UpdateOneID(u.ID).ClearRoles().AddRoleIDs(roleIDs...).Exec(ctx)
}

func (d *Deps) ChangePassword(ctx context.Context, claims *ctxdata.Claims, userID int64, newPassword string) error {
	if strings.TrimSpace(newPassword) == "" {
		return xerr.BadRequest(i18n.UserPasswordRequired)
	}
	u, err := d.mustTenantUser(ctx, claims, userID)
	if err != nil {
		return err
	}
	hash, err := HashPassword(newPassword)
	if err != nil {
		return err
	}
	if err := d.Client.User.UpdateOneID(u.ID).SetPasswordHash(hash).Exec(ctx); err != nil {
		return err
	}
	return d.RotateSalt(ctx, u.ID)
}

func (d *Deps) ChangeOwnPassword(ctx context.Context, claims *ctxdata.Claims, oldPw, newPw string) error {
	if claims == nil {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	u, err := d.ActiveUserByID(ctx, claims.UserID)
	if err != nil {
		return xerr.Unauthorized(i18n.Unauthorized)
	}
	if !CheckPassword(u.PasswordHash, oldPw) {
		return xerr.BadRequest(i18n.UserOldPasswordMismatch)
	}
	return d.ChangePassword(ctx, claims, u.ID, newPw)
}

func (d *Deps) mustTenantUser(ctx context.Context, claims *ctxdata.Claims, id int64) (*model.User, error) {
	u, err := d.ActiveUserByID(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.NotFound(i18n.UserNotFound)
		}
		return nil, err
	}
	return u, nil
}
