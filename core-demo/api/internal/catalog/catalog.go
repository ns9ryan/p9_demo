package catalog

import (
	"net/http"

	"oa.98ent.com/p9/core/rpc/coreclient"
)

const (
	menuTypeDir    int32 = 0
	menuTypeMenu   int32 = 1
	menuTypeButton int32 = 2
)

// CatalogReq 菜单、API目录、多语言数据种子数据
func catalogReq(code string) *coreclient.RegisterCatalogReq {
	return &coreclient.RegisterCatalogReq{
		Menus: []*coreclient.RegisterMenuReq{
			{Name: "Dashboard", Title: "menu.route.dashboard", Path: "/dashboard", MenuType: menuTypeMenu, Component: "dashboard/index", Sort: 1},
			{Name: "System", Title: "menu.route.system", Path: "/system", MenuType: menuTypeDir, Sort: 10},
			{Name: "User", Title: "menu.route.user", Path: "/system/user", MenuType: menuTypeMenu, Component: "system/user/index", ParentName: "System", Sort: 11},
			{Name: "Role", Title: "menu.route.role", Path: "/system/role", MenuType: menuTypeMenu, Component: "system/role/index", ParentName: "System", Sort: 12},
			{Name: "Menu", Title: "menu.route.menu", Path: "/system/menu", MenuType: menuTypeMenu, Component: "system/menu/index", ParentName: "System", Sort: 13},
			{Name: "API", Title: "menu.route.api", Path: "/system/api", MenuType: menuTypeMenu, Component: "system/api/index", ParentName: "System", Sort: 14},
			{Name: "I18n", Title: "menu.route.i18n", Path: "/system/i18n", MenuType: menuTypeDir, ParentName: "System", Sort: 16},
			{Name: "I18nEntry", Title: "menu.route.i18nEntry", Path: "/system/i18n/entry", MenuType: menuTypeMenu, Component: "system/i18n/entry/index", ParentName: "I18n", Sort: 161},
			{Name: "I18nLang", Title: "menu.route.i18nLang", Path: "/system/i18n/lang", MenuType: menuTypeMenu, Component: "system/i18n/lang/index", ParentName: "I18n", Sort: 162},
			{Name: "Log", Title: "menu.route.log", Path: "/system/log", MenuType: menuTypeDir, ParentName: "System", Sort: 15},
			{Name: "LoginLog", Title: "menu.route.loginLog", Path: "/system/log/login", MenuType: menuTypeMenu, Component: "system/log/login/index", ParentName: "Log", Sort: 151},
			{Name: "ActionLog", Title: "menu.route.actionLog", Path: "/system/log/action", MenuType: menuTypeMenu, Component: "system/log/action/index", ParentName: "Log", Sort: 152},
			{Name: "ErrorLog", Title: "menu.route.errorLog", Path: "/system/log/error", MenuType: menuTypeMenu, Component: "system/log/error/index", ParentName: "Log", Sort: 153},
			{Name: "UserCreate", Title: "menu.route.userCreate", MenuType: menuTypeButton, Permission: "user:create", ParentName: "User", Sort: 111},
			{Name: "RoleCreate", Title: "menu.route.roleCreate", MenuType: menuTypeButton, Permission: "role:create", ParentName: "Role", Sort: 121},
			{Name: "MenuCreate", Title: "menu.route.menuCreate", MenuType: menuTypeButton, Permission: "menu:create", ParentName: "Menu", Sort: 131},
			{Name: "APICreate", Title: "menu.route.apiCreate", MenuType: menuTypeButton, Permission: "api:create", ParentName: "API", Sort: 141},
			{Name: "I18nCreate", Title: "menu.route.i18nCreate", MenuType: menuTypeButton, Permission: "i18n:create", ParentName: "I18nEntry", Sort: 1611},
			{Name: "I18nLangCreate", Title: "menu.route.i18nLangCreate", MenuType: menuTypeButton, Permission: "i18nLang:create", ParentName: "I18nLang", Sort: 1621},
		},
		Apis: []*coreclient.CreateApiReq{
			{Path: "/admin/operator/self", Method: http.MethodGet, Description: "api.operatorSelf", ApiGroup: "operator", ServiceName: "core-api"},
			{Path: "/admin/operator/update", Method: http.MethodPost, Description: "api.operatorUpdate", ApiGroup: "operator", ServiceName: "core-api"},
			{Path: "/admin/user/create", Method: http.MethodPost, Description: "api.userCreate", ApiGroup: "user", ServiceName: "core-api"},
			{Path: "/admin/user/update", Method: http.MethodPost, Description: "api.userUpdate", ApiGroup: "user", ServiceName: "core-api"},
			{Path: "/admin/user/delete", Method: http.MethodPost, Description: "api.userDelete", ApiGroup: "user", ServiceName: "core-api"},
			{Path: "/admin/user/list", Method: http.MethodPost, Description: "api.userList", ApiGroup: "user", ServiceName: "core-api"},
			{Path: "/admin/user/detail", Method: http.MethodGet, Description: "api.userDetail", ApiGroup: "user", ServiceName: "core-api"},
			{Path: "/admin/user/password", Method: http.MethodPost, Description: "api.userPassword", ApiGroup: "user", ServiceName: "core-api"},
			{Path: "/admin/user/roles", Method: http.MethodPost, Description: "api.userRoles", ApiGroup: "user", ServiceName: "core-api"},
			{Path: "/admin/role/create", Method: http.MethodPost, Description: "api.roleCreate", ApiGroup: "role", ServiceName: "core-api"},
			{Path: "/admin/role/update", Method: http.MethodPost, Description: "api.roleUpdate", ApiGroup: "role", ServiceName: "core-api"},
			{Path: "/admin/role/delete", Method: http.MethodPost, Description: "api.roleDelete", ApiGroup: "role", ServiceName: "core-api"},
			{Path: "/admin/role/list", Method: http.MethodPost, Description: "api.roleList", ApiGroup: "role", ServiceName: "core-api"},
			{Path: "/admin/role/detail", Method: http.MethodGet, Description: "api.roleDetail", ApiGroup: "role", ServiceName: "core-api"},
			{Path: "/admin/menu/create", Method: http.MethodPost, Description: "api.menuCreate", ApiGroup: "menu", ServiceName: "core-api"},
			{Path: "/admin/menu/update", Method: http.MethodPost, Description: "api.menuUpdate", ApiGroup: "menu", ServiceName: "core-api"},
			{Path: "/admin/menu/delete", Method: http.MethodPost, Description: "api.menuDelete", ApiGroup: "menu", ServiceName: "core-api"},
			{Path: "/admin/menu/list", Method: http.MethodPost, Description: "api.menuList", ApiGroup: "menu", ServiceName: "core-api"},
			{Path: "/admin/api/create", Method: http.MethodPost, Description: "api.apiCreate", ApiGroup: "api", ServiceName: "core-api"},
			{Path: "/admin/api/update", Method: http.MethodPost, Description: "api.apiUpdate", ApiGroup: "api", ServiceName: "core-api"},
			{Path: "/admin/api/delete", Method: http.MethodPost, Description: "api.apiDelete", ApiGroup: "api", ServiceName: "core-api"},
			{Path: "/admin/api/list", Method: http.MethodPost, Description: "api.apiList", ApiGroup: "api", ServiceName: "core-api"},
			{Path: "/admin/authority/menu/update", Method: http.MethodPost, Description: "api.authorityMenuUpdate", ApiGroup: "authority", ServiceName: "core-api"},
			{Path: "/admin/authority/menu/role", Method: http.MethodPost, Description: "api.authorityMenuRole", ApiGroup: "authority", ServiceName: "core-api"},
			{Path: "/admin/authority/api/update", Method: http.MethodPost, Description: "api.authorityApiUpdate", ApiGroup: "authority", ServiceName: "core-api"},
			{Path: "/admin/authority/api/role", Method: http.MethodPost, Description: "api.authorityApiRole", ApiGroup: "authority", ServiceName: "core-api"},
			{Path: "/admin/log/login/list", Method: http.MethodPost, Description: "api.loginLogList", ApiGroup: "log", ServiceName: "core-api"},
			{Path: "/admin/log/action/list", Method: http.MethodPost, Description: "api.actionLogList", ApiGroup: "log", ServiceName: "core-api"},
			{Path: "/admin/log/error/list", Method: http.MethodPost, Description: "api.errorLogList", ApiGroup: "log", ServiceName: "core-api"},
			{Path: "/admin/i18n/create", Method: http.MethodPost, Description: "api.i18nCreate", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/update", Method: http.MethodPost, Description: "api.i18nUpdate", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/updateByKey", Method: http.MethodPost, Description: "api.i18nUpdateByKey", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/delete", Method: http.MethodPost, Description: "api.i18nDelete", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/list", Method: http.MethodPost, Description: "api.i18nList", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/lang/create", Method: http.MethodPost, Description: "api.i18nLangCreate", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/lang/update", Method: http.MethodPost, Description: "api.i18nLangUpdate", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/lang/reorder", Method: http.MethodPost, Description: "api.i18nLangReorder", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/lang/delete", Method: http.MethodPost, Description: "api.i18nLangDelete", ApiGroup: "i18n", ServiceName: "core-api"},
			{Path: "/admin/i18n/lang/list", Method: http.MethodPost, Description: "api.i18nLangList", ApiGroup: "i18n", ServiceName: "core-api"},
		},
		I18N:      append(append(menuI18n(code), apiI18n(code)...), frontI18n(code)...),
		I18NLangs: langSeeds(),
	}
}
