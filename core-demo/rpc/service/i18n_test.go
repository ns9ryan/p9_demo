package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/common/xerr"
	coreI18n "oa.98ent.com/p9/core/common/i18n"
)

func TestI18nCodeIsolation(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}

	coreRow, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "front", TransKey: "common.hi", Lang: i18n.LangZH, Value: "你好"})
	if err != nil {
		t.Fatal(err)
	}
	if coreRow.I18nCode != i18n.CodePlatform {
		t.Fatalf("empty code default=%q", coreRow.I18nCode)
	}

	promo, err := d.CreateI18n(ctx, CreateI18nReq{
		I18nCode: "promo", I18nGroup: "front", TransKey: "common.hi", Lang: i18n.LangZH, Value: "优惠你好",
	})
	if err != nil {
		t.Fatal(err)
	}
	if promo.I18nCode != "promo" {
		t.Fatalf("promo code=%q", promo.I18nCode)
	}

	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "front", TransKey: "common.hi", Lang: i18n.LangZH, Value: "dup"}); err == nil {
		t.Fatal("expected duplicate core key")
	} else if got := xerr.AsError(err); got.Message != coreI18n.I18nExists {
		t.Fatalf("duplicate=%v", got.Message)
	}

	key := "front.common.hi"
	dictCore, err := d.GetI18nDict(ctx, i18n.CodePlatform, "front", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if dictCore[key] != "你好" {
		t.Fatalf("core dict=%v", dictCore)
	}
	dictPromo, err := d.GetI18nDict(ctx, "promo", "front", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if dictPromo[key] != "优惠你好" {
		t.Fatalf("promo dict=%v", dictPromo)
	}
	merged, err := d.GetI18nDict(ctx, "", "front", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if merged[key] != "优惠你好" {
		t.Fatalf("merged later-write should win, got=%v", merged)
	}

	coreList, total, err := d.ListI18ns(ctx, I18nListReq{PageReq: PageReq{Page: 1, PageSize: 50}, I18nCode: i18n.CodePlatform})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(coreList) != 1 || coreList[0].I18nCode != i18n.CodePlatform {
		t.Fatalf("core list=%+v total=%d", coreList, total)
	}

	if err := d.UpdateI18nByKey(ctx, UpdateI18nByKeyReq{
		TransKey: key, Data: map[string]string{i18n.LangEN: "Hi"},
	}); err != nil {
		t.Fatal(err)
	}
	dictCore, err = d.GetI18nDict(ctx, i18n.CodePlatform, "front", i18n.LangEN)
	if err != nil {
		t.Fatal(err)
	}
	if dictCore[key] != "Hi" {
		t.Fatalf("by-key platform=%v", dictCore)
	}
	dictPromo, err = d.GetI18nDict(ctx, "promo", "front", i18n.LangEN)
	if err != nil {
		t.Fatal(err)
	}
	if dictPromo[key] != "Hi" {
		t.Fatalf("by-key promo=%v", dictPromo)
	}

	if err := d.UpdateI18nByKey(ctx, UpdateI18nByKeyReq{
		I18nCode: "promo", TransKey: key, Data: map[string]string{i18n.LangEN: "PromoHi"},
	}); err != nil {
		t.Fatal(err)
	}
	dictPromo, err = d.GetI18nDict(ctx, "promo", "front", i18n.LangEN)
	if err != nil {
		t.Fatal(err)
	}
	if dictPromo[key] != "PromoHi" {
		t.Fatalf("by-key filter promo=%v", dictPromo)
	}
	dictCore, err = d.GetI18nDict(ctx, i18n.CodePlatform, "front", i18n.LangEN)
	if err != nil {
		t.Fatal(err)
	}
	if dictCore[key] != "Hi" {
		t.Fatalf("platform should stay, got=%v", dictCore)
	}

	if err := d.UpdateI18nByKey(ctx, UpdateI18nByKeyReq{
		TransKey: "front.only.platform", Data: map[string]string{i18n.LangZH: "仅平台"},
	}); err != nil {
		t.Fatal(err)
	}
	created, _, err := d.ListI18ns(ctx, I18nListReq{PageReq: PageReq{Page: 1, PageSize: 10}, TransKey: "front.only.platform"})
	if err != nil {
		t.Fatal(err)
	}
	if len(created) != 1 || created[0].I18nCode != i18n.CodePlatform {
		t.Fatalf("new key should be platform, got=%+v", created)
	}

	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "menu", TransKey: "route.x", Lang: i18n.LangZH, Value: "菜单X"}); err != nil {
		t.Fatal(err)
	}
	allGroups, err := d.GetI18nDict(ctx, i18n.CodePlatform, "", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if allGroups["front.common.hi"] != "你好" || allGroups["menu.route.x"] != "菜单X" {
		t.Fatalf("empty group should include all groups, got=%v", allGroups)
	}
}

func findFileItem(items []I18nFileItem, code, key string) (I18nFileItem, bool) {
	for _, it := range items {
		if it.I18nCode == code && it.I18nKey == key {
			return it, true
		}
	}
	return I18nFileItem{}, false
}

func TestExportImportI18n(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "api", TransKey: "userCreate", Lang: i18n.LangZH, Value: "创建后台用户"}); err != nil {
		t.Fatal(err)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "front", TransKey: "common.search.search", Lang: i18n.LangZH, Value: "搜索"}); err != nil {
		t.Fatal(err)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "front", TransKey: "common.search.search", Lang: i18n.LangEN, Value: "Search"}); err != nil {
		t.Fatal(err)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nCode: "promo", I18nGroup: "front", TransKey: "common.search.search", Lang: i18n.LangZH, Value: "优惠搜索"}); err != nil {
		t.Fatal(err)
	}

	if _, err := d.ExportI18ns(ctx, "", "", ""); err == nil {
		t.Fatal("expected lang required")
	} else if got := xerr.AsError(err); got.Message != coreI18n.I18nLangRequired {
		t.Fatalf("export empty lang=%v", got.Message)
	}

	zh, err := d.ExportI18ns(ctx, i18n.LangZH, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := findFileItem(zh, i18n.CodePlatform, "api.userCreate"); !ok {
		t.Fatalf("zh missing api, got=%+v", zh)
	}
	front, ok := findFileItem(zh, i18n.CodePlatform, "front.common.search.search")
	if !ok || front.I18nGroup != "front" || front.I18nValue != "搜索" {
		t.Fatalf("zh front=%+v ok=%v", front, ok)
	}
	if _, ok := findFileItem(zh, "promo", "front.common.search.search"); !ok {
		t.Fatalf("zh missing promo, got=%+v", zh)
	}
	for _, it := range zh {
		if it.I18nValue == "Search" {
			t.Fatalf("zh export should not include en, got=%+v", zh)
		}
	}

	apiOnly, err := d.ExportI18ns(ctx, i18n.LangZH, "", "api")
	if err != nil {
		t.Fatal(err)
	}
	if len(apiOnly) != 1 || apiOnly[0].I18nKey != "api.userCreate" || apiOnly[0].I18nGroup != "api" {
		t.Fatalf("api filter=%+v", apiOnly)
	}

	platform, err := d.ExportI18ns(ctx, i18n.LangZH, i18n.CodePlatform, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := findFileItem(platform, "promo", "front.common.search.search"); ok {
		t.Fatalf("platform filter should drop promo, got=%+v", platform)
	}

	hk, err := d.ExportI18ns(ctx, i18n.LangHK, "", "")
	if err != nil {
		t.Fatal(err)
	}
	if len(hk) != 0 {
		t.Fatalf("empty lang should export empty, got=%+v", hk)
	}

	if _, err := d.ImportI18ns(ctx, i18n.LangEN, nil); err == nil {
		t.Fatal("expected empty items")
	} else if got := xerr.AsError(err); got.Message != coreI18n.I18nDataRequired {
		t.Fatalf("empty items=%v", got.Message)
	}
	if _, err := d.ImportI18ns(ctx, "fr-FR", []I18nFileItem{{I18nGroup: "front", I18nKey: "front.x", I18nValue: "x"}}); err == nil {
		t.Fatal("expected unsupported lang")
	} else if got := xerr.AsError(err); got.Message != coreI18n.I18nLangNotSupported {
		t.Fatalf("unsupported lang=%v", got.Message)
	}
	if _, err := d.ImportI18ns(ctx, i18n.LangEN, []I18nFileItem{{I18nKey: ""}, {I18nKey: "   "}}); err == nil {
		t.Fatal("expected all skipped")
	} else if got := xerr.AsError(err); got.Message != coreI18n.I18nDataRequired {
		t.Fatalf("all skipped=%v", got.Message)
	}

	res, err := d.ImportI18ns(ctx, i18n.LangEN, []I18nFileItem{
		{I18nCode: i18n.CodePlatform, I18nGroup: "front", I18nKey: "front.common.search.search", I18nValue: "SearchEN"},
		{I18nCode: i18n.CodePlatform, I18nGroup: "api", I18nKey: "api.userCreate", I18nValue: "Create admin user"},
		{I18nGroup: "front", I18nKey: "front.common.okText.create", I18nValue: "Create"},
		{I18nKey: ""},
	})
	if err != nil {
		t.Fatal(err)
	}
	if res.Created != 2 || res.Updated != 1 || res.Skipped != 1 {
		t.Fatalf("import counts=%+v", res)
	}

	en, err := d.GetI18nDict(ctx, i18n.CodePlatform, "", i18n.LangEN)
	if err != nil {
		t.Fatal(err)
	}
	if en["front.common.search.search"] != "SearchEN" {
		t.Fatalf("updated en=%v", en)
	}
	if en["api.userCreate"] != "Create admin user" {
		t.Fatalf("created api en=%v", en)
	}
	if en["front.common.okText.create"] != "Create" {
		t.Fatalf("created front en=%v", en)
	}

	zhDict, err := d.GetI18nDict(ctx, i18n.CodePlatform, "", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if zhDict["front.common.search.search"] != "搜索" || zhDict["api.userCreate"] != "创建后台用户" {
		t.Fatalf("import must not delete other lang, got=%v", zhDict)
	}
	promoDict, err := d.GetI18nDict(ctx, "promo", "front", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if promoDict["front.common.search.search"] != "优惠搜索" {
		t.Fatalf("import must not delete other code, got=%v", promoDict)
	}
}

func TestDeleteI18nsByKey(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "front", TransKey: "common.hi", Lang: i18n.LangZH, Value: "你好"}); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdateI18nByKey(ctx, UpdateI18nByKeyReq{
		I18nCode: i18n.CodePlatform, I18nGroup: "front", TransKey: "front.common.hi",
		Data: map[string]string{i18n.LangEN: "Hi"},
	}); err != nil {
		t.Fatal(err)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{
		I18nCode: "promo", I18nGroup: "front", TransKey: "common.hi", Lang: i18n.LangZH, Value: "优惠你好",
	}); err != nil {
		t.Fatal(err)
	}

	if err := d.DeleteI18nsByKey(ctx, "", "", ""); err == nil {
		t.Fatal("expected empty key")
	} else if got := xerr.AsError(err); got.Message != coreI18n.I18nTransKeyRequired {
		t.Fatalf("empty key=%v", got.Message)
	}
	if err := d.DeleteI18nsByKey(ctx, i18n.CodePlatform, "front", "missing.key"); err == nil {
		t.Fatal("expected not found")
	} else if got := xerr.AsError(err); got.Message != coreI18n.I18nNotFound {
		t.Fatalf("missing=%v", got.Message)
	}

	if err := d.DeleteI18nsByKey(ctx, i18n.CodePlatform, "front", "common.hi"); err != nil {
		t.Fatal(err)
	}
	coreDict, err := d.GetI18nDict(ctx, i18n.CodePlatform, "front", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := coreDict["front.common.hi"]; ok {
		t.Fatalf("platform zh should be deleted, got=%v", coreDict)
	}
	enDict, err := d.GetI18nDict(ctx, i18n.CodePlatform, "front", i18n.LangEN)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := enDict["front.common.hi"]; ok {
		t.Fatalf("platform en should be deleted, got=%v", enDict)
	}
	promoDict, err := d.GetI18nDict(ctx, "promo", "front", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if promoDict["front.common.hi"] != "优惠你好" {
		t.Fatalf("promo should stay, got=%v", promoDict)
	}

	if err := d.DeleteI18nsByKey(ctx, "promo", "front", "front.common.hi"); err != nil {
		t.Fatal(err)
	}
	promoDict, err = d.GetI18nDict(ctx, "promo", "front", i18n.LangZH)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := promoDict["front.common.hi"]; ok {
		t.Fatalf("promo should be deleted, got=%v", promoDict)
	}
}
