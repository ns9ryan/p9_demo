package logic

import (
	"oa.98ent.com/p9/core/rpc/model"
	"oa.98ent.com/p9/core/rpc/service"
	"oa.98ent.com/p9/core/rpc/types/core"
)

func ToUserPublic(u service.UserPublic) *core.UserPublic {
	out := &core.UserPublic{
		Id:                 u.ID,
		UserCode:           u.UserCode,
		Username:           u.Username,
		DisplayName:        u.DisplayName,
		IsSuperAdmin:       u.IsSuperAdmin,
		Status:             int32(u.Status),
		RoleCodes:          u.RoleCodes,
		RoleNames:          u.RoleNames,
		HomePath:           u.HomePath,
		CreatedAt:          u.CreatedAt,
		IpWhitelistEnabled: int32(u.IPWhitelistEnabled),
		IpWhitelist:        u.IPWhitelist,
	}
	if u.OperatorID != nil {
		out.OperatorId = u.OperatorID
	}
	if u.LastLoginAt != nil {
		out.LastLoginAt = u.LastLoginAt
	}
	if u.Mobile != nil {
		out.Mobile = u.Mobile
	}
	if u.Email != nil {
		out.Email = u.Email
	}
	if out.RoleCodes == nil {
		out.RoleCodes = []string{}
	}
	if out.RoleNames == nil {
		out.RoleNames = []string{}
	}
	if out.IpWhitelist == nil {
		out.IpWhitelist = []string{}
	}
	return out
}

func ToLoginResp(res *service.LoginResult) *core.LoginResp {
	return &core.LoginResp{
		Token: &core.TokenInfo{
			AccessToken:   res.Token.AccessToken,
			RefreshToken:  res.Token.RefreshToken,
			Expire:        res.Token.Expire,
			RefreshExpire: res.Token.RefreshExpire,
		},
		User: ToUserPublic(res.User),
	}
}

func ToRoleInfo(r *model.Role) *core.RoleInfo {
	out := &core.RoleInfo{
		Id:        r.ID,
		RoleCode:  r.RoleCode,
		RoleName:  r.RoleName,
		Status:    int32(r.Status),
		IsSystem:  r.IsSystem,
		SortNo:    int32(r.SortNo),
		CreatedAt: r.CreatedAt.Unix(),
		UpdatedAt: r.UpdatedAt.Unix(),
	}
	out.OperatorId = r.OperatorID
	out.Description = r.Description
	return out
}

func ToMenuInfo(m model.Menu) *core.MenuInfo {
	return &core.MenuInfo{
		Id:         m.ID,
		ParentId:   m.ParentID,
		MenuType:   int32(m.MenuType),
		Path:       m.Path,
		Name:       m.Name,
		Component:  m.Component,
		Redirect:   m.Redirect,
		Title:      m.Title,
		Icon:       m.Icon,
		Permission: m.Permission,
		HideMenu:   int32(m.HideMenu),
		Sort:       int32(m.Sort),
		Disabled:   int32(m.Disabled),
		CreatedAt:  m.CreatedAt.Unix(),
		UpdatedAt:  m.UpdatedAt.Unix(),
	}
}

func ToApiInfo(a model.API) *core.ApiInfo {
	return &core.ApiInfo{
		Id:          a.ID,
		Description: a.Description,
		ApiGroup:    a.APIGroup,
		Method:      a.Method,
		Path:        a.Path,
		IsRequired:  int32(a.IsRequired),
		ServiceName: a.ServiceName,
		CreatedAt:   a.CreatedAt.Unix(),
		UpdatedAt:   a.UpdatedAt.Unix(),
	}
}

func ToI18nInfo(row model.I18n) *core.I18NInfo {
	return &core.I18NInfo{
		Id:        row.ID,
		I18NCode:  row.I18nCode,
		I18NGroup: row.I18nGroup,
		TransKey:  row.TransKey,
		Lang:      row.Lang,
		Value:     row.Value,
		CreatedAt: row.CreatedAt.Unix(),
		UpdatedAt: row.UpdatedAt.Unix(),
	}
}

func ToI18nLangInfo(row model.I18nLang) *core.I18NLangInfo {
	return &core.I18NLangInfo{
		Id:        row.ID,
		Lang:      row.Lang,
		Name:      row.Name,
		Disabled:  int32(row.Disabled),
		SortNo:    int32(row.SortNo),
		CreatedAt: row.CreatedAt.Unix(),
		UpdatedAt: row.UpdatedAt.Unix(),
	}
}

