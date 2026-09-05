package service

import (
	"context"
	"strings"
	"time"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/casbinx"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/role"
	"oa.98ent.com/p9/core/rpc/model"
)

type CreateRoleReq struct {
	RoleCode    string `json:"role_code"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
	Status      int16  `json:"status"`
	SortNo      int    `json:"sort_no"`
}

type UpdateRoleReq struct {
	ID          int64   `json:"id"`
	RoleName    *string `json:"role_name"`
	Description *string `json:"description"`
	Status      *int16  `json:"status"`
	SortNo      *int    `json:"sort_no"`
}

func (d *Deps) CreateRole(ctx context.Context, claims *ctxdata.Claims, req CreateRoleReq) (*model.Role, error) {
	req.RoleCode = strings.TrimSpace(req.RoleCode)
	req.RoleName = strings.TrimSpace(req.RoleName)
	if req.RoleCode == "" || req.RoleName == "" {
		return nil, xerr.BadRequest(i18n.RoleCodeNameRequired)
	}
	if req.Status == 0 {
		req.Status = model.StatusNormal
	}
	if d.Mode == ModeOn && (claims == nil || claims.OperatorID == 0) {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	row, err := d.Client.Role.Create().
		SetRoleCode(req.RoleCode).
		SetRoleName(req.RoleName).
		SetNillableDescription(strPtr(req.Description)).
		SetStatus(req.Status).
		SetIsSystem(false).
		SetSortNo(req.SortNo).
		Save(ctx)
	if err != nil {
		return nil, xerr.BadRequest(i18n.RoleCreateFailed)
	}
	return roleFromEnt(row), nil
}

func (d *Deps) UpdateRole(ctx context.Context, claims *ctxdata.Claims, req UpdateRoleReq) error {
	r, err := d.mustTenantRole(ctx, claims, req.ID)
	if err != nil {
		return err
	}
	// 系统角色不能禁用
	if r.IsSystem && req.Status != nil && *req.Status == model.StatusDisabled {
		return xerr.Forbidden(i18n.RoleCannotDisableSystem)
	}
	upd := d.Client.Role.UpdateOneID(r.ID)
	if req.RoleName != nil {
		upd.SetRoleName(strings.TrimSpace(*req.RoleName))
	}
	if req.Description != nil {
		upd.SetDescription(*req.Description)
	}
	if req.SortNo != nil {
		upd.SetSortNo(*req.SortNo)
	}
	if req.Status != nil {
		upd.SetStatus(*req.Status)
		if *req.Status == model.StatusDisabled {
			dom := casbinx.Domain(r.OperatorID)
			if _, err := d.Enforcer.RemoveFilteredPolicy(0, r.RoleCode, dom); err != nil {
				return err
			}
		}
	}
	return upd.Exec(ctx)
}

func (d *Deps) DeleteRoles(ctx context.Context, claims *ctxdata.Claims, ids []int64) error {
	for _, id := range ids {
		r, err := d.mustTenantRole(ctx, claims, id)
		if err != nil {
			return err
		}
		if r.IsSystem {
			return xerr.Forbidden(i18n.RoleCannotDeleteSystem)
		}
		n, err := d.Client.Role.Query().Where(role.ID(r.ID)).QueryUsers().Count(ctx)
		if err != nil {
			return err
		}
		if n > 0 {
			return xerr.BadRequest(i18n.RoleStillBoundToUsers)
		}
		dom := casbinx.Domain(r.OperatorID)
		if _, err := d.Enforcer.RemoveFilteredPolicy(0, r.RoleCode, dom); err != nil {
			return err
		}
		if err := d.Client.Role.UpdateOneID(r.ID).ClearMenus().Exec(ctx); err != nil {
			return err
		}
		if err := d.Client.Role.UpdateOneID(r.ID).SetDeletedAt(time.Now()).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) GetRole(ctx context.Context, claims *ctxdata.Claims, id int64) (*model.Role, error) {
	return d.mustTenantRole(ctx, claims, id)
}

type RoleListReq struct {
	PageReq
	RoleName string
}

func (d *Deps) ListRoles(ctx context.Context, claims *ctxdata.Claims, req RoleListReq) ([]model.Role, int64, error) {
	if claims == nil {
		return nil, 0, xerr.Unauthorized(i18n.Unauthorized)
	}
	q := d.Client.Role.Query().Where(role.DeletedAtIsNil())
	if s := strings.TrimSpace(req.RoleName); s != "" {
		q.Where(role.Or(
			role.RoleNameContains(s),
			role.RoleCodeContains(s),
		))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalize(50)
	list, err := q.Order(ent.Asc(role.FieldSortNo), ent.Asc(role.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	return rolesFromEnt(list), int64(total), err
}

func (d *Deps) mustTenantRole(ctx context.Context, claims *ctxdata.Claims, id int64) (*model.Role, error) {
	row, err := d.Client.Role.Query().Where(role.ID(id), role.DeletedAtIsNil()).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, xerr.NotFound(i18n.RoleNotFound)
		}
		return nil, err
	}
	r := roleFromEnt(row)
	return r, nil
}

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
