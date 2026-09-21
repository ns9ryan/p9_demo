package catalog

import (
	"net/http"

	"oa.98ent.com/p9/core/rpc/coreclient"
)

const (
	menuTypeDir    int32 = 0 // 目录
	menuTypeMenu   int32 = 1 // 菜单
	menuTypeButton int32 = 2 // 按钮
)

// registerRequest 创建Platform Operator目录注册请求
func registerRequest() *coreclient.RegisterCatalogReq {
	return &coreclient.RegisterCatalogReq{
		Menus: menus(),
		Apis:  apis(),
		I18N:  platformOperatorI18n(),
	}
}

// menus 返回Platform Operator菜单目录
func menus() []*coreclient.RegisterMenuReq {
	return []*coreclient.RegisterMenuReq{
		// 分站管理
		{Name: "OperatorManagement", Title: "menu.route.operatorManagement", Path: "/operator", MenuType: menuTypeDir, Sort: 30},

		// 分站列表
		{Name: "OperatorList", Title: "menu.route.operatorList", Path: "/operator/list", MenuType: menuTypeMenu, Component: "operator/list/index", ParentName: "OperatorManagement", Sort: 31},
		{Name: "OperatorCreate", Title: "menu.route.operatorCreate", MenuType: menuTypeButton, Permission: "operator:create", ParentName: "OperatorList", Sort: 311},
		{Name: "OperatorUpdate", Title: "menu.route.operatorUpdate", MenuType: menuTypeButton, Permission: "operator:update", ParentName: "OperatorList", Sort: 312},
		{Name: "OperatorPublish", Title: "menu.route.operatorPublish", MenuType: menuTypeButton, Permission: "operator:publish", ParentName: "OperatorList", Sort: 313},
		{Name: "OperatorDelete", Title: "menu.route.operatorDelete", MenuType: menuTypeButton, Permission: "operator:delete", ParentName: "OperatorList", Sort: 314},

		// 域名管理
		{Name: "OperatorDomain", Title: "menu.route.operatorDomain", Path: "/operator/domain", MenuType: menuTypeMenu, Component: "operator/domain/index", ParentName: "OperatorManagement", Sort: 32},
		{Name: "OperatorDomainCreate", Title: "menu.route.operatorDomainCreate", MenuType: menuTypeButton, Permission: "operatorDomain:create", ParentName: "OperatorDomain", Sort: 321},
		{Name: "OperatorDomainUpdate", Title: "menu.route.operatorDomainUpdate", MenuType: menuTypeButton, Permission: "operatorDomain:update", ParentName: "OperatorDomain", Sort: 322},
		{Name: "OperatorDomainDelete", Title: "menu.route.operatorDomainDelete", MenuType: menuTypeButton, Permission: "operatorDomain:delete", ParentName: "OperatorDomain", Sort: 323},

		// 管理员账号
		{Name: "OperatorAdmin", Title: "menu.route.operatorAdmin", Path: "/operator/admin", MenuType: menuTypeMenu, Component: "operator/admin/index", ParentName: "OperatorManagement", Sort: 33},
		{Name: "OperatorAdminCreate", Title: "menu.route.operatorAdminCreate", MenuType: menuTypeButton, Permission: "operatorAdmin:create", ParentName: "OperatorAdmin", Sort: 331},
		{Name: "OperatorAdminUpdate", Title: "menu.route.operatorAdminUpdate", MenuType: menuTypeButton, Permission: "operatorAdmin:update", ParentName: "OperatorAdmin", Sort: 332},
		{Name: "OperatorAdminResetPassword", Title: "menu.route.operatorAdminResetPassword", MenuType: menuTypeButton, Permission: "operatorAdmin:resetPassword", ParentName: "OperatorAdmin", Sort: 333},
		{Name: "OperatorAdminUpdateStatus", Title: "menu.route.operatorAdminUpdateStatus", MenuType: menuTypeButton, Permission: "operatorAdmin:updateStatus", ParentName: "OperatorAdmin", Sort: 334},

		// 基础资源分配
		{Name: "BasicResourceAllocation", Title: "menu.route.basicResourceAllocation", Path: "/operator/basic-resource-allocation", MenuType: menuTypeMenu, Component: "operator/basic-resource-allocation/index", ParentName: "OperatorManagement", Sort: 34},
		{Name: "BasicResourceAllocationSave", Title: "menu.route.basicResourceAllocationSave", MenuType: menuTypeButton, Permission: "basicResourceAllocation:save", ParentName: "BasicResourceAllocation", Sort: 341},

		// 游戏资源分配
		{Name: "GameResourceAllocation", Title: "menu.route.gameResourceAllocation", Path: "/operator/game-resource-allocation", MenuType: menuTypeMenu, Component: "operator/game-resource-allocation/index", ParentName: "OperatorManagement", Sort: 35},
	}
}

