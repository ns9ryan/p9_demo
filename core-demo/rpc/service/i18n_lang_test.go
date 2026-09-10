package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
)

func testLangSeeds() []CreateI18nLangReq {
	return []CreateI18nLangReq{
		{Lang: i18n.LangZH, Name: "简体中文", SortNo: 1},
		{Lang: i18n.LangHK, Name: "繁體中文", SortNo: 2},
		{Lang: i18n.LangEN, Name: "English", SortNo: 3},
	}
}

func TestUpsertI18nLangsSeedsAndKeepsExisting(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	list, err := d.ListEnabledI18nLangs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 || list[0].Lang != i18n.LangZH || list[0].Name != "简体中文" {
		t.Fatalf("seeded=%+v", list)
	}
	zh, err := d.i18nLangByID(ctx, list[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	one := int16(1)
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: zh.ID, Disabled: &one}); err != nil {
		t.Fatal(err)
	}
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	zh, err = d.i18nLangByID(ctx, zh.ID)
	if err != nil {
		t.Fatal(err)
	}
	if zh.Disabled != 1 {
		t.Fatalf("seed must not overwrite existing: %+v", zh)
	}
	if err := d.Client.I18nLang.UpdateOneID(zh.ID).SetName("").Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	zh, err = d.i18nLangByID(ctx, zh.ID)
	if err != nil {
		t.Fatal(err)
	}
	if zh.Name != "简体中文" {
		t.Fatalf("empty name should be backfilled: %+v", zh)
	}
	if err := d.Client.I18nLang.UpdateOneID(zh.ID).SetName("中文").Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	zh, err = d.i18nLangByID(ctx, zh.ID)
	if err != nil {
		t.Fatal(err)
	}
	if zh.Name != "中文" {
		t.Fatalf("non-empty name must not be overwritten: %+v", zh)
	}
}

func TestI18nLangCreateRules(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}

	if _, err := d.CreateI18nLang(ctx, CreateI18nLangReq{Lang: "  ja-JP  "}); err == nil {
		t.Fatal("expected name required")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nNameRequired {
		t.Fatalf("name required=%v", got.Message)
	}
	row, err := d.CreateI18nLang(ctx, CreateI18nLangReq{Lang: "  ja-JP  ", Name: " 日本語 "})
	if err != nil {
		t.Fatal(err)
	}
	if row.Lang != "ja-JP" || row.Name != "日本語" {
		t.Fatalf("created=%+v", row)
	}
	if _, err := d.CreateI18nLang(ctx, CreateI18nLangReq{Lang: "ja-JP", Name: "Japanese"}); err == nil {
		t.Fatal("expected duplicate lang")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nLangExists {
		t.Fatalf("duplicate=%v", got.Message)
	}

	one := int16(1)
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: row.ID, Disabled: &one}); err != nil {
		t.Fatal(err)
	}
	got, err := d.i18nLangByID(ctx, row.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Disabled != 1 {
		t.Fatalf("disabled=%+v", got)
	}
	if err := d.DeleteI18nLangs(ctx, []int64{row.ID}); err != nil {
		t.Fatal(err)
	}
}

