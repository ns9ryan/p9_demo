package convert

import (
	"context"

	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/core/api/internal/types"
	coreI18n "oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/core/rpc/model"
)

// 将字符串转换为指针
func strPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// 将int32转换为指针
func i32Ptr(v int32) *int32 {
	if v == 0 {
		return nil
	}
	return &v
}

// 将int64转换为指针
func i64Ptr(v int64) *int64 {
	if v == 0 {
		return nil
	}
	return &v
}

// 将用户公共信息转换为API响应
func UserPublic(ctx context.Context, trans *i18n.Translator, in *coreclient.UserPublic) *types.UserPublic {
	if in == nil {
		return nil
	}
	codes := in.RoleCodes
	if codes == nil {
		codes = []string{}
	}
	roleNames := in.RoleNames
	if roleNames == nil {
		roleNames = []string{}
	}
	transRoleNames := make([]string, 0, len(roleNames))
	for _, name := range roleNames {
		transRoleNames = append(transRoleNames, trans.T(ctx, name))
	}
	return &types.UserPublic{
		Id: in.Id, UserCode: in.UserCode, Username: in.Username, DisplayName: in.DisplayName,
		OperatorCode: in.GetOperatorCode(), IsSuperAdmin: in.IsSuperAdmin, Status: in.Status,
		RoleCodes: codes, RoleNames: transRoleNames, HomePath: in.HomePath,
		CreatedAt: in.CreatedAt, LastLoginAt: in.GetLastLoginAt(),
		Mobile: in.GetMobile(), Email: in.GetEmail(),
		IpWhitelistEnabled: in.IpWhitelistEnabled, IpWhitelist: copyStrSlice(in.IpWhitelist),
	}
}

// 复制字符串切片
func copyStrSlice(in []string) []string {
	if in == nil {
		return []string{}
	}
	out := make([]string, len(in))
	copy(out, in)
	return out
}

// 登录响应转换为API响应
func LoginResp(ctx context.Context, trans *i18n.Translator, in *coreclient.LoginResp) *types.LoginResp {
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
	if u := UserPublic(ctx, trans, in.User); u != nil {
		out.User = *u
	}
	return out
}

// 颁发预览令牌响应转换为API响应
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

// 权限响应转换为API响应
func PermResp(in *coreclient.PermResp) *types.PermResp {
	codes := in.GetPermissions()
	if codes == nil {
		codes = []string{}
	}
	return &types.PermResp{Permissions: codes}
}

// 菜单节点转换为API响应
func MenuNodes(ctx context.Context, code string, in []*coreclient.MenuNode) []types.MenuNode {
	out := make([]types.MenuNode, 0, len(in))
	for _, n := range in {
		out = append(out, menuNode(ctx, code, n))
	}
	return out
}

// 菜单节点转换为API响应
func menuNode(ctx context.Context, code string, n *coreclient.MenuNode) types.MenuNode {
	if n == nil {
		return types.MenuNode{Children: []types.MenuNode{}}
	}
	return types.MenuNode{
		Id: n.Id, ParentId: n.ParentId, MenuType: n.MenuType, Path: n.Path, Name: n.Name,
		Component: n.Component, Redirect: n.Redirect, Title: i18n.TG(ctx, code, i18n.GroupMenu, n.Title),
		Icon: n.Icon, Permission: n.Permission, HideMenu: n.HideMenu, Sort: n.Sort, Children: MenuNodes(ctx, code, n.Children),
	}
}