// apis 返回Platform Operator API目录
func apis() []*coreclient.CreateApiReq {
	return []*coreclient.CreateApiReq{
		// 分站
		{Path: "/admin/operator/create", Method: http.MethodPost, Description: "api.operatorCreate", ApiGroup: "operator", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/update", Method: http.MethodPost, Description: "api.operatorUpdate", ApiGroup: "operator", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/get", Method: http.MethodGet, Description: "api.operatorGet", ApiGroup: "operator", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/list", Method: http.MethodGet, Description: "api.operatorList", ApiGroup: "operator", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/complete", Method: http.MethodPost, Description: "api.operatorComplete", ApiGroup: "operator", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/publish", Method: http.MethodPost, Description: "api.operatorPublish", ApiGroup: "operator", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/delete", Method: http.MethodPost, Description: "api.operatorDelete", ApiGroup: "operator", ServiceName: "platform-operator-api"},

		// 分站档案
		{Path: "/admin/operator/profile/create", Method: http.MethodPost, Description: "api.operatorProfileCreate", ApiGroup: "operator_profile", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/profile/update", Method: http.MethodPost, Description: "api.operatorProfileUpdate", ApiGroup: "operator_profile", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/profile/get", Method: http.MethodGet, Description: "api.operatorProfileGet", ApiGroup: "operator_profile", ServiceName: "platform-operator-api"},

		// 分站域名
		{Path: "/admin/operator/domain/create", Method: http.MethodPost, Description: "api.operatorDomainCreate", ApiGroup: "operator_domain", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/domain/update", Method: http.MethodPost, Description: "api.operatorDomainUpdate", ApiGroup: "operator_domain", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/domain/get", Method: http.MethodGet, Description: "api.operatorDomainGet", ApiGroup: "operator_domain", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/domain/list", Method: http.MethodGet, Description: "api.operatorDomainList", ApiGroup: "operator_domain", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/domain/delete", Method: http.MethodPost, Description: "api.operatorDomainDelete", ApiGroup: "operator_domain", ServiceName: "platform-operator-api"},

		// 分站管理员
		{Path: "/admin/operator/admin/list", Method: http.MethodGet, Description: "api.operatorAdminList", ApiGroup: "operator_admin", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/admin/create", Method: http.MethodPost, Description: "api.operatorAdminCreate", ApiGroup: "operator_admin", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/admin/update", Method: http.MethodPost, Description: "api.operatorAdminUpdate", ApiGroup: "operator_admin", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/admin/resetPassword", Method: http.MethodPost, Description: "api.operatorAdminResetPassword", ApiGroup: "operator_admin", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/admin/updateStatus", Method: http.MethodPost, Description: "api.operatorAdminUpdateStatus", ApiGroup: "operator_admin", ServiceName: "platform-operator-api"},

		// 基础资源分配
		{Path: "/admin/operator/basic-resource-allocation/list", Method: http.MethodGet, Description: "api.basicResourceAllocationList", ApiGroup: "basic_resource_allocation", ServiceName: "platform-operator-api"},

		// 语言分配
		{Path: "/admin/operator/language-allocation/list", Method: http.MethodGet, Description: "api.languageAllocationList", ApiGroup: "language_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/language-allocation/save", Method: http.MethodPost, Description: "api.languageAllocationSave", ApiGroup: "language_allocation", ServiceName: "platform-operator-api"},

		// 经营地区分配
		{Path: "/admin/operator/region-allocation/list", Method: http.MethodGet, Description: "api.regionAllocationList", ApiGroup: "region_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/region-allocation/save", Method: http.MethodPost, Description: "api.regionAllocationSave", ApiGroup: "region_allocation", ServiceName: "platform-operator-api"},

		// 代理子线路分配
		{Path: "/admin/operator/agent-line-allocation/list", Method: http.MethodGet, Description: "api.agentLineAllocationList", ApiGroup: "agent_line_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/agent-line-allocation/save", Method: http.MethodPost, Description: "api.agentLineAllocationSave", ApiGroup: "agent_line_allocation", ServiceName: "platform-operator-api"},

		// 分站游戏分配
		{Path: "/admin/operator/game-allocation/list", Method: http.MethodGet, Description: "api.operatorGameAllocationList", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game/list", Method: http.MethodGet, Description: "api.operatorGameList", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game/batch-update-status", Method: http.MethodPost, Description: "api.operatorGameBatchUpdateStatus", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-category/list", Method: http.MethodGet, Description: "api.operatorGameCategoryList", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-category/batch-update-status", Method: http.MethodPost, Description: "api.operatorGameCategoryBatchUpdateStatus", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-provider/list", Method: http.MethodGet, Description: "api.operatorGameProviderList", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-provider/batch-update-status", Method: http.MethodPost, Description: "api.operatorGameProviderBatchUpdateStatus", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-channel/list", Method: http.MethodGet, Description: "api.operatorGameChannelList", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-channel/batch-update-status", Method: http.MethodPost, Description: "api.operatorGameChannelBatchUpdateStatus", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},

		{Path: "/admin/operator/game/save-allocation", Method: http.MethodPost, Description: "api.operatorGameSaveAllocation", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-category/save-allocation", Method: http.MethodPost, Description: "api.operatorGameCategorySaveAllocation", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-provider/save-allocation", Method: http.MethodPost, Description: "api.operatorGameProviderSaveAllocation", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
		{Path: "/admin/operator/game-channel/save-allocation", Method: http.MethodPost, Description: "api.operatorGameChannelSaveAllocation", ApiGroup: "game_allocation", ServiceName: "platform-operator-api"},
	}
}
