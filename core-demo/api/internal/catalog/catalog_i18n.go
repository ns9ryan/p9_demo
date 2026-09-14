package catalog

import (
	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

func langSeeds() []*coreclient.CreateI18NLangReq {
	return []*coreclient.CreateI18NLangReq{
		{Lang: i18n.LangZH, Name: "简体中文", SortNo: 1},
		{Lang: i18n.LangHK, Name: "繁體中文", SortNo: 2},
		{Lang: i18n.LangEN, Name: "English", SortNo: 3},
	}
}

func menuI18n(i18nCode string) []*coreclient.I18NItem {
	var out []*coreclient.I18NItem
	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupMenu, TransKey: key, Lang: i18n.LangEN, Value: en},
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

func apiI18n(i18nCode string) []*coreclient.I18NItem {
	var out []*coreclient.I18NItem
	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupAPI, TransKey: key, Lang: i18n.LangEN, Value: en},
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
	add("api.userIpWhitelist", "修改用户IP白名单", "修改用戶IP白名單", "Update user IP whitelist")
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
	add("api.i18nLangReorder", "调整语言排序", "調整語言排序", "Reorder languages")
	add("api.i18nLangDelete", "删除支持的语言", "刪除支持的語言", "Delete language")
	add("api.i18nLangList", "支持的语言列表", "支持的語言列表", "Language list")
	return out
}

func frontI18n(i18nCode string) []*coreclient.I18NItem {
	var out []*coreclient.I18NItem
	add := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupFront, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupFront, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupFront, TransKey: key, Lang: i18n.LangEN, Value: en},
		)
	}
	addLogin := func(key, zh, hk, en string) {
		out = append(out,
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupLogin, TransKey: key, Lang: i18n.LangZH, Value: zh},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupLogin, TransKey: key, Lang: i18n.LangHK, Value: hk},
			&coreclient.I18NItem{I18NCode: i18nCode, I18NGroup: i18n.GroupLogin, TransKey: key, Lang: i18n.LangEN, Value: en},
		)
	}
	add("common.column.operations", "操作", "操作", "Operations")
	add("common.search.reset", "重置", "重置", "Reset")
	add("common.search.search", "搜索", "搜尋", "Search")
	add("common.placeholder.select", "请选择", "請選擇", "Please select")
	add("common.action.delete", "删除", "刪除", "Delete")
	add("common.column.createdAt", "创建时间", "建立時間", "Created At")
	add("common.column.index", "序号", "序號", "No.")
	add("common.okText.save", "保存", "儲存", "Save")
	add("common.message.saveSuccess", "保存成功", "儲存成功", "Saved successfully")
	add("common.action.edit", "编辑", "編輯", "Edit")
	add("common.message.deleteSuccess", "删除成功", "刪除成功", "Deleted successfully")
	add("common.delete.title", "确认删除", "確認刪除", "Confirm Delete")
	add("common.okText.create", "创建", "建立", "Create")
	add("common.message.createSuccess", "创建成功", "建立成功", "Created successfully")
	add("common.column.status", "状态", "狀態", "Status")
	add("common.toolbar.refresh", "刷新", "重新整理", "Refresh")
	add("common.toolbar.add", "添加", "新增", "Add")
	add("common.action.back", "返回", "返回", "Back")
	add("common.status.enabled", "启用", "啟用", "Enabled")
	add("common.message.enabled", "已启用", "已啟用", "Enabled")
	add("common.placeholder.input", "请输入", "請輸入", "Please input")
	add("common.search.statusDisabled", "关闭", "關閉", "Close")
	add("common.all", "全部", "全部", "All")
	add("common.action.view", "查看", "檢視", "View")
	add("common.column.updatedAt", "更新时间", "更新時間", "Updated At")
	add("common.message.submitSuccess", "提交成功", "提交成功", "Submitted successfully")
	add("common.status.disabled", "停用", "停用", "Disabled")
	add("common.delete.okText", "删除", "刪除", "Delete")
	add("common.message.updateSuccess", "更新成功", "更新成功", "Updated successfully")
	add("common.switch.on", "开", "開", "On")
	add("common.switch.off", "关", "關", "Off")
	add("common.yes", "是", "是", "Yes")
	add("common.no", "否", "否", "No")
	addLogin("login.brand.title", "智能科技平台", "智能科技平台", "Intelligent Technology Platform")
	addLogin("login.brand.subtitle", "ZHINENGKEJIPINGTAI", "ZHINENGKEJIPINGTAI", "ZHINENGKEJIPINGTAI")
	addLogin("login.form.title", "用户登录", "用戶登錄", "Sign in")
	addLogin("login.form.userName", "用户名", "用戶名", "Username")
	addLogin("login.form.userName.placeholder", "请输入用户名", "請輸入用戶名", "Enter username")
	addLogin("login.form.userName.errMsg", "请输入用户名", "請輸入用戶名", "Please enter username")
	addLogin("login.form.password", "密码", "密碼", "Password")
	addLogin("login.form.password.placeholder", "请输入密码", "請輸入密碼", "Enter password")
	addLogin("login.form.password.errMsg", "请输入密码", "請輸入密碼", "Please enter password")
	addLogin("login.form.rememberPassword", "记住我", "記住我", "Remember me")
	addLogin("login.form.forgetPassword", "忘记密码?", "忘記密碼?", "Forgot password?")
	addLogin("login.form.login", "登录", "登錄", "Sign in")
	addLogin("login.form.login.success", "登录成功", "登錄成功", "Signed in successfully")
	return out
}
