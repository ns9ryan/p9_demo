package i18n

import (
	"context"
	"testing"
)

func TestCodeByPartnerMode(t *testing.T) {
	cases := map[string]string{
		"on":       CodeOperator,
		"ON":       CodeOperator,
		" on ":     CodeOperator,
		"off":      CodePlatform,
		"":         CodePlatform,
		"platform": CodePlatform,
	}
	for in, want := range cases {
		if got := CodeByPartnerMode(in); got != want {
			t.Fatalf("CodeByPartnerMode(%q)=%q want %q", in, got, want)
		}
	}
}

func TestParseLang(t *testing.T) {
	cases := map[string]string{
		"":               LangZH,
		"zh":             LangZH,
		"zh-CN":          LangZH,
		"zh-CN,en;q=0.8": LangZH,
		"en":             LangEN,
		"en-US":          LangEN,
		"en-GB":          LangEN,
		"ja-JP":          "ja-JP",
		"ja":             "ja",
		"ko-KR":          "ko-KR",
	}
	for in, want := range cases {
		if got := ParseLang(in); got != want {
			t.Fatalf("ParseLang(%q)=%q want %q", in, got, want)
		}
	}
}

func TestT(t *testing.T) {
	zh := WithLang(context.Background(), LangZH)
	en := WithLang(context.Background(), LangEN)
	if got := T(zh, LoginLogResultSuccess); got != "成功" {
		t.Fatalf("zh loginLog.resultSuccess=%q", got)
	}
	if got := T(en, LoginLogResultFail); got != "Failed" {
		t.Fatalf("en loginLog.resultFail=%q", got)
	}
	if got := T(zh, "运营"); got != "运营" {
		t.Fatalf("passthrough=%q", got)
	}
}

func TestTG(t *testing.T) {
	InvalidateAll()
	t.Cleanup(func() {
		SetDictLoader(nil)
		InvalidateAll()
	})
	SetDictLoader(func(_ context.Context, code, group, lang string) (map[string]string, error) {
		if group != GroupMenu {
			return map[string]string{}, nil
		}
		if lang == LangZH {
			return map[string]string{
				"menu.route.dashboard": "工作台",
				"legacy":               "旧",
			}, nil
		}
		return map[string]string{"menu.route.dashboard": "Dashboard"}, nil
	})
	zh := WithLang(context.Background(), LangZH)
	en := WithLang(context.Background(), LangEN)
	if got := TG(zh, "platform", GroupMenu, "route.dashboard"); got != "工作台" {
		t.Fatalf("zh=%q", got)
	}
	if got := TG(en, "platform", GroupMenu, "route.dashboard"); got != "Dashboard" {
		t.Fatalf("en=%q", got)
	}
	if got := TG(zh, "platform", GroupMenu, "legacy"); got != "旧" {
		t.Fatalf("short key=%q", got)
	}
	if got := TG(zh, "platform", GroupMenu, "missing"); got != "missing" {
		t.Fatalf("missing=%q", got)
	}
}

func TestTf(t *testing.T) {
	zh := WithLang(context.Background(), LangZH)
	en := WithLang(context.Background(), LangEN)
	data := map[string]any{"Field": "id"}
	if got := Tf(zh, ParamRequired, data); got != "缺少必填参数: id" {
		t.Fatalf("zh ParamRequired=%q", got)
	}
	if got := Tf(en, ParamRequired, data); got != "required parameter missing: id" {
		t.Fatalf("en ParamRequired=%q", got)
	}
	if got := Tf(zh, ParamTypeMismatch, data); got != "参数类型错误: id" {
		t.Fatalf("zh ParamTypeMismatch=%q", got)
	}
	if got := Tf(en, ParamTypeMismatch, data); got != "parameter type mismatch: id" {
		t.Fatalf("en ParamTypeMismatch=%q", got)
	}
}
