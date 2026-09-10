package service

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
	"oa.98ent.com/p9/core/rpc/casbinx"
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/ent/api"
	"oa.98ent.com/p9/core/rpc/ent/menu"
	"oa.98ent.com/p9/core/rpc/ent/role"
	"oa.98ent.com/p9/core/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateMenuReq struct {
	ParentID   int64  `json:"parent_id"`
	MenuType   int16  `json:"menu_type"`
	Path       string `json:"path"`
	Name       string `json:"name"`
	Component  string `json:"component"`
	Redirect   string `json:"redirect"`
	Title      string `json:"title"`
	Icon       string `json:"icon"`
	Permission string `json:"permission"`
	HideMenu   int16  `json:"hide_menu"`
	Sort       int    `json:"sort"`
	Disabled   int16  `json:"disabled"`
}

type UpdateMenuReq struct {
	ID         int64   `json:"id"`
	ParentID   *int64  `json:"parent_id"`
	MenuType   *int16  `json:"menu_type"`
	Path       *string `json:"path"`
	Name       *string `json:"name"`
	Component  *string `json:"component"`
	Redirect   *string `json:"redirect"`
	Title      *string `json:"title"`
	Icon       *string `json:"icon"`
	Permission *string `json:"permission"`
	HideMenu   *int16  `json:"hide_menu"`
	Sort       *int    `json:"sort"`
	Disabled   *int16  `json:"disabled"`
}