func TestI18nLangBlockedWhenEntriesExist(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	row, err := d.CreateI18nLang(ctx, CreateI18nLangReq{Lang: "ja-JP", Name: "日本語"})
	if err != nil {
		t.Fatal(err)
	}
	next := "ko-KR"
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: row.ID, Lang: &next}); err != nil {
		t.Fatal(err)
	}
	got, err := d.i18nLangByID(ctx, row.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Lang != "ko-KR" {
		t.Fatalf("lang=%s", got.Lang)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "menu", TransKey: "route.demo", Lang: "ko-KR", Value: "demo"}); err != nil {
		t.Fatal(err)
	}
	same := "ko-KR"
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: row.ID, Lang: &same}); err != nil {
		t.Fatal(err)
	}
	back := "ja-JP"
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: got.ID, Lang: &back}); err == nil {
		t.Fatal("expected cannot change lang with entries")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nCannotChangeLangWithEntries {
		t.Fatalf("change lang=%v", got.Message)
	}
	if err := d.DeleteI18nLangs(ctx, []int64{row.ID}); err == nil {
		t.Fatal("expected cannot delete lang with entries")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nCannotDeleteLangWithEntries {
		t.Fatalf("delete lang=%v", got.Message)
	}
}

func TestCreateI18nRequiresSupportedLang(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "menu", TransKey: "route.demo", Lang: "ja-JP", Value: "demo"}); err == nil {
		t.Fatal("expected unsupported lang")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nLangNotSupported {
		t.Fatalf("unsupported=%v", got.Message)
	}
	if _, err := d.CreateI18nLang(ctx, CreateI18nLangReq{Lang: "ja-JP", Name: "日本語"}); err != nil {
		t.Fatal(err)
	}
	if _, err := d.CreateI18n(ctx, CreateI18nReq{I18nGroup: "menu", TransKey: "route.demo", Lang: "ja-JP", Value: "demo"}); err != nil {
		t.Fatal(err)
	}
	if err := d.UpdateI18nByKey(ctx, UpdateI18nByKeyReq{TransKey: "menu.route.demo", Data: map[string]string{"fr-FR": "x"}}); err == nil {
		t.Fatal("expected unsupported lang on by-key create")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nLangNotSupported {
		t.Fatalf("by-key=%v", got.Message)
	}
}

func TestI18nLangSortNo(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	late, err := d.CreateI18nLang(ctx, CreateI18nLangReq{Lang: "ja-JP", Name: "日本語", SortNo: 20})
	if err != nil {
		t.Fatal(err)
	}
	if late.SortNo != 20 {
		t.Fatalf("create sort=%d", late.SortNo)
	}
	list, _, err := d.ListI18nLangs(ctx, I18nLangListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 4 || list[0].Lang != i18n.LangZH || list[1].Lang != i18n.LangHK || list[2].Lang != i18n.LangEN || list[3].ID != late.ID {
		t.Fatalf("order=%+v", list)
	}
	zero := 0
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: late.ID, SortNo: &zero}); err != nil {
		t.Fatal(err)
	}
	got, err := d.i18nLangByID(ctx, late.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.SortNo != 0 {
		t.Fatalf("update 0 got=%d", got.SortNo)
	}
}

func TestReorderI18nLang(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.UpsertI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	list, _, err := d.ListI18nLangs(ctx, I18nLangListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 {
		t.Fatalf("seeded=%+v", list)
	}
	zh, hk, en := list[0], list[1], list[2]
	if zh.Lang != i18n.LangZH || hk.Lang != i18n.LangHK || en.Lang != i18n.LangEN {
		t.Fatalf("order=%+v", list)
	}

	if err := d.ReorderI18nLang(ctx, en.ID, zh.ID); err != nil {
		t.Fatal(err)
	}
	list, _, err = d.ListI18nLangs(ctx, I18nLangListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 || list[0].ID != en.ID || list[1].ID != zh.ID || list[2].ID != hk.ID {
		t.Fatalf("move en before zh: %+v", list)
	}
	for i, row := range list {
		if row.SortNo != i+1 {
			t.Fatalf("sort_no[%d]=%d", i, row.SortNo)
		}
	}

	if err := d.ReorderI18nLang(ctx, zh.ID, hk.ID); err != nil {
		t.Fatal(err)
	}
	list, _, err = d.ListI18nLangs(ctx, I18nLangListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if list[0].ID != en.ID || list[1].ID != hk.ID || list[2].ID != zh.ID {
		t.Fatalf("move zh after hk: %+v", list)
	}

	if err := d.ReorderI18nLang(ctx, en.ID, en.ID); err != nil {
		t.Fatal(err)
	}
	same, _, err := d.ListI18nLangs(ctx, I18nLangListReq{})
	if err != nil {
		t.Fatal(err)
	}
	if same[0].ID != list[0].ID || same[1].ID != list[1].ID || same[2].ID != list[2].ID {
		t.Fatalf("noop changed order: %+v", same)
	}

	if err := d.ReorderI18nLang(ctx, 0, zh.ID); err == nil {
		t.Fatal("expected invalid param")
	} else if got := xerr.AsError(err); got.Message != i18n.InvalidParam {
		t.Fatalf("invalid=%v", got.Message)
	}
	if err := d.ReorderI18nLang(ctx, en.ID, 99999); err == nil {
		t.Fatal("expected not found")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nLangNotFound {
		t.Fatalf("not found=%v", got.Message)
	}
}
