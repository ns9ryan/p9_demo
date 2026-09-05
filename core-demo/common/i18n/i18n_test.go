package i18n

import (
	"context"
	"testing"
)

func TestParseLang(t *testing.T) {
	cases := map[string]string{
		"":               LangZH,
		"zh":             LangZH,
		"zh-CN":          LangZH,
		"zh-CN,en;q=0.8": LangZH,
		"en":             LangEN,
		"en-US":          LangEN,
		"en-GB":          LangEN,
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
	if got := T(zh, "route.dashboard"); got != "工作台" {
		t.Fatalf("zh route.dashboard=%q", got)
	}
	if got := T(zh, "route.log"); got != "日志管理" {
		t.Fatalf("zh route.log=%q", got)
	}
	if got := T(zh, "route.loginLog"); got != "登录日志" {
		t.Fatalf("zh route.loginLog=%q", got)
	}
	if got := T(zh, "route.actionLog"); got != "操作日志" {
		t.Fatalf("zh route.actionLog=%q", got)
	}
	if got := T(zh, "route.errorLog"); got != "错误日志" {
		t.Fatalf("zh route.errorLog=%q", got)
	}
	if got := T(zh, LoginLogResultSuccess); got != "成功" {
		t.Fatalf("zh loginLog.resultSuccess=%q", got)
	}
	if got := T(en, LoginLogResultFail); got != "Failed" {
		t.Fatalf("en loginLog.resultFail=%q", got)
	}
	if got := T(en, "route.dashboard"); got != "Dashboard" {
		t.Fatalf("en route.dashboard=%q", got)
	}
	if got := T(zh, "运营"); got != "运营" {
		t.Fatalf("passthrough=%q", got)
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
