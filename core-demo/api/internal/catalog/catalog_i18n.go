package catalog

import (
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

func langSeeds() []*coreclient.CreateI18NLangReq {
	return []*coreclient.CreateI18NLangReq{
		{Lang: i18n.LangZH, Name: "简体中文", IsDefault: 1},
		{Lang: i18n.LangHK, Name: "繁體中文"},
		{Lang: i18n.LangEN, Name: "English"},
	}
}

func menuI18n() []*coreclient.I18NItem {
	var out []*coreclient.I18NItem
	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangEN, Value: en},
		)
	}
	add("menu.route.dashboard", "工作台", "工作台", "Dashboard")
	add("menu.route.system", "系统管理", "系統管理", "System")
	add("menu.route.user", "用户管理", "用戶管理", "Users")
	add("menu.route.role", "角色管理", "角色管理", "Roles")
	add("menu.route.menu", "菜单管理", "菜單管理", "Menus")
	add("menu.route.api", "接口管理", "接口管理", "APIs")
	add("menu.route.i18n", "多语言", "多語言", "I18n")
	add("menu.route.i18nEntry", "词条管理", "詞條管理", "Entries")
	add("menu.route.i18nLang", "语言管理", "語言管理", "Languages")
	add("menu.route.userCreate", "新建用户", "新建用戶", "Create user")
	add("menu.route.roleCreate", "新建角色", "新建角色", "Create role")
	add("menu.route.menuCreate", "新建菜单", "新建菜單", "Create menu")
	add("menu.route.apiCreate", "新建接口", "新建接口", "Create API")
	add("menu.route.i18nCreate", "新建多语言", "新建多語言", "Create i18n")
	add("menu.route.i18nLangCreate", "新建语言", "新建語言", "Create language")
	add("menu.route.log", "日志管理", "日志管理", "Logs")
	add("menu.route.loginLog", "登录日志", "登錄日志", "Login logs")
	add("menu.route.actionLog", "操作日志", "操作日志", "Action logs")
	add("menu.route.errorLog", "错误日志", "錯誤日志", "Error logs")
	return out
}

func apiI18n() []*coreclient.I18NItem {
	var out []*coreclient.I18NItem
	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangEN, Value: en},
		)
	}
	add("api.operatorSelf", "当前厅", "當前廳", "Current operator")
	add("api.operatorUpdate", "更新当前厅", "更新當前廳", "Update current operator")
	add("api.userCreate", "创建后台用户", "創建後台用戶", "Create admin user")
	add("api.userUpdate", "更新后台用户", "更新後台用戶", "Update admin user")
	add("api.userDelete", "删除后台用户", "刪除後台用戶", "Delete admin user")
	add("api.userList", "后台用户列表", "後台用戶列表", "Admin user list")
	add("api.userDetail", "后台用户详情", "後台用戶詳情", "Admin user detail")
	add("api.userPassword", "修改他人密码", "修改他人密碼", "Change another user's password")
	add("api.userRoles", "绑定用户角色", "綁定用戶角色", "Bind user roles")
	add("api.roleCreate", "创建角色", "創建角色", "Create role")
	add("api.roleUpdate", "更新角色", "更新角色", "Update role")
	add("api.roleDelete", "删除角色", "刪除角色", "Delete role")
	add("api.roleList", "角色列表", "角色列表", "Role list")
	add("api.roleDetail", "角色详情", "角色詳情", "Role detail")
	add("api.menuCreate", "创建菜单", "創建菜單", "Create menu")
	add("api.menuUpdate", "更新菜单", "更新菜單", "Update menu")
	add("api.menuDelete", "删除菜单", "刪除菜單", "Delete menu")
	add("api.menuList", "菜单目录", "菜單目錄", "Menu catalog")
	add("api.apiCreate", "创建API", "創建API", "Create API")
	add("api.apiUpdate", "更新API", "更新API", "Update API")
	add("api.apiDelete", "删除API", "刪除API", "Delete API")
	add("api.apiList", "API目录", "API目錄", "API catalog")
	add("api.authorityMenuUpdate", "分配菜单", "分配菜單", "Assign menus")
	add("api.authorityMenuRole", "查询角色菜单", "查詢角色菜單", "Query role menus")
	add("api.authorityApiUpdate", "分配API", "分配API", "Assign APIs")
	add("api.authorityApiRole", "查询角色API", "查詢角色API", "Query role APIs")
	add("api.loginLogList", "登录日志列表", "登錄日志列表", "Login log list")
	add("api.actionLogList", "操作日志列表", "操作日志列表", "Action log list")
	add("api.errorLogList", "错误日志列表", "錯誤日志列表", "Error log list")
	add("api.i18nCreate", "创建多语言", "創建多語言", "Create i18n")
	add("api.i18nUpdate", "更新多语言", "更新多語言", "Update i18n")
	add("api.i18nUpdateByKey", "按词条key更新多语言", "按詞條key更新多語言", "Update i18n by key")
	add("api.i18nDelete", "删除多语言", "刪除多語言", "Delete i18n")
	add("api.i18nList", "多语言列表", "多語言列表", "I18n list")
	add("api.i18nLangCreate", "创建支持的语言", "創建支持的語言", "Create language")
	add("api.i18nLangUpdate", "更新支持的语言", "更新支持的語言", "Update language")
	add("api.i18nLangDelete", "删除支持的语言", "刪除支持的語言", "Delete language")
	add("api.i18nLangList", "支持的语言列表", "支持的語言列表", "Language list")
	return out
}