// 角色信息转换为API响应
func RoleInfo(ctx context.Context, trans *i18n.Translator, in *coreclient.RoleInfo) *types.RoleInfo {
	if in == nil {
		return nil
	}
	return &types.RoleInfo{
		Id: in.Id, OperatorCode: in.GetOperatorCode(), RoleCode: in.RoleCode, RoleName: trans.T(ctx, in.RoleName),
		Description: trans.T(ctx, in.GetDescription()), Status: in.Status, IsSystem: in.IsSystem, SortNo: in.SortNo,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

// 菜单信息转换为API响应
func MenuInfo(ctx context.Context, code string, in *coreclient.MenuInfo) *types.MenuInfo {
	if in == nil {
		return nil
	}
	return &types.MenuInfo{
		Id: in.Id, ParentId: in.ParentId, MenuType: in.MenuType, Path: in.Path, Name: in.Name, Component: in.Component,
		Redirect: in.Redirect, TransTitle: i18n.TG(ctx, code, i18n.GroupMenu, in.Title), Title: in.Title, Icon: in.Icon,
		Permission: in.Permission, HideMenu: in.HideMenu, Sort: in.Sort, Disabled: in.Disabled,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

// API信息转换为API响应
func ApiInfo(ctx context.Context, code string, in *coreclient.ApiInfo) *types.ApiInfo {
	if in == nil {
		return nil
	}
	return &types.ApiInfo{
		Id: in.Id, TransDescription: i18n.TG(ctx, code, i18n.GroupAPI, in.Description), Description: in.Description, ApiGroup: in.ApiGroup,
		Method: in.Method, Path: in.Path, IsRequired: in.IsRequired, ServiceName: in.ServiceName, CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

// 用户列表响应转换为API响应
func UserList(ctx context.Context, trans *i18n.Translator, in *coreclient.UserListResp) *types.UserListResp {
	list := make([]types.UserPublic, 0, len(in.GetList()))
	for _, u := range in.GetList() {
		if p := UserPublic(ctx, trans, u); p != nil {
			list = append(list, *p)
		}
	}
	return &types.UserListResp{List: list, Total: in.GetTotal()}
}

// 角色列表响应转换为API响应
func RoleList(ctx context.Context, trans *i18n.Translator, in *coreclient.RoleListResp) *types.RoleListResp {
	list := make([]types.RoleInfo, 0, len(in.GetList()))
	for _, r := range in.GetList() {
		if p := RoleInfo(ctx, trans, r); p != nil {
			list = append(list, *p)
		}
	}
	return &types.RoleListResp{List: list, Total: in.GetTotal()}
}

// 菜单信息列表转换为API响应
func MenuInfos(ctx context.Context, code string, in []*coreclient.MenuInfo) []types.MenuInfo {
	out := make([]types.MenuInfo, 0, len(in))
	for _, m := range in {
		if p := MenuInfo(ctx, code, m); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

// API信息列表转换为API响应
func ApiInfos(ctx context.Context, code string, in []*coreclient.ApiInfo) []types.ApiInfo {
	out := make([]types.ApiInfo, 0, len(in))
	for _, a := range in {
		if p := ApiInfo(ctx, code, a); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

// API列表响应转换为API响应
func ApiList(ctx context.Context, code string, in *coreclient.ApiListResp) *types.ApiListResp {
	return &types.ApiListResp{List: ApiInfos(ctx, code, in.GetList()), Total: in.GetTotal()}
}

// API列表请求转换为RPC请求
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

// API授权项列表转换为API响应
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

// 创建用户请求转换为RPC请求
func CreateUserReq(in *types.CreateUserReq) *coreclient.CreateUserReq {
	out := &coreclient.CreateUserReq{
		Username: in.Username, Password: in.Password, DisplayName: in.DisplayName,
		Status: in.Status, RoleIds: in.RoleIds,
	}
	out.Mobile = strPtr(in.Mobile)
	out.Email = strPtr(in.Email)
	return out
}

// 更新用户请求转换为RPC请求
func UpdateUserReq(in *types.UpdateUserReq) *coreclient.UpdateUserReq {
	return &coreclient.UpdateUserReq{
		Id: in.Id, DisplayName: strPtr(in.DisplayName), Mobile: strPtr(in.Mobile),
		Email: strPtr(in.Email), Status: i32Ptr(in.Status),
	}
}

// 更新用户IP白名单请求转换为RPC请求
func UpdateUserIpWhitelistReq(in *types.UpdateUserIpWhitelistReq) *coreclient.UpdateUserIpWhitelistReq {
	list := in.IpWhitelist
	if list == nil {
		list = []string{}
	}
	return &coreclient.UpdateUserIpWhitelistReq{
		Id: in.Id, IpWhitelistEnabled: in.IpWhitelistEnabled, IpWhitelist: list,
	}
}

// 创建角色请求转换为RPC请求
func CreateRoleReq(in *types.CreateRoleReq) *coreclient.CreateRoleReq {
	return &coreclient.CreateRoleReq{
		RoleCode: in.RoleCode, RoleName: in.RoleName, Description: in.Description,
		Status: in.Status, SortNo: in.SortNo,
	}
}

// 更新角色请求转换为RPC请求
func UpdateRoleReq(in *types.UpdateRoleReq) *coreclient.UpdateRoleReq {
	return &coreclient.UpdateRoleReq{
		Id: in.Id, RoleName: strPtr(in.RoleName), Description: strPtr(in.Description),
		Status: i32Ptr(in.Status), SortNo: i32Ptr(in.SortNo),
	}
}

// 创建菜单请求转换为RPC请求
func CreateMenuReq(in *types.CreateMenuReq) *coreclient.CreateMenuReq {
	return &coreclient.CreateMenuReq{
		ParentId: in.ParentId, MenuType: in.MenuType, Path: in.Path, Name: in.Name, Component: in.Component,
		Redirect: in.Redirect, Title: in.Title, Icon: in.Icon, Permission: in.Permission,
		HideMenu: in.HideMenu, Sort: in.Sort, Disabled: in.Disabled,
	}
}

// 更新菜单请求转换为RPC请求
func UpdateMenuReq(in *types.UpdateMenuReq) *coreclient.UpdateMenuReq {
	return &coreclient.UpdateMenuReq{
		Id: in.Id, ParentId: in.ParentId, MenuType: in.MenuType, Path: strPtr(in.Path),
		Name: strPtr(in.Name), Component: strPtr(in.Component), Redirect: strPtr(in.Redirect),
		Title: strPtr(in.Title), Icon: strPtr(in.Icon), Permission: strPtr(in.Permission),
		HideMenu: in.HideMenu, Sort: in.Sort, Disabled: in.Disabled,
	}
}

// 创建API请求转换为RPC请求
func CreateApiReq(in *types.CreateApiReq) *coreclient.CreateApiReq {
	return &coreclient.CreateApiReq{
		Description: in.Description, ApiGroup: in.ApiGroup, Method: in.Method, Path: in.Path,
		IsRequired: in.IsRequired, ServiceName: in.ServiceName,
	}
}

// 更新API请求转换为RPC请求
func UpdateApiReq(in *types.UpdateApiReq) *coreclient.UpdateApiReq {
	return &coreclient.UpdateApiReq{
		Id: in.Id, Description: strPtr(in.Description), ApiGroup: strPtr(in.ApiGroup),
		Method: strPtr(in.Method), Path: strPtr(in.Path), IsRequired: i32Ptr(in.IsRequired),
		ServiceName: strPtr(in.ServiceName),
	}
}

// IDs请求转换为RPC请求
func IDsReq(in *types.IDsReq) *coreclient.IDsReq {
	return &coreclient.IDsReq{Id: in.Id, Ids: in.Ids}
}

// 分页请求转换为RPC请求
func PageReq(in *types.PageReq) *coreclient.PageReq {
	return &coreclient.PageReq{Page: in.Page, PageSize: in.PageSize}
}

// 角色列表请求转换为RPC请求
func RoleListReq(in *types.RoleListReq) *coreclient.RoleListReq {
	return &coreclient.RoleListReq{
		Page:     in.Page,
		PageSize: in.PageSize,
		RoleName: in.RoleName,
	}
}

// 用户列表请求转换为RPC请求
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

// API授权请求转换为RPC请求
func ApiAuthReq(in *types.ApiAuthReq) *coreclient.ApiAuthReq {
	data := make([]*coreclient.ApiAuthItem, 0, len(in.Data))
	for _, it := range in.Data {
		data = append(data, &coreclient.ApiAuthItem{Path: it.Path, Method: it.Method})
	}
	return &coreclient.ApiAuthReq{RoleId: in.RoleId, Data: data}
}

// 登录日志列表请求转换为RPC请求
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

// 登录日志列表响应转换为API响应
func LoginLogList(ctx context.Context, trans *i18n.Translator, in *coreclient.LoginLogListResp) *types.LoginLogListResp {
	list := make([]types.LoginLogInfo, 0, len(in.GetList()))
	for _, row := range in.GetList() {
		list = append(list, loginLogInfo(ctx, trans, row))
	}
	return &types.LoginLogListResp{List: list, Total: in.GetTotal()}
}

// 登录日志信息转换为API响应
func loginLogInfo(ctx context.Context, trans *i18n.Translator, in *coreclient.LoginLogInfo) types.LoginLogInfo {
	if in == nil {
		return types.LoginLogInfo{}
	}
	out := types.LoginLogInfo{
		Id:          in.Id,
		UserId:      in.GetUserId(),
		Username:    in.Username,
		LoginResult: loginResultText(ctx, trans, in.LoginResult),
		LoginIp:     in.LoginIp,
		DeviceId:    in.GetDeviceId(),
		UserAgent:   in.GetUserAgent(),
		LoginAt:     in.LoginAt,
	}
	if reason := in.GetFailureReason(); reason != "" {
		out.FailureReason = trans.T(ctx, reason)
	}
	return out
}

// 登录结果文本转换为字符串
func loginResultText(ctx context.Context, trans *i18n.Translator, result int32) string {
	switch int16(result) {
	case model.LoginResultSuccess:
		return trans.T(ctx, coreI18n.LoginLogResultSuccess)
	case model.LoginResultFail:
		return trans.T(ctx, coreI18n.LoginLogResultFail)
	default:
		return ""
	}
}

// 管理员操作日志列表请求转换为RPC请求
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

// 管理员操作日志列表响应转换为API响应
func AdminActionLogList(ctx context.Context, trans *i18n.Translator, in *coreclient.AdminActionLogListResp) *types.AdminActionLogListResp {
	list := make([]types.AdminActionLogInfo, 0, len(in.GetList()))
	for _, row := range in.GetList() {
		list = append(list, adminActionLogInfo(ctx, trans, row))
	}
	return &types.AdminActionLogListResp{List: list, Total: in.GetTotal()}
}

// 管理员操作日志信息转换为API响应
func adminActionLogInfo(ctx context.Context, trans *i18n.Translator, in *coreclient.AdminActionLogInfo) types.AdminActionLogInfo {
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
		ActionResult:   loginResultText(ctx, trans, in.ActionResult),
		ResponseStatus: in.ResponseStatus,
		ResponseBody:   in.GetResponseBody(),
		DurationMs:     in.DurationMs,
		ClientIp:       in.ClientIp,
		UserAgent:      in.GetUserAgent(),
		CreatedAt:      in.CreatedAt,
	}
}

// 错误日志列表请求转换为RPC请求
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

// 错误日志列表响应转换为API响应
func ErrorLogList(in *coreclient.ErrorLogListResp) *types.ErrorLogListResp {
	list := make([]types.ErrorLogInfo, 0, len(in.GetList()))
	for _, row := range in.GetList() {
		list = append(list, errorLogInfo(row))
	}
	return &types.ErrorLogListResp{List: list, Total: in.GetTotal()}
}

// I18n信息转换为API响应
func I18nInfo(in *coreclient.I18NInfo) *types.I18nInfo {
	if in == nil {
		return nil
	}
	return &types.I18nInfo{
		Id: in.Id, I18nCode: in.I18NCode, I18nGroup: in.I18NGroup, TransKey: in.TransKey, Lang: in.Lang, Value: in.Value,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt,
	}
}

// I18n信息列表转换为API响应
func I18nInfos(in []*coreclient.I18NInfo) []types.I18nInfo {
	out := make([]types.I18nInfo, 0, len(in))
	for _, row := range in {
		if p := I18nInfo(row); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

// I18n列表响应转换为API响应
func I18nList(in *coreclient.I18NListResp) *types.I18nListResp {
	return &types.I18nListResp{List: I18nInfos(in.GetList()), Total: in.GetTotal()}
}

// 创建I18n请求转换为RPC请求
func CreateI18nReq(in *types.CreateI18nReq) *coreclient.CreateI18NReq {
	return &coreclient.CreateI18NReq{
		I18NCode: in.I18nCode, I18NGroup: in.I18nGroup, TransKey: in.TransKey, Lang: in.Lang, Value: in.Value,
	}
}

// 更新I18n请求转换为RPC请求
func UpdateI18nReq(in *types.UpdateI18nReq) *coreclient.UpdateI18NReq {
	return &coreclient.UpdateI18NReq{
		Id: in.Id, I18NCode: strPtr(in.I18nCode), I18NGroup: strPtr(in.I18nGroup), TransKey: strPtr(in.TransKey),
		Lang: strPtr(in.Lang), Value: strPtr(in.Value),
	}
}

// 更新I18nByKey请求转换为RPC请求
func UpdateI18nByKeyReq(in *types.UpdateI18nByKeyReq) *coreclient.UpdateI18NByKeyReq {
	return &coreclient.UpdateI18NByKeyReq{I18NCode: in.I18nCode, I18NGroup: in.I18nGroup, TransKey: in.TransKey, Data: in.Data}
}

func DeleteI18nByKeyReq(in *types.DeleteI18nByKeyReq) *coreclient.DeleteI18NByKeyReq {
	if in == nil {
		return &coreclient.DeleteI18NByKeyReq{}
	}
	return &coreclient.DeleteI18NByKeyReq{I18NCode: in.I18nCode, I18NGroup: in.I18nGroup, TransKey: in.TransKey}
}

// 获取I18n列表请求转换为RPC请求
func I18nListReq(in *types.I18nListReq) *coreclient.I18NListReq {
	return &coreclient.I18NListReq{
		Page: in.Page, PageSize: in.PageSize, I18NCode: in.I18nCode, I18NGroup: in.I18nGroup, TransKey: in.TransKey, Lang: in.Lang,
	}
}

// 获取I18n字典请求转换为RPC请求
func I18nDictReq(in *types.GetI18nDictReq) *coreclient.GetI18NDictReq {
	return &coreclient.GetI18NDictReq{I18NCode: in.I18nCode, I18NGroup: in.I18nGroup, Lang: in.Lang}
}

// 将I18n字典响应转换为API响应
func I18nDict(in *coreclient.I18NDictResp) *types.I18nDictResp {
	items := in.GetItems()
	if items == nil {
		items = map[string]string{}
	}
	return &types.I18nDictResp{Items: items}
}

// 导出I18n请求转换为RPC请求
func ExportI18nReq(in *types.ExportI18nReq) *coreclient.ExportI18NReq {
	return &coreclient.ExportI18NReq{Lang: in.Lang, I18NCode: in.I18nCode, I18NGroup: in.I18nGroup}
}

// 将I18n文件项转换为API响应
func I18nFileItems(in []*coreclient.I18NFileItem) []types.I18nFileItem {
	out := make([]types.I18nFileItem, 0, len(in))
	for _, row := range in {
		if row == nil {
			continue
		}
		out = append(out, types.I18nFileItem{
			I18nCode: row.I18NCode, I18nGroup: row.I18NGroup, I18nKey: row.I18NKey, I18nValue: row.I18NValue,
		})
	}
	return out
}

// 导出I18n响应转换为API响应
func ExportI18n(in *coreclient.ExportI18NResp) *types.ExportI18nResp {
	items := I18nFileItems(in.GetI18NItems())
	if items == nil {
		items = []types.I18nFileItem{}
	}
	return &types.ExportI18nResp{Lang: in.GetLang(), I18nItems: items}
}

// 导入I18n响应转换为API响应
func ImportI18n(in *coreclient.ImportI18NResp) *types.ImportI18nResp {
	return &types.ImportI18nResp{Created: in.GetCreated(), Updated: in.GetUpdated(), Skipped: in.GetSkipped()}
}

// 获取I18n语言信息
func I18nLangInfo(ctx context.Context, code string, in *coreclient.I18NLangInfo) *types.I18nLangInfo {
	if in == nil {
		return nil
	}
	i18nName := in.Name
	if in.I18NKey != "" {
		if t := i18n.TG(ctx, code, i18n.GroupLang, in.I18NKey); t != "" && t != in.I18NKey {
			i18nName = t
		}
	}
	return &types.I18nLangInfo{
		Id: in.Id, Lang: in.Lang, Name: in.Name, I18nKey: in.I18NKey, I18nName: i18nName, Disabled: in.Disabled,
		CreatedAt: in.CreatedAt, UpdatedAt: in.UpdatedAt, SortNo: in.SortNo,
	}
}

// 获取I18n语言列表
func I18nLangInfos(ctx context.Context, code string, in []*coreclient.I18NLangInfo) []types.I18nLangInfo {
	out := make([]types.I18nLangInfo, 0, len(in))
	for _, row := range in {
		if p := I18nLangInfo(ctx, code, row); p != nil {
			out = append(out, *p)
		}
	}
	return out
}

// 获取I18n语言列表响应
func I18nLangList(ctx context.Context, code string, in *coreclient.I18NLangListResp) *types.I18nLangListResp {
	return &types.I18nLangListResp{List: I18nLangInfos(ctx, code, in.GetList()), Total: in.GetTotal()}
}

// 创建I18n语言请求
func CreateI18nLangReq(in *types.CreateI18nLangReq) *coreclient.CreateI18NLangReq {
	return &coreclient.CreateI18NLangReq{Lang: in.Lang, Name: in.Name, I18NKey: in.I18nKey, Disabled: in.Disabled, SortNo: in.SortNo}
}

// 更新I18n语言请求
func UpdateI18nLangReq(in *types.UpdateI18nLangReq) *coreclient.UpdateI18NLangReq {
	return &coreclient.UpdateI18NLangReq{
		Id: in.Id, Lang: strPtr(in.Lang), Name: strPtr(in.Name), I18NKey: strPtr(in.I18nKey), Disabled: in.Disabled, SortNo: in.SortNo,
	}
}

// 重新排序I18n语言请求
func ReorderI18nLangReq(in *types.ReorderI18nLangReq) *coreclient.ReorderI18NLangReq {
	return &coreclient.ReorderI18NLangReq{Id: in.Id, TargetId: in.TargetId}
}

// 获取I18n语言列表请求
func I18nLangListReq(in *types.I18nLangListReq) *coreclient.I18NLangListReq {
	return &coreclient.I18NLangListReq{
		Page: in.Page, PageSize: in.PageSize, Lang: in.Lang, Disabled: in.Disabled,
	}
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
		OperatorCode:   in.GetOperatorCode(),
		CreatedAt:      in.CreatedAt,
	}
}
