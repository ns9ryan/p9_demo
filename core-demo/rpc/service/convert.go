package service

import (
	"oa.98ent.com/p9/core/rpc/ent"
	"oa.98ent.com/p9/core/rpc/model"
)

func userFromEnt(u *ent.User) *model.User {
	if u == nil {
		return nil
	}
	return &model.User{
		ID:                 u.ID,
		UserCode:           u.UserCode,
		OperatorID:         u.OperatorID,
		Username:           u.Username,
		PasswordHash:       u.PasswordHash,
		Salt:               u.Salt,
		DisplayName:        u.DisplayName,
		Mobile:             u.Mobile,
		Email:              u.Email,
		Status:             u.Status,
		IsSuperAdmin:       u.IsSuperAdmin,
		LastLoginAt:        u.LastLoginAt,
		LastLoginIP:        u.LastLoginIP,
		IPWhitelistEnabled: u.IPWhitelistEnabled,
		IPWhitelist:        copyStrings(u.IPWhitelist),
		CreatedAt:          u.CreatedAt,
		UpdatedAt:          u.UpdatedAt,
		DeletedAt:          u.DeletedAt,
	}
}

func copyStrings(in []string) []string {
	if in == nil {
		return []string{}
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

func usersFromEnt(list []*ent.User) []model.User {
	out := make([]model.User, 0, len(list))
	for _, u := range list {
		out = append(out, *userFromEnt(u))
	}
	return out
}

func roleFromEnt(r *ent.Role) *model.Role {
	if r == nil {
		return nil
	}
	return &model.Role{
		ID:          r.ID,
		OperatorID:  r.OperatorID,
		RoleCode:    r.RoleCode,
		RoleName:    r.RoleName,
		Description: r.Description,
		Status:      r.Status,
		IsSystem:    r.IsSystem,
		SortNo:      r.SortNo,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
		DeletedAt:   r.DeletedAt,
	}
}

func rolesFromEnt(list []*ent.Role) []model.Role {
	out := make([]model.Role, 0, len(list))
	for _, r := range list {
		out = append(out, *roleFromEnt(r))
	}
	return out
}

func operatorFromEnt(op *ent.Operator) *model.Operator {
	if op == nil {
		return nil
	}
	return &model.Operator{
		ID:                     op.ID,
		OperatorCode:           op.OperatorCode,
		TimezoneCode:           op.TimezoneCode,
		SettlementCurrencyCode: op.SettlementCurrencyCode,
		Status:                 op.Status,
		RequiredConfigVersion:  op.RequiredConfigVersion,
		CompletedConfigVersion: op.CompletedConfigVersion,
		ConfigCompletedAt:      op.ConfigCompletedAt,
		CreatedAt:              op.CreatedAt,
		UpdatedAt:              op.UpdatedAt,
	}
}

func menuFromEnt(m *ent.Menu) model.Menu {
	return model.Menu{
		ID:         m.ID,
		ParentID:   m.ParentID,
		MenuType:   m.MenuType,
		Path:       m.Path,
		Name:       m.Name,
		Component:  m.Component,
		Redirect:   m.Redirect,
		Title:      m.Title,
		Icon:       m.Icon,
		Permission: m.Permission,
		HideMenu:   m.HideMenu,
		Sort:       m.Sort,
		Disabled:   m.Disabled,
		CreatedAt:  m.CreatedAt,
		UpdatedAt:  m.UpdatedAt,
	}
}

func menusFromEnt(list []*ent.Menu) []model.Menu {
	out := make([]model.Menu, 0, len(list))
	for _, m := range list {
		out = append(out, menuFromEnt(m))
	}
	return out
}

func apiFromEnt(a *ent.API) model.API {
	return model.API{
		ID:          a.ID,
		Description: a.Description,
		APIGroup:    a.APIGroup,
		Method:      a.Method,
		Path:        a.Path,
		IsRequired:  a.IsRequired,
		ServiceName: a.ServiceName,
		CreatedAt:   a.CreatedAt,
		UpdatedAt:   a.UpdatedAt,
	}
}

func apisFromEnt(list []*ent.API) []model.API {
	out := make([]model.API, 0, len(list))
	for _, a := range list {
		out = append(out, apiFromEnt(a))
	}
	return out
}

func i18nFromEnt(row *ent.I18n) model.I18n {
	return model.I18n{
		ID:        row.ID,
		I18nCode:  row.I18nCode,
		I18nGroup: row.I18nGroup,
		TransKey:  row.TransKey,
		Lang:      row.Lang,
		Value:     row.Value,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func i18nsFromEnt(list []*ent.I18n) []model.I18n {
	out := make([]model.I18n, 0, len(list))
	for _, row := range list {
		out = append(out, i18nFromEnt(row))
	}
	return out
}

func i18nLangFromEnt(row *ent.I18nLang) model.I18nLang {
	return model.I18nLang{
		ID:        row.ID,
		Lang:      row.Lang,
		Name:      row.Name,
		Disabled:  row.Disabled,
		SortNo:    row.SortNo,
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
}

func i18nLangsFromEnt(list []*ent.I18nLang) []model.I18nLang {
	out := make([]model.I18nLang, 0, len(list))
	for _, row := range list {
		out = append(out, i18nLangFromEnt(row))
	}
	return out
}

func loginLogFromEnt(row *ent.LoginLog) model.LoginLog {
	return model.LoginLog{
		ID:            row.ID,
		UserID:        row.UserID,
		OperatorID:    row.OperatorID,
		Username:      row.Username,
		LoginResult:   row.LoginResult,
		FailureReason: row.FailureReason,
		LoginIP:       row.LoginIP,
		DeviceID:      row.DeviceID,
		UserAgent:     row.UserAgent,
		LoginAt:       row.LoginAt,
	}
}

func loginLogsFromEnt(list []*ent.LoginLog) []model.LoginLog {
	out := make([]model.LoginLog, 0, len(list))
	for _, row := range list {
		out = append(out, loginLogFromEnt(row))
	}
	return out
}

func adminActionLogFromEnt(row *ent.AdminActionLog) model.AdminActionLog {
	out := model.AdminActionLog{
		ID:             row.ID,
		UserID:         row.UserID,
		OperatorID:     row.OperatorID,
		RequestMethod:  row.RequestMethod,
		RequestPath:    row.RequestPath,
		RequestQuery:   row.RequestQuery,
		RequestBody:    row.RequestBody,
		ActionResult:   row.ActionResult,
		ResponseStatus: row.ResponseStatus,
		ResponseBody:   row.ResponseBody,
		DurationMS:     row.DurationMs,
		ClientIP:       row.ClientIP,
		UserAgent:      row.UserAgent,
		CreatedAt:      row.CreatedAt,
	}
	if row.Edges.User != nil {
		out.Username = row.Edges.User.Username
	}
	return out
}

func adminActionLogsFromEnt(list []*ent.AdminActionLog) []model.AdminActionLog {
	out := make([]model.AdminActionLog, 0, len(list))
	for _, row := range list {
		out = append(out, adminActionLogFromEnt(row))
	}
	return out
}

func errorLogFromEnt(row *ent.ErrorLog) model.ErrorLog {
	out := model.ErrorLog{
		ID:             row.ID,
		UserID:         row.UserID,
		OperatorID:     row.OperatorID,
		RequestMethod:  row.RequestMethod,
		RequestPath:    row.RequestPath,
		RequestQuery:   row.RequestQuery,
		RequestBody:    row.RequestBody,
		ServiceName:    row.ServiceName,
		ResponseStatus: row.ResponseStatus,
		ResponseBody:   row.ResponseBody,
		Subject:        row.Subject,
		Detail:         row.Detail,
		DurationMS:     row.DurationMs,
		ClientIP:       row.ClientIP,
		UserAgent:      row.UserAgent,
		CreatedAt:      row.CreatedAt,
	}
	if row.Edges.User != nil {
		out.Username = row.Edges.User.Username
	}
	return out
}

func errorLogsFromEnt(list []*ent.ErrorLog) []model.ErrorLog {
	out := make([]model.ErrorLog, 0, len(list))
	for _, row := range list {
		out = append(out, errorLogFromEnt(row))
	}
	return out
}
