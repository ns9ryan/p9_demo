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

// PlatformBaseReq 返回Platform Base菜单和API目录
func PlatformBaseReq() *coreclient.RegisterCatalogReq {
	return &coreclient.RegisterCatalogReq{
		Menus: []*coreclient.RegisterMenuReq{
			// 基础数据
			{
				Name:     "BaseData",
				Title:    "route.baseData",
				Path:     "/base-data",
				MenuType: menuTypeDir,
				Sort:     20,
			},

			// 语言管理
			{
				Name:       "Language",
				Title:      "route.language",
				Path:       "/base-data/language",
				MenuType:   menuTypeMenu,
				Component:  "base-data/language/index",
				ParentName: "BaseData",
				Sort:       21,
			},
			{
				Name:       "LanguageCreate",
				Title:      "route.languageCreate",
				MenuType:   menuTypeButton,
				Permission: "language:create",
				ParentName: "Language",
				Sort:       211,
			},
			{
				Name:       "LanguageUpdate",
				Title:      "route.languageUpdate",
				MenuType:   menuTypeButton,
				Permission: "language:update",
				ParentName: "Language",
				Sort:       212,
			},
			{
				Name:       "LanguageReorder",
				Title:      "route.languageReorder",
				MenuType:   menuTypeButton,
				Permission: "language:reorder",
				ParentName: "Language",
				Sort:       213,
			},

			// 时区管理
			{
				Name:       "Timezone",
				Title:      "route.timezone",
				Path:       "/base-data/timezone",
				MenuType:   menuTypeMenu,
				Component:  "base-data/timezone/index",
				ParentName: "BaseData",
				Sort:       22,
			},
			{
				Name:       "TimezoneCreate",
				Title:      "route.timezoneCreate",
				MenuType:   menuTypeButton,
				Permission: "timezone:create",
				ParentName: "Timezone",
				Sort:       221,
			},
			{
				Name:       "TimezoneUpdate",
				Title:      "route.timezoneUpdate",
				MenuType:   menuTypeButton,
				Permission: "timezone:update",
				ParentName: "Timezone",
				Sort:       222,
			},
			{
				Name:       "TimezoneReorder",
				Title:      "route.timezoneReorder",
				MenuType:   menuTypeButton,
				Permission: "timezone:reorder",
				ParentName: "Timezone",
				Sort:       223,
			},

			// 货币管理
			{
				Name:       "Currency",
				Title:      "route.currency",
				Path:       "/base-data/currency",
				MenuType:   menuTypeMenu,
				Component:  "base-data/currency/index",
				ParentName: "BaseData",
				Sort:       23,
			},
			{
				Name:       "CurrencyCreate",
				Title:      "route.currencyCreate",
				MenuType:   menuTypeButton,
				Permission: "currency:create",
				ParentName: "Currency",
				Sort:       231,
			},
			{
				Name:       "CurrencyUpdate",
				Title:      "route.currencyUpdate",
				MenuType:   menuTypeButton,
				Permission: "currency:update",
				ParentName: "Currency",
				Sort:       232,
			},
			{
				Name:       "CurrencyReorder",
				Title:      "route.currencyReorder",
				MenuType:   menuTypeButton,
				Permission: "currency:reorder",
				ParentName: "Currency",
				Sort:       233,
			},

			// 国家地区管理
			{
				Name:       "Region",
				Title:      "route.region",
				Path:       "/base-data/region",
				MenuType:   menuTypeMenu,
				Component:  "base-data/region/index",
				ParentName: "BaseData",
				Sort:       24,
			},
			{
				Name:       "RegionCreate",
				Title:      "route.regionCreate",
				MenuType:   menuTypeButton,
				Permission: "region:create",
				ParentName: "Region",
				Sort:       241,
			},
			{
				Name:       "RegionUpdate",
				Title:      "route.regionUpdate",
				MenuType:   menuTypeButton,
				Permission: "region:update",
				ParentName: "Region",
				Sort:       242,
			},
			{
				Name:       "RegionReorder",
				Title:      "route.regionReorder",
				MenuType:   menuTypeButton,
				Permission: "region:reorder",
				ParentName: "Region",
				Sort:       243,
			},
		},

		Apis: []*coreclient.CreateApiReq{
			// 语言
			{Path: "/language/create", Method: http.MethodPost, Description: "api.languageCreate", ApiGroup: "language", ServiceName: "platform-base-api"},
			{Path: "/language/update", Method: http.MethodPost, Description: "api.languageUpdate", ApiGroup: "language", ServiceName: "platform-base-api"},
			{Path: "/language/get", Method: http.MethodGet, Description: "api.languageGet", ApiGroup: "language", ServiceName: "platform-base-api"},
			{Path: "/language/list", Method: http.MethodGet, Description: "api.languageList", ApiGroup: "language", ServiceName: "platform-base-api"},
			{Path: "/language/list-all", Method: http.MethodGet, Description: "api.languageListAll", ApiGroup: "language", ServiceName: "platform-base-api"},
			{Path: "/language/reorder", Method: http.MethodPost, Description: "api.languageReorder", ApiGroup: "language", ServiceName: "platform-base-api"},

			// 时区
			{Path: "/timezone/create", Method: http.MethodPost, Description: "api.timezoneCreate", ApiGroup: "timezone", ServiceName: "platform-base-api"},
			{Path: "/timezone/update", Method: http.MethodPost, Description: "api.timezoneUpdate", ApiGroup: "timezone", ServiceName: "platform-base-api"},
			{Path: "/timezone/get", Method: http.MethodGet, Description: "api.timezoneGet", ApiGroup: "timezone", ServiceName: "platform-base-api"},
			{Path: "/timezone/list", Method: http.MethodGet, Description: "api.timezoneList", ApiGroup: "timezone", ServiceName: "platform-base-api"},
			{Path: "/timezone/list-all", Method: http.MethodGet, Description: "api.timezoneListAll", ApiGroup: "timezone", ServiceName: "platform-base-api"},
			{Path: "/timezone/reorder", Method: http.MethodPost, Description: "api.timezoneReorder", ApiGroup: "timezone", ServiceName: "platform-base-api"},

			// 货币
			{Path: "/currency/create", Method: http.MethodPost, Description: "api.currencyCreate", ApiGroup: "currency", ServiceName: "platform-base-api"},
			{Path: "/currency/update", Method: http.MethodPost, Description: "api.currencyUpdate", ApiGroup: "currency", ServiceName: "platform-base-api"},
			{Path: "/currency/get", Method: http.MethodGet, Description: "api.currencyGet", ApiGroup: "currency", ServiceName: "platform-base-api"},
			{Path: "/currency/list", Method: http.MethodGet, Description: "api.currencyList", ApiGroup: "currency", ServiceName: "platform-base-api"},
			{Path: "/currency/list-all", Method: http.MethodGet, Description: "api.currencyListAll", ApiGroup: "currency", ServiceName: "platform-base-api"},
			{Path: "/currency/reorder", Method: http.MethodPost, Description: "api.currencyReorder", ApiGroup: "currency", ServiceName: "platform-base-api"},

			// 国家地区
			{Path: "/region/create", Method: http.MethodPost, Description: "api.regionCreate", ApiGroup: "region", ServiceName: "platform-base-api"},
			{Path: "/region/update", Method: http.MethodPost, Description: "api.regionUpdate", ApiGroup: "region", ServiceName: "platform-base-api"},
			{Path: "/region/get", Method: http.MethodGet, Description: "api.regionGet", ApiGroup: "region", ServiceName: "platform-base-api"},
			{Path: "/region/list", Method: http.MethodGet, Description: "api.regionList", ApiGroup: "region", ServiceName: "platform-base-api"},
			{Path: "/region/list-all", Method: http.MethodGet, Description: "api.regionListAll", ApiGroup: "region", ServiceName: "platform-base-api"},
			{Path: "/region/reorder", Method: http.MethodPost, Description: "api.regionReorder", ApiGroup: "region", ServiceName: "platform-base-api"},
		},
	}
}