type CreateAPIReq struct {
	Description string `json:"description"`
	APIGroup    string `json:"api_group"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	IsRequired  int16  `json:"is_required"`
	ServiceName string `json:"service_name"`
}

type UpdateAPIReq struct {
	ID          int64   `json:"id"`
	Description *string `json:"description"`
	APIGroup    *string `json:"api_group"`
	Method      *string `json:"method"`
	Path        *string `json:"path"`
	IsRequired  *int16  `json:"is_required"`
	ServiceName *string `json:"service_name"`
}

var apiMethods = map[string]struct{}{
	"GET": {}, "POST": {}, "PUT": {}, "PATCH": {}, "DELETE": {},
}

func (d *Deps) CreateMenu(ctx context.Context, req CreateMenuReq) (*model.Menu, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Title = strings.TrimSpace(req.Title)
	if req.Name == "" || req.Title == "" {
		return nil, xerr.BadRequest(i18n.MenuNameTitleRequired)
	}
	if !validMenuType(req.MenuType) {
		return nil, xerr.BadRequest(i18n.MenuInvalidType)
	}
	if err := d.assertMenuParent(ctx, 0, req.ParentID); err != nil {
		return nil, err
	}
	taken, err := d.menuNameTaken(ctx, req.Name, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, xerr.BadRequest(i18n.MenuNameExists)
	}
	row, err := d.Client.Menu.Create().
		SetParentID(req.ParentID).SetMenuType(req.MenuType).SetPath(strings.TrimSpace(req.Path)).
		SetName(req.Name).SetComponent(strings.TrimSpace(req.Component)).
		SetRedirect(strings.TrimSpace(req.Redirect)).SetTitle(req.Title).
		SetIcon(strings.TrimSpace(req.Icon)).SetPermission(strings.TrimSpace(req.Permission)).
		SetHideMenu(req.HideMenu).SetSort(req.Sort).SetDisabled(req.Disabled).
		Save(ctx)
	if err != nil {
		return nil, xerr.BadRequest(i18n.MenuCreateFailed)
	}
	if err := d.grantMenuToSupers(ctx, row.ID); err != nil {
		logx.Errorf("grant menu to supers failed: %v", err)
		return nil, err
	}
	m := menuFromEnt(row)
	return &m, nil
}

func (d *Deps) UpdateMenu(ctx context.Context, req UpdateMenuReq) error {
	row, err := d.menuByID(ctx, req.ID)
	if err != nil {
		return err
	}
	if err := d.prepareMenuUpdate(ctx, row, &req); err != nil {
		return err
	}
	return d.applyMenuUpdate(ctx, row.ID, req)
}

func (d *Deps) DeleteMenus(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := d.deleteMenu(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) CreateAPI(ctx context.Context, req CreateAPIReq) (*model.API, error) {
	norm, err := normalizeAPI(req.Method, req.Path, req.Description, req.APIGroup, req.ServiceName, req.IsRequired)
	if err != nil {
		return nil, err
	}
	taken, err := d.apiKeyTaken(ctx, norm.Method, norm.Path, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, xerr.BadRequest(i18n.APIExists)
	}
	row, err := d.Client.API.Create().
		SetDescription(norm.Description).SetAPIGroup(norm.APIGroup).SetMethod(norm.Method).
		SetPath(norm.Path).SetIsRequired(norm.IsRequired).SetServiceName(norm.ServiceName).
		Save(ctx)
	if err != nil {
		return nil, xerr.BadRequest(i18n.APICreateFailed)
	}
	if err := d.grantAPIToSupers(ctx, row.Path, row.Method); err != nil {
		return nil, err
	}
	a := apiFromEnt(row)
	return &a, nil
}

type RegisterMenuReq struct {
	Name       string
	Title      string
	MenuType   int16
	Path       string
	Component  string
	Redirect   string
	Icon       string
	Permission string
	HideMenu   int16
	Sort       int
	Disabled   int16
	ParentName string
}

// 注册目录。menus: 菜单，apis: API，items: 多语言词条，langs: 语言列表
func (d *Deps) RegisterCatalog(ctx context.Context, menus []RegisterMenuReq, apis []CreateAPIReq, items []I18nItem, langs []CreateI18nLangReq) error {
	if err := d.UpsertI18nLangs(ctx, langs); err != nil {
		return err
	}
	for _, m := range menus {
		if _, err := d.upsertRegisterMenu(ctx, m); err != nil {
			return err
		}
	}
	if err := d.linkRegisterMenuParents(ctx, menus); err != nil {
		return err
	}
	for _, a := range apis {
		if _, err := d.RegisterAPI(ctx, a); err != nil {
			return err
		}
	}
	for _, it := range items {
		if err := d.UpsertI18n(ctx, it); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) upsertRegisterMenu(ctx context.Context, req RegisterMenuReq) (*ent.Menu, error) {
	req.Name = strings.TrimSpace(req.Name)
	req.Title = strings.TrimSpace(req.Title)
	if req.Name == "" || req.Title == "" {
		return nil, xerr.BadRequest(i18n.MenuNameTitleRequired)
	}
	if !validMenuType(req.MenuType) {
		return nil, xerr.BadRequest(i18n.MenuInvalidType)
	}
	exist, err := d.Client.Menu.Query().Where(menu.NameEQ(req.Name)).Only(ctx)
	if ent.IsNotFound(err) {
		row, err := d.Client.Menu.Create().
			SetParentID(0).SetMenuType(req.MenuType).SetPath(strings.TrimSpace(req.Path)).
			SetName(req.Name).SetComponent(strings.TrimSpace(req.Component)).
			SetRedirect(strings.TrimSpace(req.Redirect)).SetTitle(req.Title).
			SetIcon(strings.TrimSpace(req.Icon)).SetPermission(strings.TrimSpace(req.Permission)).
			SetHideMenu(req.HideMenu).SetSort(req.Sort).SetDisabled(req.Disabled).
			Save(ctx)
		if err != nil {
			return nil, err
		}
		if err := d.grantMenuToSupers(ctx, row.ID); err != nil {
			return nil, err
		}
		return row, nil
	}
	if err != nil {
		return nil, err
	}
	if err := d.Client.Menu.UpdateOne(exist).
		SetMenuType(req.MenuType).SetPath(strings.TrimSpace(req.Path)).
		SetComponent(strings.TrimSpace(req.Component)).SetRedirect(strings.TrimSpace(req.Redirect)).
		SetTitle(req.Title).SetIcon(strings.TrimSpace(req.Icon)).SetPermission(strings.TrimSpace(req.Permission)).
		SetHideMenu(req.HideMenu).SetSort(req.Sort).SetDisabled(req.Disabled).
		Exec(ctx); err != nil {
		return nil, err
	}
	if err := d.grantMenuToSupers(ctx, exist.ID); err != nil {
		return nil, err
	}
	return exist, nil
}

func (d *Deps) linkRegisterMenuParents(ctx context.Context, menus []RegisterMenuReq) error {
	all, err := d.Client.Menu.Query().All(ctx)
	if err != nil {
		return err
	}
	ids := map[string]int64{}
	for _, m := range all {
		ids[m.Name] = m.ID
	}
	for _, m := range menus {
		if strings.TrimSpace(m.ParentName) == "" {
			continue
		}
		pid, ok := ids[m.ParentName]
		if !ok {
			return xerr.BadRequest(i18n.MenuParentNotFound)
		}
		if _, err := d.Client.Menu.Update().Where(menu.NameEQ(m.Name)).SetParentID(pid).Save(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) RegisterAPI(ctx context.Context, req CreateAPIReq) (*model.API, error) {
	norm, err := normalizeAPI(req.Method, req.Path, req.Description, req.APIGroup, req.ServiceName, req.IsRequired)
	if err != nil {
		return nil, err
	}
	row, err := d.Client.API.Query().Where(api.MethodEQ(norm.Method), api.PathEQ(norm.Path)).Only(ctx)
	if ent.IsNotFound(err) {
		row, err = d.Client.API.Create().
			SetDescription(norm.Description).SetAPIGroup(norm.APIGroup).SetMethod(norm.Method).
			SetPath(norm.Path).SetIsRequired(norm.IsRequired).SetServiceName(norm.ServiceName).
			Save(ctx)
	} else if err == nil {
		err = d.Client.API.UpdateOne(row).
			SetDescription(norm.Description).SetAPIGroup(norm.APIGroup).
			SetIsRequired(norm.IsRequired).SetServiceName(norm.ServiceName).
			Exec(ctx)
	}
	if err != nil {
		return nil, err
	}
	if err := d.grantAPIToSupers(ctx, norm.Path, norm.Method); err != nil {
		return nil, err
	}
	if row == nil {
		return nil, xerr.InternalServerError(i18n.APIRegisterFailed)
	}
	a := apiFromEnt(row)
	if a.Description != norm.Description {
		a.Description = norm.Description
		a.APIGroup = norm.APIGroup
		a.IsRequired = norm.IsRequired
		a.ServiceName = norm.ServiceName
	}
	return &a, nil
}

func (d *Deps) UpdateAPI(ctx context.Context, req UpdateAPIReq) error {
	row, err := d.apiByID(ctx, req.ID)
	if err != nil {
		return err
	}
	next := apiFromUpdate(row, req)
	norm, err := normalizeAPI(next.Method, next.Path, next.Description, next.APIGroup, next.ServiceName, next.IsRequired)
	if err != nil {
		return err
	}
	taken, err := d.apiKeyTaken(ctx, norm.Method, norm.Path, row.ID)
	if err != nil {
		return err
	}
	if taken {
		return xerr.BadRequest(i18n.APIExists)
	}
	if err := d.Client.API.UpdateOneID(row.ID).
		SetDescription(norm.Description).SetAPIGroup(norm.APIGroup).SetMethod(norm.Method).
		SetPath(norm.Path).SetIsRequired(norm.IsRequired).SetServiceName(norm.ServiceName).
		Exec(ctx); err != nil {
		return err
	}
	if row.Path != norm.Path || row.Method != norm.Method {
		return d.rewriteAPIPolicies(row.Path, row.Method, norm.Path, norm.Method)
	}
	return nil
}

func (d *Deps) DeleteAPIs(ctx context.Context, ids []int64) error {
	for _, id := range ids {
		if err := d.deleteAPI(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) SuperRoles(ctx context.Context) ([]*ent.Role, error) {
	return d.Client.Role.Query().
		Where(role.RoleCodeEQ(RoleSuperAdmin), role.IsSystemEQ(true), role.DeletedAtIsNil()).
		All(ctx)
}

func (d *Deps) prepareMenuUpdate(ctx context.Context, row model.Menu, req *UpdateMenuReq) error {
	if req.MenuType != nil && !validMenuType(*req.MenuType) {
		return xerr.BadRequest(i18n.MenuInvalidType)
	}
	if req.Name != nil {
		name := strings.TrimSpace(*req.Name)
		if name == "" {
			return xerr.BadRequest(i18n.MenuNameRequired)
		}
		taken, err := d.menuNameTaken(ctx, name, row.ID)
		if err != nil {
			return err
		}
		if taken {
			return xerr.BadRequest(i18n.MenuNameExists)
		}
		req.Name = &name
	}
	parentID := row.ParentID
	if req.ParentID != nil {
		parentID = *req.ParentID
	}
	return d.assertMenuParent(ctx, row.ID, parentID)
}

func (d *Deps) applyMenuUpdate(ctx context.Context, id int64, req UpdateMenuReq) error {
	upd := d.Client.Menu.UpdateOneID(id)
	if req.ParentID != nil {
		upd.SetParentID(*req.ParentID)
	}
	if req.MenuType != nil {
		upd.SetMenuType(*req.MenuType)
	}
	if req.Path != nil {
		upd.SetPath(strings.TrimSpace(*req.Path))
	}
	if req.Name != nil {
		upd.SetName(*req.Name)
	}
	if req.Component != nil {
		upd.SetComponent(strings.TrimSpace(*req.Component))
	}
	if req.Redirect != nil {
		upd.SetRedirect(strings.TrimSpace(*req.Redirect))
	}
	if req.Title != nil {
		upd.SetTitle(strings.TrimSpace(*req.Title))
	}
	if req.Icon != nil {
		upd.SetIcon(strings.TrimSpace(*req.Icon))
	}
	if req.Permission != nil {
		upd.SetPermission(strings.TrimSpace(*req.Permission))
	}
	if req.HideMenu != nil {
		upd.SetHideMenu(*req.HideMenu)
	}
	if req.Sort != nil {
		upd.SetSort(*req.Sort)
	}
	if req.Disabled != nil {
		upd.SetDisabled(*req.Disabled)
	}
	return upd.Exec(ctx)
}

func (d *Deps) deleteMenu(ctx context.Context, id int64) error {
	if _, err := d.menuByID(ctx, id); err != nil {
		return err
	}
	n, err := d.Client.Menu.Query().Where(menu.ParentIDEQ(id)).Count(ctx)
	if err != nil {
		return err
	}
	if n > 0 {
		return xerr.BadRequest(i18n.MenuHasChildren)
	}
	if _, err := d.Client.Menu.UpdateOneID(id).ClearRoles().Save(ctx); err != nil {
		return err
	}
	return d.Client.Menu.DeleteOneID(id).Exec(ctx)
}

func (d *Deps) deleteAPI(ctx context.Context, id int64) error {
	row, err := d.apiByID(ctx, id)
	if err != nil {
		return err
	}
	if err := d.removeAPIPolicies(row.Path, row.Method); err != nil {
		return err
	}
	return d.Client.API.DeleteOneID(id).Exec(ctx)
}

func (d *Deps) menuByID(ctx context.Context, id int64) (model.Menu, error) {
	row, err := d.Client.Menu.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.Menu{}, xerr.NotFound(i18n.MenuNotFound)
		}
		return model.Menu{}, err
	}
	return menuFromEnt(row), nil
}

func (d *Deps) apiByID(ctx context.Context, id int64) (model.API, error) {
	row, err := d.Client.API.Get(ctx, id)
	if err != nil {
		if ent.IsNotFound(err) {
			return model.API{}, xerr.NotFound(i18n.APINotFound)
		}
		return model.API{}, err
	}
	return apiFromEnt(row), nil
}

func (d *Deps) menuNameTaken(ctx context.Context, name string, exceptID int64) (bool, error) {
	q := d.Client.Menu.Query().Where(menu.NameEQ(name))
	if exceptID > 0 {
		q = q.Where(menu.IDNEQ(exceptID))
	}
	return q.Exist(ctx)
}

func (d *Deps) apiKeyTaken(ctx context.Context, method, path string, exceptID int64) (bool, error) {
	q := d.Client.API.Query().Where(api.MethodEQ(method), api.PathEQ(path))
	if exceptID > 0 {
		q = q.Where(api.IDNEQ(exceptID))
	}
	return q.Exist(ctx)
}

func (d *Deps) assertMenuParent(ctx context.Context, id, parentID int64) error {
	if parentID == 0 {
		return nil
	}
	if parentID == id {
		return xerr.BadRequest(i18n.MenuInvalidParentID)
	}
	parent, err := d.Client.Menu.Get(ctx, parentID)
	if err != nil {
		if ent.IsNotFound(err) {
			return xerr.BadRequest(i18n.MenuParentNotFound)
		}
		return err
	}
	if id == 0 {
		return nil
	}
	for i, cur := 0, parent.ParentID; i < 64 && cur != 0; i++ {
		if cur == id {
			return xerr.BadRequest(i18n.MenuInvalidParentID)
		}
		row, err := d.Client.Menu.Get(ctx, cur)
		if err != nil {
			return err
		}
		cur = row.ParentID
	}
	return nil
}

func (d *Deps) grantMenuToSupers(ctx context.Context, menuID int64) error {
	roles, err := d.SuperRoles(ctx)
	if err != nil {
		return err
	}
	for _, row := range roles {
		exist, err := d.Client.Role.Query().Where(role.ID(row.ID)).QueryMenus().Where(menu.ID(menuID)).Exist(ctx)
		if err != nil {
			return err
		}
		if exist {
			continue
		}
		if err := d.Client.Role.UpdateOneID(row.ID).AddMenuIDs(menuID).Exec(ctx); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) grantAPIToSupers(ctx context.Context, path, method string) error {
	if d.Mode == ModeOff && isOperatorAPI(path) {
		return nil
	}
	roles, err := d.SuperRoles(ctx)
	if err != nil {
		return err
	}
	for _, row := range roles {
		if _, err := d.Enforcer.AddPolicy(row.RoleCode, casbinx.Domain(row.OperatorID), path, method); err != nil {
			return err
		}
	}
	return nil
}

func (d *Deps) removeAPIPolicies(path, method string) error {
	return d.rewriteAPIPolicies(path, method, "", "")
}

func (d *Deps) rewriteAPIPolicies(oldPath, oldMethod, newPath, newMethod string) error {
	list, err := d.Enforcer.GetPolicy()
	if err != nil {
		return err
	}
	var remove, add [][]string
	for _, p := range list {
		if len(p) < 4 || p[2] != oldPath || p[3] != oldMethod {
			continue
		}
		remove = append(remove, p)
		if newPath != "" {
			np := append([]string(nil), p...)
			np[2], np[3] = newPath, newMethod
			add = append(add, np)
		}
	}
	if len(remove) == 0 {
		return nil
	}
	if _, err := d.Enforcer.RemovePolicies(remove); err != nil {
		return err
	}
	if len(add) == 0 {
		return nil
	}
	_, err = d.Enforcer.AddPolicies(add)
	return err
}

func validMenuType(t int16) bool {
	return t == model.MenuTypeDir || t == model.MenuTypeMenu || t == model.MenuTypeButton
}

type apiNorm struct {
	Description string
	APIGroup    string
	Method      string
	Path        string
	IsRequired  int16
	ServiceName string
}

func normalizeAPI(method, path, desc, group, svc string, required int16) (apiNorm, error) {
	out := apiNorm{
		Method: strings.ToUpper(strings.TrimSpace(method)),
		Path:   strings.TrimSpace(path), Description: strings.TrimSpace(desc),
		APIGroup: strings.TrimSpace(group), ServiceName: strings.TrimSpace(svc),
		IsRequired: required,
	}
	if _, ok := apiMethods[out.Method]; !ok {
		return apiNorm{}, xerr.BadRequest(i18n.APIInvalidMethod)
	}
	if out.Path == "" || out.Path[0] != '/' {
		return apiNorm{}, xerr.BadRequest(i18n.APIPathMustStartWithSlash)
	}
	if out.Description == "" {
		out.Description = out.Method + " " + out.Path
	}
	if out.APIGroup == "" {
		out.APIGroup = "biz"
	}
	if out.ServiceName == "" {
		out.ServiceName = "admin"
	}
	if out.IsRequired != 0 && out.IsRequired != 1 {
		return apiNorm{}, xerr.BadRequest(i18n.APIInvalidIsRequired)
	}
	return out, nil
}

func apiFromUpdate(row model.API, req UpdateAPIReq) model.API {
	if req.Description != nil {
		row.Description = *req.Description
	}
	if req.APIGroup != nil {
		row.APIGroup = *req.APIGroup
	}
	if req.Method != nil {
		row.Method = *req.Method
	}
	if req.Path != nil {
		row.Path = *req.Path
	}
	if req.IsRequired != nil {
		row.IsRequired = *req.IsRequired
	}
	if req.ServiceName != nil {
		row.ServiceName = *req.ServiceName
	}
	return row
}
