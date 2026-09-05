package service

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/common/ctxdata"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/casbinx"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/api"
	"oa.98ent.com/p9/core/rpc/ent/menu"
	"oa.98ent.com/p9/core/rpc/ent/role"
	"oa.98ent.com/p9/core/rpc/model"
)

type MenuAuthReq struct {
	RoleID  int64   `json:"role_id"`
	MenuIDs []int64 `json:"menu_ids"`
}

type APIAuthItem struct {
	Path   string `json:"path"`
	Method string `json:"method"`
}

type APIAuthReq struct {
	RoleID int64         `json:"role_id"`
	Data   []APIAuthItem `json:"data"`
}

func (d *Deps) UpdateMenuAuthority(ctx context.Context, claims *ctxdata.Claims, req MenuAuthReq) error {
	r, err := d.mustTenantRole(ctx, claims, req.RoleID)
	if err != nil {
		return err
	}
	if r.IsSystem {
		return xerr.Forbidden(i18n.AuthorityCannotChangeSystemMenus)
	}
	upd := d.Client.Role.UpdateOneID(r.ID).ClearMenus()
	if len(req.MenuIDs) > 0 {
		upd = upd.AddMenuIDs(req.MenuIDs...)
	}
	return upd.Exec(ctx)
}

func (d *Deps) GetMenuAuthority(ctx context.Context, claims *ctxdata.Claims, roleID int64) ([]int64, error) {
	r, err := d.mustTenantRole(ctx, claims, roleID)
	if err != nil {
		return nil, err
	}
	return d.Client.Role.Query().Where(role.ID(r.ID)).QueryMenus().IDs(ctx)
}

func (d *Deps) UpdateAPIAuthority(ctx context.Context, claims *ctxdata.Claims, req APIAuthReq) error {
	r, err := d.mustTenantRole(ctx, claims, req.RoleID)
	if err != nil {
		return err
	}
	if r.IsSystem {
		return xerr.Forbidden(i18n.AuthorityCannotChangeSystemAPIs)
	}
	required, err := d.Client.API.Query().Where(api.IsRequiredEQ(1)).All(ctx)
	if err != nil {
		return err
	}
	seen := map[string]APIAuthItem{}
	for _, it := range req.Data {
		seen[it.Method+" "+it.Path] = it
	}
	for _, a := range required {
		seen[a.Method+" "+a.Path] = APIAuthItem{Path: a.Path, Method: a.Method}
	}
	dom := casbinx.Domain(r.OperatorID)
	var policies [][]string
	for _, it := range seen {
		policies = append(policies, []string{r.RoleCode, dom, it.Path, it.Method})
	}
	return d.ReplaceRoleAPIPolicies(r.RoleCode, dom, policies)
}

func (d *Deps) GetAPIAuthority(ctx context.Context, claims *ctxdata.Claims, roleID int64) ([]APIAuthItem, error) {
	r, err := d.mustTenantRole(ctx, claims, roleID)
	if err != nil {
		return nil, err
	}
	dom := casbinx.Domain(r.OperatorID)
	list, err := d.Enforcer.GetFilteredPolicy(0, r.RoleCode, dom)
	if err != nil {
		return nil, err
	}
	out := make([]APIAuthItem, 0, len(list))
	for _, p := range list {
		if len(p) >= 4 {
			out = append(out, APIAuthItem{Path: p[2], Method: p[3]})
		}
	}
	return out, nil
}

func (d *Deps) ListMenus(ctx context.Context) ([]model.Menu, error) {
	list, err := d.Client.Menu.Query().Order(ent.Asc(menu.FieldSort), ent.Asc(menu.FieldID)).All(ctx)
	return menusFromEnt(list), err
}

type APIListReq struct {
	PageReq
	Path        string
	Method      string
	APIGroup    string
	ServiceName string
	Description string
}

