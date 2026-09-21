package response

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
)

func TestFromErrorMissingField(t *testing.T) {
	got := FromError(errors.New(`field "id" is not set`))
	if got.Status != http.StatusBadRequest || got.Message != i18n.ParamRequired {
		t.Fatalf("got %+v", got)
	}
	if got.Params["Field"] != "id" {
		t.Fatalf("Field=%v", got.Params["Field"])
	}
}

func TestFromErrorNestedFieldNotSet(t *testing.T) {
	got := FromError(errors.New(`"token" is not set`))
	if got.Status != http.StatusBadRequest || got.Message != i18n.ParamRequired {
		t.Fatalf("got %+v", got)
	}
	if got.Params["Field"] != "token" {
		t.Fatalf("Field=%v", got.Params["Field"])
	}
}

func TestFromErrorTypeMismatch(t *testing.T) {
	got := FromError(errors.New(`type mismatch for field "page", expect "int32", actual "string"`))
	if got.Status != http.StatusBadRequest || got.Message != i18n.ParamTypeMismatch {
		t.Fatalf("got %+v", got)
	}
	if got.Params["Field"] != "page" {
		t.Fatalf("Field=%v", got.Params["Field"])
	}
}

func TestFromErrorValueRange(t *testing.T) {
	got := FromError(errors.New("wrong number range setting"))
	if got.Status != http.StatusBadRequest || got.Message != i18n.ParamValueInvalid {
		t.Fatalf("got %+v", got)
	}
}

func TestFromErrorValueOptions(t *testing.T) {
	got := FromError(errors.New(`value "0" for field "status" is not defined in options "[1 2]"`))
	if got.Status != http.StatusBadRequest || got.Message != i18n.ParamValueInvalidField {
		t.Fatalf("got %+v", got)
	}
	if got.Params["Field"] != "status" {
		t.Fatalf("Field=%v", got.Params["Field"])
	}
}

func TestFromErrorValueOptionsNoField(t *testing.T) {
	got := FromError(errors.New(`error: value "x" is not defined in options "[a b]"`))
	if got.Status != http.StatusBadRequest || got.Message != i18n.ParamValueInvalid {
		t.Fatalf("got %+v", got)
	}
}

func TestLocalizeParamValueInvalid(t *testing.T) {
	zh := i18n.WithLang(context.Background(), i18n.LangZH)
	en := i18n.WithLang(context.Background(), i18n.LangEN)
	e := FromError(errors.New("wrong number range setting"))
	if got := localize(zh, e); got != "参数值错误" {
		t.Fatalf("zh=%q", got)
	}
	if got := localize(en, e); got != "invalid parameter value" {
		t.Fatalf("en=%q", got)
	}
	e = FromError(errors.New(`value "0" for field "id" is not defined in options "[1 2]"`))
	if got := localize(zh, e); got != "参数值错误: id" {
		t.Fatalf("zh field=%q", got)
	}
	if got := localize(en, e); got != "invalid parameter value: id" {
		t.Fatalf("en field=%q", got)
	}
}

func TestFromErrorJSONSyntax(t *testing.T) {
	got := FromError(&json.SyntaxError{})
	if got.Status != http.StatusBadRequest || got.Message != i18n.InvalidParam {
		t.Fatalf("got %+v", got)
	}
}

func TestFromErrorJSONType(t *testing.T) {
	got := FromError(&json.UnmarshalTypeError{Field: "username"})
	if got.Status != http.StatusBadRequest || got.Message != i18n.InvalidParam {
		t.Fatalf("got %+v", got)
	}
}

func TestFromErrorHidesInternal(t *testing.T) {
	got := FromError(errors.New("ERROR: duplicate key value violates unique constraint (SQLSTATE 23505)"))
	if got.Status != http.StatusInternalServerError || got.Message != i18n.InternalError {
		t.Fatalf("got %+v", got)
	}
	if got.Cause == nil || got.Stack == "" {
		t.Fatalf("missing cause/stack %+v", got)
	}
}

func TestFromErrorKeepsDomainError(t *testing.T) {
	orig := xerr.NotFound(i18n.UserNotFound)
	got := FromError(orig)
	if got.Status != http.StatusNotFound || got.Message != i18n.UserNotFound {
		t.Fatalf("got %+v", got)
	}
}

func TestLocalizeParamRequired(t *testing.T) {
	zh := i18n.WithLang(context.Background(), i18n.LangZH)
	en := i18n.WithLang(context.Background(), i18n.LangEN)
	e := FromError(errors.New(`field "id" is not set`))
	if got := localize(zh, e); got != "缺少必填参数: id" {
		t.Fatalf("zh=%q", got)
	}
	if got := localize(en, e); got != "required parameter missing: id" {
		t.Fatalf("en=%q", got)
	}
}
