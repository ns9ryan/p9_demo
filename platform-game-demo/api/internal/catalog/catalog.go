package catalog

import (
	"net/http"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

const (
	menuTypeDir    int32 = 0
	menuTypeMenu   int32 = 1
	menuTypeButton int32 = 2
)

// AdminReq 管理员请求(菜单相关种子数据)
func AdminReq(code string) *coreclient.RegisterCatalogReq {
	return &coreclient.RegisterCatalogReq{
		Menus: []*coreclient.RegisterMenuReq{
			{Name: "Game", Title: "menu.route.game", Path: "/gameManage", MenuType: menuTypeDir, Sort: 40},
			{Name: "GameCategory", Title: "menu.route.gameCategory", Path: "/gameManage/category", MenuType: menuTypeMenu, Component: "gameManage/category/index", ParentName: "Game", Sort: 1},
			{Name: "GameChannel", Title: "menu.route.gameChannel", Path: "/gameManage/channel", MenuType: menuTypeMenu, Component: "gameManage/channel/index", ParentName: "Game", Sort: 2},
			{Name: "GameManufacturer", Title: "menu.route.gameManufacturer", Path: "/gameManage/manufacturer", MenuType: menuTypeMenu, Component: "gameManage/manufacturer/index", ParentName: "Game", Sort: 3},
			{Name: "IndieGame", Title: "menu.route.gameIndie", Path: "/gameManage/indieGame", MenuType: menuTypeMenu, Component: "gameManage/indieGame/index", ParentName: "Game", Sort: 4},
			{Name: "GameCurrency", Title: "menu.route.gameCurrency", Path: "/gameManage/currency", MenuType: menuTypeMenu, Component: "gameManage/currency/index", ParentName: "Game", Sort: 5},
			{Name: "GameSyncList", Title: "menu.route.gameSyncList", Path: "/gameManage/syncList", MenuType: menuTypeMenu, Component: "gameManage/syncList/index", ParentName: "Game", Sort: 6},
		},
		Apis: []*coreclient.CreateApiReq{
			{Path: "/admin/game-category/list", Method: http.MethodGet, Description: "api.gameCategoryList", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-channel/list", Method: http.MethodGet, Description: "api.gameChannelList", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-provider/list", Method: http.MethodGet, Description: "api.gameProviderList", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game/list", Method: http.MethodGet, Description: "api.gameList", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-currency/list", Method: http.MethodGet, Description: "api.gameCurrencyList", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-sync-checkpoint/list", Method: http.MethodGet, Description: "api.gameSyncCheckpointList", ApiGroup: "game", ServiceName: "platform-game-api"},

			{Path: "/admin/game-sync-checkpoint/get", Method: http.MethodGet, Description: "api.gameSyncCheckpointGet", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game/get", Method: http.MethodGet, Description: "api.gameGet", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-category/get", Method: http.MethodGet, Description: "api.gameCategoryGet", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-channel/get", Method: http.MethodGet, Description: "api.gameChannelGet", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-provider/get", Method: http.MethodGet, Description: "api.gameProviderGet", ApiGroup: "game", ServiceName: "platform-game-api"},

			{Path: "/admin/game-category/update", Method: http.MethodPost, Description: "api.gameCategoryUpdate", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-channel/update", Method: http.MethodPost, Description: "api.gameChannelUpdate", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-provider/update", Method: http.MethodPost, Description: "api.gameProviderUpdate", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game/update", Method: http.MethodPost, Description: "api.gameUpdate", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-currency/update", Method: http.MethodPost, Description: "api.gameCurrencyUpdate", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/game-sync-checkpoint/update", Method: http.MethodPost, Description: "api.gameSyncCheckpointUpdate", ApiGroup: "game", ServiceName: "platform-game-api"},

			{Path: "/admin/sync/categories/preview", Method: http.MethodPost, Description: "api.syncCategoriesPreview", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/channel/preview", Method: http.MethodPost, Description: "api.syncChannelPreview", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/vendors/preview", Method: http.MethodPost, Description: "api.syncProviderPreview", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/games/preview", Method: http.MethodPost, Description: "api.syncGamePreview", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/currencies/preview", Method: http.MethodPost, Description: "api.syncCurrencyPreview", ApiGroup: "game", ServiceName: "platform-game-api"},

			{Path: "/admin/sync/categories/run", Method: http.MethodPost, Description: "api.syncCategoriesRun", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/channel/run", Method: http.MethodPost, Description: "api.syncChannelRun", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/vendors/run", Method: http.MethodPost, Description: "api.syncProviderRun", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/games/run", Method: http.MethodPost, Description: "api.syncGameRun", ApiGroup: "game", ServiceName: "platform-game-api"},
			{Path: "/admin/sync/currencies/run", Method: http.MethodPost, Description: "api.syncCurrencyRun", ApiGroup: "game", ServiceName: "platform-game-api"},
		},
		I18N:      append(append(menuI18n(code), apiI18n(code)...)),
		I18NLangs: langSeeds(),
	}
}

func AppendGameI18nItems(nameMap map[string]string) *coreclient.RegisterCatalogReq {
	var out []*coreclient.I18NItem
	for key, value := range nameMap {
		out = append(out,
			&coreclient.I18NItem{
				I18NCode:  i18n.CodePlatform,
				I18NGroup: i18nGroupGame,
				TransKey:  key + ".name",
				Lang:      i18n.LangZH,
				Value:     value,
			},
		)
	}
	return &coreclient.RegisterCatalogReq{
		I18N: out,
	}
}
