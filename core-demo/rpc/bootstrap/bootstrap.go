package bootstrap

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/casbinx"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/menu"
	"oa.98ent.com/p9/core/rpc/ent/role"
	"oa.98ent.com/p9/core/rpc/ent/user"
	"oa.98ent.com/p9/core/rpc/model"
	"oa.98ent.com/p9/core/rpc/service"
)

type CreateAdminReq struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	DisplayName string `json:"display_name"`
}

type CreateOperatorAdminReq struct {
	OperatorCode string `json:"operator_code"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	DisplayName  string `json:"display_name"`
}

// CreatePlatformAdmin 创建总网超级管理员
func CreatePlatformAdmin(ctx context.Context, d *service.Deps, req CreateAdminReq) (*model.User, error) {
	if err := d.RequireMode(service.ModeOff); err != nil {
		return nil, err
	}
	return createRoot(ctx, d, createRootReq{
		Username: req.Username, Password: req.Password, DisplayName: req.DisplayName,
	})
}

// CreateOperatorAdmin 创建分站超级管理员
func CreateOperatorAdmin(ctx context.Context, d *service.Deps, req CreateOperatorAdminReq) (*model.User, error) {
	if err := d.RequireMode(service.ModeOn); err != nil {
		return nil, err
	}
	req.OperatorCode = strings.TrimSpace(req.OperatorCode)
	if req.OperatorCode == "" || req.Username == "" || req.Password == "" {
		return nil, xerr.BadRequest(i18n.AuthBootstrapFieldsRequired)
	}
	if req.DisplayName == "" {
		req.DisplayName = req.Username
	}
	code := req.OperatorCode
	return createRoot(ctx, d, createRootReq{
		OperatorCode: &code, Username: req.Username, Password: req.Password, DisplayName: req.DisplayName,
	})
}

type createRootReq struct {
	OperatorCode *string
	Username     string
	Password     string
	DisplayName  string
}

func createRoot(ctx context.Context, d *service.Deps, req createRootReq) (*model.User, error) {
	req.Username = strings.TrimSpace(req.Username)
	req.DisplayName = strings.TrimSpace(req.DisplayName)
	if req.Username == "" || req.Password == "" {
		return nil, xerr.BadRequest(i18n.AuthUsernamePasswordRequired)
	}
	if req.DisplayName == "" {
		req.DisplayName = req.Username
	}
	if err := assertNoRoot(ctx, d, req.OperatorCode); err != nil {
		return nil, err
	}
	hash, err := service.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}
	salt, err := service.RandomSalt()
	if err != nil {
		return nil, err
	}
	var created *ent.User
	err = withTx(ctx, d, func(nd *service.Deps) error {
		u, err := insertRoot(ctx, nd, req, hash, salt)
		created = u
		return err
	})
	if err != nil {
		return nil, err
	}
	return userModel(created), nil
}

func assertNoRoot(ctx context.Context, d *service.Deps, operatorCode *string) error {
	q := d.Client.User.Query().Where(user.DeletedAtIsNil(), user.IsSuperAdminEQ(true))
	if operatorCode == nil || *operatorCode == "" {
		q = q.Where(user.OperatorCodeIsNil())
	} else {
		q = q.Where(user.OperatorCodeEQ(*operatorCode))
	}
	n, err := q.Count(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return xerr.BadRequest(i18n.AuthRootUserExists)
	}
	return nil
}

func insertRoot(ctx context.Context, d *service.Deps, req createRootReq, hash, salt string) (*ent.User, error) {
	roleRow, err := ensureSuperRole(ctx, d.Client, req.OperatorCode)
	if err != nil {
		return nil, err
	}
	if err := grantAllMenus(ctx, d.Client, roleRow.ID); err != nil {
		return nil, err
	}
	u, err := d.Client.User.Create().
		SetUserCode(service.NewUserCode()).
		SetNillableOperatorCode(req.OperatorCode).
		SetUsername(req.Username).
		SetPasswordHash(hash).
		SetSalt(salt).
		SetDisplayName(req.DisplayName).
		SetStatus(model.StatusNormal).
		SetIsSuperAdmin(true).
		Save(ctx)
	if err != nil {
		return nil, err
	}
	if err := d.Client.User.UpdateOneID(u.ID).AddRoleIDs(roleRow.ID).Exec(ctx); err != nil {
		return nil, err
	}
	dom := casbinx.Domain(req.OperatorCode)
	policies, err := d.AllAPIPolicies(ctx, dom)
	if err != nil {
		return nil, err
	}
	if err := d.ReplaceRoleAPIPolicies(service.RoleSuperAdmin, dom, policies); err != nil {
		return nil, err
	}
	return u, nil
}

func grantAllMenus(ctx context.Context, c *ent.Client, roleID int64) error {
	menus, err := c.Menu.Query().Where(menu.DisabledEQ(0)).All(ctx)
	if err != nil {
		return err
	}
	ids := make([]int64, 0, len(menus))
	for _, m := range menus {
		ids = append(ids, m.ID)
	}
	upd := c.Role.UpdateOneID(roleID).ClearMenus()
	if len(ids) > 0 {
		upd = upd.AddMenuIDs(ids...)
	}
	return upd.Exec(ctx)
}

func ensureSuperRole(ctx context.Context, c *ent.Client, operatorCode *string) (*ent.Role, error) {
	q := c.Role.Query().Where(role.RoleCodeEqualFold(service.RoleSuperAdmin))
	if operatorCode == nil || *operatorCode == "" {
		q = q.Where(role.OperatorCodeIsNil())
	} else {
		q = q.Where(role.OperatorCodeEQ(*operatorCode))
	}
	row, err := q.Only(ctx)
	if ent.IsNotFound(err) {
		return c.Role.Create().
			SetNillableOperatorCode(operatorCode).
			SetRoleCode(service.RoleSuperAdmin).
			SetRoleName("role.superAdmin").
			SetStatus(model.StatusNormal).
			SetIsSystem(true).
			SetSortNo(0).
			Save(ctx)
	}
	if err != nil {
		return nil, err
	}
	return c.Role.UpdateOne(row).
		SetRoleName("role.superAdmin").
		SetIsSystem(true).
		SetStatus(model.StatusNormal).
		ClearDeletedAt().
		Save(ctx)
}

func withTx(ctx context.Context, d *service.Deps, fn func(*service.Deps) error) error {
	tx, err := d.Client.Tx(ctx)
	if err != nil {
		return err
	}
	nd := *d
	nd.Client = tx.Client()
	if err := fn(&nd); err != nil {
		_ = tx.Rollback()
		return err
	}
	return tx.Commit()
}

func userModel(u *ent.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		ID:           u.ID,
		UserCode:     u.UserCode,
		OperatorCode: u.OperatorCode,
		Username:     u.Username,
		PasswordHash: u.PasswordHash,
		Salt:         u.Salt,
		DisplayName:  u.DisplayName,
		Status:       u.Status,
		IsSuperAdmin: u.IsSuperAdmin,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}
}