func ToOperatorInfo(op *model.Operator) *core.OperatorInfo {
	out := &core.OperatorInfo{
		Id:                     op.ID,
		OperatorCode:           op.OperatorCode,
		TimezoneCode:           op.TimezoneCode,
		SettlementCurrencyCode: op.SettlementCurrencyCode,
		Status:                 int32(op.Status),
		RequiredConfigVersion:  int32(op.RequiredConfigVersion),
		CompletedConfigVersion: int32(op.CompletedConfigVersion),
		CreatedAt:              op.CreatedAt.Unix(),
		UpdatedAt:              op.UpdatedAt.Unix(),
	}
	if op.ConfigCompletedAt != nil {
		out.ConfigCompletedAt = op.ConfigCompletedAt.Unix()
	}
	return out
}

func ToMenuNode(n service.MenuNode) *core.MenuNode {
	children := make([]*core.MenuNode, 0, len(n.Children))
	for _, c := range n.Children {
		children = append(children, ToMenuNode(c))
	}
	return &core.MenuNode{
		Id:         n.ID,
		ParentId:   n.ParentID,
		MenuType:   int32(n.MenuType),
		Path:       n.Path,
		Name:       n.Name,
		Component:  n.Component,
		Redirect:   n.Redirect,
		Title:      n.Title,
		Icon:       n.Icon,
		Permission: n.Permission,
		HideMenu:   int32(n.HideMenu),
		Sort:       int32(n.Sort),
		Children:   children,
	}
}

func ToInt16Ptr(v *int32) *int16 {
	if v == nil {
		return nil
	}
	n := int16(*v)
	return &n
}

func ToIntPtr(v *int32) *int {
	if v == nil {
		return nil
	}
	n := int(*v)
	return &n
}

func ToLoginLogInfo(row *model.LoginLog) *core.LoginLogInfo {
	if row == nil {
		return nil
	}
	out := &core.LoginLogInfo{
		Id:          row.ID,
		Username:    row.Username,
		LoginResult: int32(row.LoginResult),
		LoginIp:     row.LoginIP,
		LoginAt:     row.LoginAt.Unix(),
		UserId:      row.UserID,
		DeviceId:    row.DeviceID,
		UserAgent:   row.UserAgent,
	}
	if row.FailureReason != nil {
		out.FailureReason = row.FailureReason
	}
	return out
}

func ToAdminActionLogInfo(row *model.AdminActionLog) *core.AdminActionLogInfo {
	if row == nil {
		return nil
	}
	out := &core.AdminActionLogInfo{
		Id:             row.ID,
		UserId:         row.UserID,
		Username:       row.Username,
		RequestMethod:  row.RequestMethod,
		RequestPath:    row.RequestPath,
		ActionResult:   int32(row.ActionResult),
		ResponseStatus: int32(row.ResponseStatus),
		DurationMs:     int32(row.DurationMS),
		ClientIp:       row.ClientIP,
		CreatedAt:      row.CreatedAt.Unix(),
		RequestQuery:   row.RequestQuery,
		RequestBody:    row.RequestBody,
		ResponseBody:   row.ResponseBody,
		UserAgent:      row.UserAgent,
	}
	return out
}

func ToErrorLogInfo(row *model.ErrorLog) *core.ErrorLogInfo {
	if row == nil {
		return nil
	}
	out := &core.ErrorLogInfo{
		Id:             row.ID,
		Username:       row.Username,
		RequestMethod:  row.RequestMethod,
		RequestPath:    row.RequestPath,
		ServiceName:    row.ServiceName,
		ResponseStatus: int32(row.ResponseStatus),
		DurationMs:     int32(row.DurationMS),
		ClientIp:       row.ClientIP,
		CreatedAt:      row.CreatedAt.Unix(),
		UserId:         row.UserID,
		RequestQuery:   row.RequestQuery,
		RequestBody:    row.RequestBody,
		ResponseBody:   row.ResponseBody,
		Subject:        row.Subject,
		Detail:         row.Detail,
		UserAgent:      row.UserAgent,
		OperatorId:     row.OperatorID,
	}
	return out
}
