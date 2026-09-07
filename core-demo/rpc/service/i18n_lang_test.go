package service

import (
	"context"
	"testing"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
)

func testLangSeeds() []CreateI18nLangReq {
	return []CreateI18nLangReq{
		{Lang: i18n.LangZH, Name: "简体中文", IsDefault: 1},
		{Lang: i18n.LangHK, Name: "繁體中文"},
		{Lang: i18n.LangEN, Name: "English"},
	}
}

func TestEnsureI18nLangsSeedsAndKeepsExisting(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.EnsureI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	list, err := d.ListEnabledI18nLangs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 3 || list[0].Lang != i18n.LangZH || list[0].IsDefault != 1 || list[0].Name != "简体中文" {
		t.Fatalf("seeded=%+v", list)
	}
	zh, err := d.i18nLangByID(ctx, list[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	one := int16(1)
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: zh.ID, Disabled: &one}); err == nil {
		t.Fatal("expected cannot disable default")
	}
	if err := d.Client.I18nLang.UpdateOneID(zh.ID).SetDisabled(1).SetIsDefault(0).Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.EnsureI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	zh, err = d.i18nLangByID(ctx, zh.ID)
	if err != nil {
		t.Fatal(err)
	}
	if zh.Disabled != 1 || zh.IsDefault != 0 {
		t.Fatalf("seed must not overwrite existing: %+v", zh)
	}
	if err := d.Client.I18nLang.UpdateOneID(zh.ID).SetName("").Exec(ctx); err != nil {
		t.Fatal(err)
	}
	if err := d.EnsureI18nLangs(ctx, testLangSeeds()); err != nil {
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
	if err := d.EnsureI18nLangs(ctx, testLangSeeds()); err != nil {
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

func TestI18nLangDefaultRules(t *testing.T) {
	d := testDeps(t, ModeOff)
	ctx := context.Background()
	if err := d.EnsureI18nLangs(ctx, testLangSeeds()); err != nil {
		t.Fatal(err)
	}
	enabled, err := d.ListEnabledI18nLangs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	zh := enabled[0]
	other := enabled[1]
	if other.IsDefault == 1 {
		t.Fatalf("expected only one default, second=%+v", other)
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

	zero := int16(0)
	one := int16(1)
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: zh.ID, IsDefault: &zero}); err == nil {
		t.Fatal("expected cannot unset default")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nCannotDisableDefault {
		t.Fatalf("unset default=%v", got.Message)
	}
	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: zh.ID, Disabled: &one}); err == nil {
		t.Fatal("expected cannot disable default")
	}
	if err := d.DeleteI18nLangs(ctx, []int64{zh.ID}); err == nil {
		t.Fatal("expected cannot delete default")
	} else if got := xerr.AsError(err); got.Message != i18n.I18nCannotDeleteDefault {
		t.Fatalf("delete default=%v", got.Message)
	}

	if err := d.UpdateI18nLang(ctx, UpdateI18nLangReq{ID: other.ID, IsDefault: &one}); err != nil {
		t.Fatal(err)
	}
	zh, err = d.i18nLangByID(ctx, zh.ID)
	if err != nil {
		t.Fatal(err)
	}
	other, err = d.i18nLangByID(ctx, other.ID)
	if err != nil {
		t.Fatal(err)
	}
	if zh.IsDefault != 0 || other.IsDefault != 1 {
		t.Fatalf("switch default zh=%+v other=%+v", zh, other)
	}

	list, err := d.ListEnabledI18nLangs(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if list[0].ID != other.ID || list[0].IsDefault != 1 {
		t.Fatalf("enabled order %+v", list)
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
