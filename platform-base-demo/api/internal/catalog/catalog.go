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

// registerRequest 创建Platform Base目录注册请求
func registerRequest() *coreclient.RegisterCatalogReq {
	return &coreclient.RegisterCatalogReq{
		Menus: menus(),
		Apis:  apis(),
		I18N:  platformBaseI18n(),
	}
}

// menus 返回Platform Base菜单目录
func menus() []*coreclient.RegisterMenuReq {
	return []*coreclient.RegisterMenuReq{
		// 基础数据
		{Name: "BaseData", Title: "menu.route.baseData", Path: "/base-data", MenuType: menuTypeDir, Sort: 20},

		// 时区管理
		{Name: "Timezone", Title: "menu.route.timezone", Path: "/base-data/timezone", MenuType: menuTypeMenu, Component: "base-data/timezone/index", ParentName: "BaseData", Sort: 21},
		{Name: "TimezoneUpdate", Title: "menu.route.timezoneUpdate", MenuType: menuTypeButton, Permission: "timezone:update", ParentName: "Timezone", Sort: 211},
		{Name: "TimezoneReorder", Title: "menu.route.timezoneReorder", MenuType: menuTypeButton, Permission: "timezone:reorder", ParentName: "Timezone", Sort: 212},

		// 货币管理
		{Name: "Currency", Title: "menu.route.currency", Path: "/base-data/currency", MenuType: menuTypeMenu, Component: "base-data/currency/index", ParentName: "BaseData", Sort: 22},
		{Name: "CurrencyUpdate", Title: "menu.route.currencyUpdate", MenuType: menuTypeButton, Permission: "currency:update", ParentName: "Currency", Sort: 221},
		{Name: "CurrencyReorder", Title: "menu.route.currencyReorder", MenuType: menuTypeButton, Permission: "currency:reorder", ParentName: "Currency", Sort: 222},

		// 国家地区管理
		{Name: "Region", Title: "menu.route.region", Path: "/base-data/region", MenuType: menuTypeMenu, Component: "base-data/region/index", ParentName: "BaseData", Sort: 23},
		{Name: "RegionUpdate", Title: "menu.route.regionUpdate", MenuType: menuTypeButton, Permission: "region:update", ParentName: "Region", Sort: 231},
		{Name: "RegionReorder", Title: "menu.route.regionReorder", MenuType: menuTypeButton, Permission: "region:reorder", ParentName: "Region", Sort: 232},
	}
}

// apis 返回Platform Base API目录
func apis() []*coreclient.CreateApiReq {
	return []*coreclient.CreateApiReq{
		// 时区
		{Path: "/admin/timezone/update", Method: http.MethodPost, Description: "api.timezoneUpdate", ApiGroup: "timezone", ServiceName: "platform-base-api"},
		{Path: "/admin/timezone/get", Method: http.MethodGet, Description: "api.timezoneGet", ApiGroup: "timezone", ServiceName: "platform-base-api"},
		{Path: "/admin/timezone/list", Method: http.MethodGet, Description: "api.timezoneList", ApiGroup: "timezone", ServiceName: "platform-base-api"},
		{Path: "/admin/timezone/list-all", Method: http.MethodGet, Description: "api.timezoneListAll", ApiGroup: "timezone", ServiceName: "platform-base-api"},
		{Path: "/admin/timezone/reorder", Method: http.MethodPost, Description: "api.timezoneReorder", ApiGroup: "timezone", ServiceName: "platform-base-api"},

		// 货币
		{Path: "/admin/currency/update", Method: http.MethodPost, Description: "api.currencyUpdate", ApiGroup: "currency", ServiceName: "platform-base-api"},
		{Path: "/admin/currency/get", Method: http.MethodGet, Description: "api.currencyGet", ApiGroup: "currency", ServiceName: "platform-base-api"},
		{Path: "/admin/currency/list", Method: http.MethodGet, Description: "api.currencyList", ApiGroup: "currency", ServiceName: "platform-base-api"},
		{Path: "/admin/currency/list-all", Method: http.MethodGet, Description: "api.currencyListAll", ApiGroup: "currency", ServiceName: "platform-base-api"},
		{Path: "/admin/currency/reorder", Method: http.MethodPost, Description: "api.currencyReorder", ApiGroup: "currency", ServiceName: "platform-base-api"},

		// 国家地区
		{Path: "/admin/region/update", Method: http.MethodPost, Description: "api.regionUpdate", ApiGroup: "region", ServiceName: "platform-base-api"},
		{Path: "/admin/region/get", Method: http.MethodGet, Description: "api.regionGet", ApiGroup: "region", ServiceName: "platform-base-api"},
		{Path: "/admin/region/list", Method: http.MethodGet, Description: "api.regionList", ApiGroup: "region", ServiceName: "platform-base-api"},
		{Path: "/admin/region/list-all", Method: http.MethodGet, Description: "api.regionListAll", ApiGroup: "region", ServiceName: "platform-base-api"},
		{Path: "/admin/region/reorder", Method: http.MethodPost, Description: "api.regionReorder", ApiGroup: "region", ServiceName: "platform-base-api"},
	}
}
