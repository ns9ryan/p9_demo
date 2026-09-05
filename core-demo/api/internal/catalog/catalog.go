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

func AdminReq() *coreclient.RegisterCatalogReq {
	return &coreclient.RegisterCatalogReq{
		Menus: []*coreclient.RegisterMenuReq{
			{Name: "Dashboard", Title: "route.dashboard", Path: "/dashboard", MenuType: menuTypeMenu, Component: "dashboard/index", Sort: 1},
			{Name: "System", Title: "route.system", Path: "/system", MenuType: menuTypeDir, Sort: 10},
			{Name: "User", Title: "route.user", Path: "/system/user", MenuType: menuTypeMenu, Component: "system/user/index", ParentName: "System", Sort: 11},
			{Name: "Role", Title: "route.role", Path: "/system/role", MenuType: menuTypeMenu, Component: "system/role/index", ParentName: "System", Sort: 12},
			{Name: "Menu", Title: "route.menu", Path: "/system/menu", MenuType: menuTypeMenu, Component: "system/menu/index", ParentName: "System", Sort: 13},
			{Name: "API", Title: "route.api", Path: "/system/api", MenuType: menuTypeMenu, Component: "system/api/index", ParentName: "System", Sort: 14},
			{Name: "Log", Title: "route.log", Path: "/system/log", MenuType: menuTypeDir, ParentName: "System", Sort: 15},
			{Name: "LoginLog", Title: "route.loginLog", Path: "/system/log/login", MenuType: menuTypeMenu, Component: "system/log/login/index", ParentName: "Log", Sort: 151},
			{Name: "ActionLog", Title: "route.actionLog", Path: "/system/log/action", MenuType: menuTypeMenu, Component: "system/log/action/index", ParentName: "Log", Sort: 152},
			{Name: "ErrorLog", Title: "route.errorLog", Path: "/system/log/error", MenuType: menuTypeMenu, Component: "system/log/error/index", ParentName: "Log", Sort: 153},
			{Name: "UserCreate", Title: "route.userCreate", MenuType: menuTypeButton, Permission: "user:create", ParentName: "User", Sort: 111},
			{Name: "RoleCreate", Title: "route.roleCreate", MenuType: menuTypeButton, Permission: "role:create", ParentName: "Role", Sort: 121},
			{Name: "MenuCreate", Title: "route.menuCreate", MenuType: menuTypeButton, Permission: "menu:create", ParentName: "Menu", Sort: 131},
			{Name: "APICreate", Title: "route.apiCreate", MenuType: menuTypeButton, Permission: "api:create", ParentName: "API", Sort: 141},
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
		},
	}
}
