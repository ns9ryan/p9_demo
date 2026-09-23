package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/core/rpc/ent/api"
	enti18n "oa.98ent.com/p9/core/rpc/ent/i18n"
	enti18nlang "oa.98ent.com/p9/core/rpc/ent/i18nlang"
	"oa.98ent.com/p9/core/rpc/ent/menu"
	"oa.98ent.com/p9/core/rpc/ent/sysinit"
	"oa.98ent.com/p9/core/rpc/model"
)

func TestRegisterCatalogFirstUpsertThenInsertOnly(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()

	if _, err := d.Client.Menu.Create().
		SetParentID(0).SetMenuType(model.MenuTypeDir).SetName("Dashboard").
		SetTitle("old-title").SetPath("/old").SetSort(9).
		Save(ctx); err != nil {
		t.Fatal(err)
	}

	firstMenus := []RegisterMenuReq{
		{Name: "Dashboard", Title: "menu.route.dashboard", MenuType: model.MenuTypeDir, Path: "/dashboard", Sort: 1},
		{Name: "User", Title: "menu.route.user", MenuType: model.MenuTypeMenu, Path: "/user", ParentName: "Dashboard", Sort: 2},
	}
	firstAPIs := []CreateAPIReq{
		{Method: "GET", Path: "/admin/user/list", Description: "api.userList", APIGroup: "user", ServiceName: "core-api"},
	}
	firstItems := []I18nItem{
		{I18nCode: "platform", I18nGroup: "menu", TransKey: "route.dashboard", Lang: i18n.LangZH, Value: "仪表盘"},
	}
	if err := d.RegisterCatalog(ctx, firstMenus, firstAPIs, firstItems, testLangSeeds()); err != nil {
		t.Fatal(err)
	}

	dash, err := d.Client.Menu.Query().Where(menu.NameEQ("Dashboard")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if dash.Title != "menu.route.dashboard" || dash.Path != "/dashboard" || dash.Sort != 1 {
		t.Fatalf("first run should upsert existing menu: %+v", dash)
	}
	user, err := d.Client.Menu.Query().Where(menu.NameEQ("User")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if user.ParentID != dash.ID {
		t.Fatalf("first run should link parent: parent=%d want=%d", user.ParentID, dash.ID)
	}

	keys, err := d.Client.SysInit.Query().Select(sysinit.FieldInitKey).All(ctx)
	if err != nil {
		t.Fatal(err)
	}
	got := map[string]bool{}
	for _, row := range keys {
		got[row.InitKey] = true
	}
	for _, k := range []string{InitKeyMenu, InitKeyAPI, InitKeyI18nLangs, InitKeyI18nDict} {
		if !got[k] {
			t.Fatalf("missing sys_init key %s: %+v", k, got)
		}
	}

	if err := d.Client.Menu.UpdateOne(dash).SetTitle("custom-title").SetPath("/custom").Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.Menu.UpdateOne(user).SetParentID(0).SetTitle("custom-user").Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.API.Update().Where(api.MethodEQ("GET"), api.PathEQ("/admin/user/list")).
		SetDescription("custom-api").Exec(ctx); err != nil {
		t.Fatal(err)
	}
	zh, err := d.Client.I18nLang.Query().Where(enti18nlang.LangEQ(i18n.LangZH)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if err := d.Client.I18nLang.UpdateOne(zh).SetName("中文定制").Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.Client.I18n.Update().
		Where(enti18n.I18nCodeEQ("platform"), enti18n.TransKeyEQ("menu.route.dashboard"), enti18n.LangEQ(i18n.LangZH)).
		SetValue("定制仪表盘").Exec(ctx); err != nil {
		t.Fatal(err)
	}

	secondMenus := []RegisterMenuReq{
		{Name: "Dashboard", Title: "should-not-overwrite", MenuType: model.MenuTypeDir, Path: "/dashboard", Sort: 1},
		{Name: "User", Title: "should-not-overwrite-user", MenuType: model.MenuTypeMenu, Path: "/user", ParentName: "Dashboard", Sort: 2},
		{Name: "Role", Title: "menu.route.role", MenuType: model.MenuTypeMenu, Path: "/role", ParentName: "Dashboard", Sort: 3},
	}
	secondAPIs := []CreateAPIReq{
		{Method: "GET", Path: "/admin/user/list", Description: "should-not-overwrite-api", APIGroup: "user", ServiceName: "core-api"},
		{Method: "POST", Path: "/admin/user/create", Description: "api.userCreate", APIGroup: "user", ServiceName: "core-api"},
	}
	secondItems := []I18nItem{
		{I18nCode: "platform", I18nGroup: "menu", TransKey: "route.dashboard", Lang: i18n.LangZH, Value: "should-not-overwrite-i18n"},
		{I18nCode: "platform", I18nGroup: "menu", TransKey: "route.role", Lang: i18n.LangZH, Value: "角色"},
	}
	secondLangs := append(testLangSeeds(), CreateI18nLangReq{Lang: "ja-JP", Name: "日本語", I18nKey: "lang.ja-JP", SortNo: 4})
	secondLangs[0].Name = "should-not-overwrite-lang"

	if err := d.RegisterCatalog(ctx, secondMenus, secondAPIs, secondItems, secondLangs); err != nil {
		t.Fatal(err)
	}

	dash, err = d.Client.Menu.Query().Where(menu.NameEQ("Dashboard")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if dash.Title != "custom-title" || dash.Path != "/custom" {
		t.Fatalf("initialized menu must not be updated: %+v", dash)
	}
	user, err = d.Client.Menu.Query().Where(menu.NameEQ("User")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if user.ParentID != 0 || user.Title != "custom-user" {
		t.Fatalf("initialized menu parent/title must not be updated: %+v", user)
	}
	role, err := d.Client.Menu.Query().Where(menu.NameEQ("Role")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if role.ParentID != dash.ID || role.Title != "menu.route.role" {
		t.Fatalf("new menu should be inserted and linked: %+v", role)
	}

	oldAPI, err := d.Client.API.Query().Where(api.MethodEQ("GET"), api.PathEQ("/admin/user/list")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if oldAPI.Description != "custom-api" {
		t.Fatalf("initialized api must not be updated: %+v", oldAPI)
	}
	newAPI, err := d.Client.API.Query().Where(api.MethodEQ("POST"), api.PathEQ("/admin/user/create")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if newAPI.Description != "api.userCreate" {
		t.Fatalf("new api should be inserted: %+v", newAPI)
	}

	zh, err = d.Client.I18nLang.Query().Where(enti18nlang.LangEQ(i18n.LangZH)).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if zh.Name != "中文定制" {
		t.Fatalf("initialized lang must not be updated: %+v", zh)
	}
	ja, err := d.Client.I18nLang.Query().Where(enti18nlang.LangEQ("ja-JP")).Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if ja.Name != "日本語" {
		t.Fatalf("new lang should be inserted: %+v", ja)
	}

	oldI18n, err := d.Client.I18n.Query().
		Where(enti18n.I18nCodeEQ("platform"), enti18n.TransKeyEQ("menu.route.dashboard"), enti18n.LangEQ(i18n.LangZH)).
		Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if oldI18n.Value != "定制仪表盘" {
		t.Fatalf("initialized i18n must not be updated: %+v", oldI18n)
	}
	newI18n, err := d.Client.I18n.Query().
		Where(enti18n.I18nCodeEQ("platform"), enti18n.TransKeyEQ("menu.route.role"), enti18n.LangEQ(i18n.LangZH)).
		Only(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if newI18n.Value != "角色" {
		t.Fatalf("new i18n should be inserted: %+v", newI18n)
	}
}
