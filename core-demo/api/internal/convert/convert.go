package convert

import (
	"context"

	"oa.98ent.com/p9/core/api/internal/types"
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/core/rpc/model"
)

func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func i32Ptr(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}

func i64Ptr(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

func UserPublic(in *coreclient.UserPublic) *types.UserPublic {
	if in == nil {
		return nil
	}
	codes := in.RoleCodes
	if codes == nil {
		codes = []string{}
	}
	return &types.UserPublic{
		Id: in.Id, UserCode: in.UserCode, Username: in.Username, DisplayName: in.DisplayName,
		OperatorId: in.GetOperatorId(), IsSuperAdmin: in.IsSuperAdmin, Status: in.Status,
		RoleCodes: codes, HomePath: in.HomePath,
		CreatedAt: in.CreatedAt, LastLoginAt: in.GetLastLoginAt(),
		Mobile: in.GetMobile(), Email: in.GetEmail(),
	}
}

func LoginResp(in *coreclient.LoginResp) *types.LoginResp {
	if in == nil {
		return nil
	}
	out := &types.LoginResp{}
	if in.Token != nil {
		out.Token = types.TokenInfo{
			AccessToken: in.Token.AccessToken, RefreshToken: in.Token.RefreshToken,
			Expire: in.Token.Expire, RefreshExpire: in.Token.RefreshExpire,
		}
	}
	if u := UserPublic(in.User); u != nil {
		out.User = *u
	}
	return out
}

func IssuePreviewTokenResp(in *coreclient.IssuePreviewTokenResp) *types.IssuePreviewTokenResp {
	if in == nil {
		return nil
	}
	return &types.IssuePreviewTokenResp{
		AccessToken:  in.AccessToken,
		Expire:       in.Expire,
		OperatorCode: in.OperatorCode,
		HomePath:     in.HomePath,
	}
}

func PermResp(in *coreclient.PermResp) *types.PermResp {
	codes := in.GetPermissions()
	if codes == nil {
		codes = []string{}
	}
	return &types.PermResp{Permissions: codes}
}

func MenuNodes(ctx context.Context, in []*coreclient.MenuNode) []types.MenuNode {
	out := make([]types.MenuNode, 0, len(in))
	for _, n := range in {
		out = append(out, menuNode(ctx, n))
	}
	return out
}

func menuNode(ctx context.Context, n *coreclient.MenuNode) types.MenuNode {
	if n == nil {
		return types.MenuNode{Children: []types.MenuNode{}}
	}
	return types.MenuNode{
		Id: n.Id, ParentId: n.ParentId, MenuType: n.MenuType, Path: n.Path, Name: n.Name,
		Component: n.Component, Redirect: n.Redirect, Title: i18n.T(ctx, n.Title), Icon: n.Icon,
		Permission: n.Permission, HideMenu: n.HideMenu, Sort: n.Sort, Children: MenuNodes(ctx, n.Children),
	}
}

func OperatorInfo(in *coreclient.OperatorInfo) *types.OperatorInfo {
	if in == nil {
		return nil
	}
	return &types.OperatorInfo{
		Id: in.Id, OperatorCode: in.OperatorCode, TimezoneCode: in.TimezoneCode,
		SettlementCurrencyCode: in.SettlementCurrencyCode, Status: in.Status,
		RequiredConfigVersion: in.RequiredConfigVersion, CompletedConfigVersion: in.CompletedConfigVersion,
		ConfigCompletedAt: in.ConfigCompletedAt, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

func RoleInfo(ctx context.Context, in *coreclient.RoleInfo) *types.RoleInfo {
	if in == nil {
		return nil
	}
	return &types.RoleInfo{
		Id: in.Id, OperatorId: in.GetOperatorId(), RoleCode: in.RoleCode, RoleName: i18n.T(ctx, in.RoleName),
		Description: i18n.T(ctx, in.GetDescription()), Status: in.Status, IsSystem: in.IsSystem, SortNo: in.SortNo,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

func MenuInfo(ctx context.Context, in *coreclient.MenuInfo) *types.MenuInfo {
	if in == nil {
		return nil
	}
	return &types.MenuInfo{
		Id: in.Id, ParentId: in.ParentId, MenuType: in.MenuType, Path: in.Path, Name: in.Name,
		Component: in.Component, Redirect: in.Redirect, Title: i18n.T(ctx, in.Title), Icon: in.Icon,
		Permission: in.Permission, HideMenu: in.HideMenu, Sort: in.Sort, Disabled: in.Disabled,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

func ApiInfo(ctx context.Context, in *coreclient.ApiInfo) *types.ApiInfo {
	if in == nil {
		return nil
	}
	return &types.ApiInfo{
		Id: in.Id, Description: i18n.T(ctx, in.Description), ApiGroup: in.ApiGroup, Method: in.Method, Path: in.Path,
		IsRequired: in.IsRequired, ServiceName: in.ServiceName, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

func UserList(in *coreclient.UserListResp) *types.UserListResp {
	list := make([]types.UserPublic, 0, len(in.GetList()))
	for _, u := range in.GetList() {
		if p := UserPublic(u); p != nil {
			list = append(list, *p)
		}
	}
	return &types.UserListResp{List: list, Total: in.GetTotal()}
}

func RoleList(ctx context.Context, in *coreclient.RoleListResp) *types.RoleListResp {
	list := make([]types.RoleInfo, 0, len(in.GetList()))
	for _, r := range in.GetList() {
		if p := RoleInfo(ctx, r); p != nil {
			list = append(list, *p)
		}
	}
	return &types.RoleListResp{List: list, Total: in.GetTotal()}
}

func MenuInfos(ctx context.Context, in []*coreclient.MenuInfo) []types.MenuInfo {
	out := make([]types.MenuInfo, 0, len(in))
	for _, m := range in {
		if p := MenuInfo(ctx, m); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

func ApiInfos(ctx context.Context, in []*coreclient.ApiInfo) []types.ApiInfo {
	out := make([]types.ApiInfo, 0, len(in))
	for _, a := range in {
		if p := ApiInfo(ctx, a); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

func ApiList(ctx context.Context, in *coreclient.ApiListResp) *types.ApiListResp {
	return &types.ApiListResp{List: ApiInfos(ctx, in.GetList()), Total: in.GetTotal()}
}

func ApiListReq(in *types.ApiListReq) *coreclient.ApiListReq {
	return &coreclient.ApiListReq{
		Page:        in.Page,
		PageSize:    in.PageSize,
		Path:        in.Path,
		Method:      in.Method,
		ApiGroup:    in.ApiGroup,
		ServiceName: in.ServiceName,
		Description: in.Description,
	}
}

func ApiAuthItems(in []*coreclient.ApiAuthItem) []types.ApiAuthItem {
	out := make([]types.ApiAuthItem, 0, len(in))
	for _, it := range in {
		if it == nil {
			continue
		}
		out = append(out, types.ApiAuthItem{Path: it.Path, Method: it.Method})
	}
	return out
}

func CreateUserReq(in *types.CreateUserReq) *coreclient.CreateUserReq {
	out := &coreclient.CreateUserReq{
		Username: in.Username, Password: in.Password, DisplayName: in.DisplayName,
		Status: in.Status, RoleIds: in.RoleIds,
	}
	out.Mobile = strPtr(in.Mobile)
	out.Email = strPtr(in.Email)
	return out
}

func UpdateUserReq(in *types.UpdateUserReq) *coreclient.UpdateUserReq {
	return &coreclient.UpdateUserReq{
		Id: in.Id, DisplayName: strPtr(in.DisplayName), Mobile: strPtr(in.Mobile),
		Email: strPtr(in.Email), Status: i32Ptr(in.Status),
	}
}

func UpdateOperatorReq(in *types.UpdateOperatorReq) *coreclient.UpdateOperatorReq {
	return &coreclient.UpdateOperatorReq{
		TimezoneCode: strPtr(in.TimezoneCode), SettlementCurrencyCode: strPtr(in.SettlementCurrencyCode),
	}
}

func CreateRoleReq(in *types.CreateRoleReq) *coreclient.CreateRoleReq {
	return &coreclient.CreateRoleReq{
		RoleCode: in.RoleCode, RoleName: in.RoleName, Description: in.Description,
		Status: in.Status, SortNo: in.SortNo,
	}
}

func UpdateRoleReq(in *types.UpdateRoleReq) *coreclient.UpdateRoleReq {
	return &coreclient.UpdateRoleReq{
		Id: in.Id, RoleName: strPtr(in.RoleName), Description: strPtr(in.Description),
		Status: i32Ptr(in.Status), SortNo: i32Ptr(in.SortNo),
	}
}

func CreateMenuReq(in *types.CreateMenuReq) *coreclient.CreateMenuReq {
	return &coreclient.CreateMenuReq{
		ParentId: in.ParentId, MenuType: in.MenuType, Path: in.Path, Name: in.Name, Component: in.Component,
		Redirect: in.Redirect, Title: in.Title, Icon: in.Icon, Permission: in.Permission,
		HideMenu: in.HideMenu, Sort: in.Sort, Disabled: in.Disabled,
	}
}

func UpdateMenuReq(in *types.UpdateMenuReq) *coreclient.UpdateMenuReq {
	return &coreclient.UpdateMenuReq{
		Id: in.Id, ParentId: in.ParentId, MenuType: in.MenuType, Path: strPtr(in.Path),
		Name: strPtr(in.Name), Component: strPtr(in.Component), Redirect: strPtr(in.Redirect),
		Title: strPtr(in.Title), Icon: strPtr(in.Icon), Permission: strPtr(in.Permission),
		HideMenu: in.HideMenu, Sort: in.Sort, Disabled: in.Disabled,
	}
}

func CreateApiReq(in *types.CreateApiReq) *coreclient.CreateApiReq {
	return &coreclient.CreateApiReq{
		Description: in.Description, ApiGroup: in.ApiGroup, Method: in.Method, Path: in.Path,
		IsRequired: in.IsRequired, ServiceName: in.ServiceName,
	}
}

func UpdateApiReq(in *types.UpdateApiReq) *coreclient.UpdateApiReq {
	return &coreclient.UpdateApiReq{
		Id: in.Id, Description: strPtr(in.Description), ApiGroup: strPtr(in.ApiGroup),
		Method: strPtr(in.Method), Path: strPtr(in.Path), IsRequired: i32Ptr(in.IsRequired),
		ServiceName: strPtr(in.ServiceName),
	}
}

func IDsReq(in *types.IDsReq) *coreclient.IDsReq {
	return &coreclient.IDsReq{Id: in.Id, Ids: in.Ids}
}

func PageReq(in *types.PageReq) *coreclient.PageReq {
	return &coreclient.PageReq{Page: in.Page, PageSize: in.PageSize}
}

func RoleListReq(in *types.RoleListReq) *coreclient.RoleListReq {
	return &coreclient.RoleListReq{
		Page:     in.Page,
		PageSize: in.PageSize,
		RoleName: in.RoleName,
	}
}

func UserListReq(in *types.UserListReq) *coreclient.UserListReq {
	return &coreclient.UserListReq{
		Page:        in.Page,
		PageSize:    in.PageSize,
		Username:    in.Username,
		Mobile:      in.Mobile,
		Email:       in.Email,
		DisplayName: in.DisplayName,
		RoleIds:     in.RoleIds,
	}
}

func ApiAuthReq(in *types.ApiAuthReq) *coreclient.ApiAuthReq {
	data := make([]*coreclient.ApiAuthItem, 0, len(in.Data))
	for _, it := range in.Data {
		data = append(data, &coreclient.ApiAuthItem{Path: it.Path, Method: it.Method})
	}
	return &coreclient.ApiAuthReq{RoleId: in.RoleId, Data: data}
}

func LoginLogListReq(in *types.LoginLogListReq) *coreclient.LoginLogListReq {
	return &coreclient.LoginLogListReq{
		Page:        in.Page,
		PageSize:    in.PageSize,
		Username:    in.Username,
		LoginResult: in.LoginResult,
		UserId:      in.UserId,
		LoginAtFrom: in.LoginAtFrom,
		LoginAtTo:   in.LoginAtTo,
	}
}

func LoginLogList(ctx context.Context, in *coreclient.LoginLogListResp) *types.LoginLogListResp {
	list := make([]types.LoginLogInfo, 0, len(in.GetList()))
	for _, row := range in.GetList() {
		list = append(list, loginLogInfo(ctx, row))
	}
	return &types.LoginLogListResp{List: list, Total: in.GetTotal()}
}

func loginLogInfo(ctx context.Context, in *coreclient.LoginLogInfo) types.LoginLogInfo {
	if in == nil {
		return types.LoginLogInfo{}
	}
	out := types.LoginLogInfo{
		Id:          in.Id,
		UserId:      in.GetUserId(),
		Username:    in.Username,
		LoginResult: loginResultText(ctx, in.LoginResult),
		LoginIp:     in.LoginIp,
		DeviceId:    in.GetDeviceId(),
		UserAgent:   in.GetUserAgent(),
		LoginAt:     in.LoginAt,
	}
	if reason := in.GetFailureReason(); reason != "" {
		out.FailureReason = i18n.T(ctx, reason)
	}
	return out
}

func loginResultText(ctx context.Context, result int32) string {
	switch int16(result) {
	case model.LoginResultSuccess:
		return i18n.T(ctx, i18n.LoginLogResultSuccess)
	case model.LoginResultFail:
		return i18n.T(ctx, i18n.LoginLogResultFail)
	default:
		return ""
	}
}

func AdminActionLogListReq(in *types.AdminActionLogListReq) *coreclient.AdminActionLogListReq {
	return &coreclient.AdminActionLogListReq{
		Page:          in.Page,
		PageSize:      in.PageSize,
		UserId:        in.UserId,
		Username:      in.Username,
		RequestMethod: in.RequestMethod,
		RequestPath:   in.RequestPath,
		ActionResult:  in.ActionResult,
		CreatedAtFrom: in.CreatedAtFrom,
		CreatedAtTo:   in.CreatedAtTo,
	}
}

func AdminActionLogList(ctx context.Context, in *coreclient.AdminActionLogListResp) *types.AdminActionLogListResp {
	list := make([]types.AdminActionLogInfo, 0, len(in.GetList()))
	for _, row := range in.GetList() {
		list = append(list, adminActionLogInfo(ctx, row))
	}
	return &types.AdminActionLogListResp{List: list, Total: in.GetTotal()}
}

func adminActionLogInfo(ctx context.Context, in *coreclient.AdminActionLogInfo) types.AdminActionLogInfo {
	if in == nil {
		return types.AdminActionLogInfo{}
	}
	return types.AdminActionLogInfo{
		Id:             in.Id,
		UserId:         in.UserId,
		Username:       in.Username,
		RequestMethod:  in.RequestMethod,
		RequestPath:    in.RequestPath,
		RequestQuery:   in.GetRequestQuery(),
		RequestBody:    in.GetRequestBody(),
		ActionResult:   loginResultText(ctx, in.ActionResult),
		ResponseStatus: in.ResponseStatus,
		ResponseBody:   in.GetResponseBody(),
		DurationMs:     in.DurationMs,
		ClientIp:       in.ClientIp,
		UserAgent:      in.GetUserAgent(),
		CreatedAt:      in.CreatedAt,
	}
}

func ErrorLogListReq(in *types.ErrorLogListReq) *coreclient.ErrorLogListReq {
	return &coreclient.ErrorLogListReq{
		Page:           in.Page,
		PageSize:       in.PageSize,
		UserId:         in.UserId,
		RequestPath:    in.RequestPath,
		ServiceName:    in.ServiceName,
		ResponseStatus: in.ResponseStatus,
		CreatedAtFrom:  in.CreatedAtFrom,
		CreatedAtTo:    in.CreatedAtTo,
	}
}

func ErrorLogList(in *coreclient.ErrorLogListResp) *types.ErrorLogListResp {
	list := make([]types.ErrorLogInfo, 0, len(in.GetList()))
	for _, row := range in.GetList() {
		list = append(list, errorLogInfo(row))
	}
	return &types.ErrorLogListResp{List: list, Total: in.GetTotal()}
}

func errorLogInfo(in *coreclient.ErrorLogInfo) types.ErrorLogInfo {
	if in == nil {
		return types.ErrorLogInfo{}
	}
	return types.ErrorLogInfo{
		Id:             in.Id,
		UserId:         in.GetUserId(),
		Username:       in.Username,
		RequestMethod:  in.RequestMethod,
		RequestPath:    in.RequestPath,
		RequestQuery:   in.GetRequestQuery(),
		RequestBody:    in.GetRequestBody(),
		ServiceName:    in.ServiceName,
		ResponseStatus: in.ResponseStatus,
		ResponseBody:   in.GetResponseBody(),
		Subject:        in.GetSubject(),
		Detail:         in.GetDetail(),
		DurationMs:     in.DurationMs,
		ClientIp:       in.ClientIp,
		UserAgent:      in.GetUserAgent(),
		OperatorId:     in.GetOperatorId(),
		CreatedAt:      in.CreatedAt,
	}
}