func (d *Deps) ListAPIs(ctx context.Context, req APIListReq) ([]model.API, int64, error) {
	q := d.Client.API.Query()
	if s := strings.TrimSpace(req.Path); s != "" {
		q.Where(api.PathContains(s))
	}
	if s := strings.TrimSpace(req.Method); s != "" {
		q.Where(api.MethodEqualFold(s))
	}
	if s := strings.TrimSpace(req.APIGroup); s != "" {
		q.Where(api.APIGroupContains(s))
	}
	if s := strings.TrimSpace(req.ServiceName); s != "" {
		q.Where(api.ServiceNameContains(s))
	}
	if s := strings.TrimSpace(req.Description); s != "" {
		q.Where(api.DescriptionContains(s))
	}
	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	req.normalize(50)
	list, err := q.Order(ent.Asc(api.FieldAPIGroup), ent.Asc(api.FieldID)).
		Offset((req.Page - 1) * req.PageSize).
		Limit(req.PageSize).
		All(ctx)
	return apisFromEnt(list), int64(total), err
}

func (d *Deps) MenusByRole(ctx context.Context, claims *ctxdata.Claims) ([]model.Menu, error) {
	if claims == nil {
		return nil, xerr.Unauthorized(i18n.Unauthorized)
	}
	if len(claims.RoleCodes) == 0 {
		return make([]model.Menu, 0), nil
	}
	roles, err := d.Client.Role.Query().
		Where(role.RoleCodeIn(claims.RoleCodes...), role.DeletedAtIsNil(), role.StatusEQ(model.StatusNormal)).
		All(ctx)
	if err != nil {
		return nil, err
	}
	if len(roles) == 0 {
		return make([]model.Menu, 0), nil
	}
	ids := make([]int64, 0, len(roles))
	for _, r := range roles {
		ids = append(ids, r.ID)
	}
	q := d.Client.Menu.Query().Where(menu.DisabledEQ(0), menu.HideMenuEQ(0), menu.HasRolesWith(role.IDIn(ids...))).
		Unique(true).Order(ent.Asc(menu.FieldSort), ent.Asc(menu.FieldID))
	list, err := q.All(ctx)
	if err != nil {
		return nil, err
	}
	return uniqMenus(menusFromEnt(list)), nil
}

type MenuNode struct {
	ID         int64      `json:"id"`
	ParentID   int64      `json:"parent_id"`
	MenuType   int16      `json:"menu_type"`
	Path       string     `json:"path"`
	Name       string     `json:"name"`
	Component  string     `json:"component"`
	Redirect   string     `json:"redirect"`
	Title      string     `json:"title"`
	Icon       string     `json:"icon"`
	Permission string     `json:"permission"`
	HideMenu   int16      `json:"hide_menu"`
	Sort       int        `json:"sort"`
	Children   []MenuNode `json:"children"`
}

func (d *Deps) MenuTreeByRole(ctx context.Context, claims *ctxdata.Claims) ([]MenuNode, error) {
	menus, err := d.MenusByRole(ctx, claims)
	if err != nil {
		return nil, err
	}
	pages := make([]model.Menu, 0)
	for _, m := range menus {
		if m.Permission == "" {
			pages = append(pages, m)
		}
	}
	return buildMenuTree(pages, 0), nil
}

func buildMenuTree(menus []model.Menu, parentID int64) []MenuNode {
	nodes := make([]MenuNode, 0)
	for _, m := range menus {
		if m.ParentID != parentID {
			continue
		}
		n := MenuNode{
			ID: m.ID, ParentID: m.ParentID, MenuType: m.MenuType, Path: m.Path,
			Name: m.Name, Component: m.Component, Redirect: m.Redirect, Title: m.Title,
			Icon: m.Icon, Permission: m.Permission, HideMenu: m.HideMenu, Sort: m.Sort,
		}
		n.Children = buildMenuTree(menus, m.ID)
		nodes = append(nodes, n)
	}
	return nodes
}

func (d *Deps) PermCodes(ctx context.Context, claims *ctxdata.Claims) ([]string, error) {
	menus, err := d.MenusByRole(ctx, claims)
	if err != nil {
		return nil, err
	}
	codes := make([]string, 0)
	for _, m := range menus {
		if m.Permission != "" {
			codes = append(codes, m.Permission)
		}
	}
	return codes, nil
}

func uniqMenus(in []model.Menu) []model.Menu {
	seen := map[int64]struct{}{}
	out := make([]model.Menu, 0, len(in))
	for _, m := range in {
		if _, ok := seen[m.ID]; ok {
			continue
		}
		seen[m.ID] = struct{}{}
		out = append(out, m)
	}
	return out
}
