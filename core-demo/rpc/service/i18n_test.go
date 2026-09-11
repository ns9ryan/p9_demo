package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
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
	} else if got := xerr.AsError(err); got.Message != i18n.I18nExists {
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
