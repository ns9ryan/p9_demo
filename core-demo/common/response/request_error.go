package response

import (
	"encoding/json"
	"errors"
	"regexp"

	"oa.98ent.com/p9/core/common/i18n"
	"oa.98ent.com/p9/core/common/xerr"
)

var (
	reFieldNotSet  = regexp.MustCompile(`^(?:field )?"([^"]+)" is not set$`)
	reTypeMismatch = regexp.MustCompile(`^type mismatch for field "([^"]+)"`)
)

// asRequestError 将错误转换为xerr.Error
func asRequestError(err error) *xerr.Error {
	if err == nil {
		return nil
	}
	var syn *json.SyntaxError
	var typ *json.UnmarshalTypeError
	if errors.As(err, &syn) || errors.As(err, &typ) {
		return xerr.BadRequest(i18n.InvalidParam)
	}
	msg := err.Error()
	if m := reFieldNotSet.FindStringSubmatch(msg); len(m) == 2 {
		return xerr.BadRequestWith(i18n.ParamRequired, map[string]any{"Field": m[1]})
	}
	if m := reTypeMismatch.FindStringSubmatch(msg); len(m) == 2 {
		return xerr.BadRequestWith(i18n.ParamTypeMismatch, map[string]any{"Field": m[1]})
	}
	return nil
}
