package response

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"strings"

	"oa.98ent.com/p9/common/xerr"
)

const errNumberRange = "wrong number range setting"

var (
	reFieldNotSet  = regexp.MustCompile(`^(?:field )?"([^"]+)" is not set$`)
	reTypeMismatch = regexp.MustCompile(`^type mismatch for field "([^"]+)"`)
	reValueOptions = regexp.MustCompile(`value "[^"]*" for field "([^"]+)" is not defined in options`)
	reFieldOptions = regexp.MustCompile(`^field "([^"]+)" not in options$`)
)

// asRequestError 将错误转换为xerr.Error
func asRequestError(err error) *xerr.Error {
	if err == nil {
		return nil
	}
	// 处理JSON语法错误
	var syn *json.SyntaxError
	var typ *json.UnmarshalTypeError
	if errors.As(err, &syn) || errors.As(err, &typ) {
		return xerr.BadRequest("common.invalidParam")
	}
	// 处理EOF错误
	if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
		return xerr.BadRequest("common.invalidParam")
	}
	// 处理其他错误
	msg := err.Error()
	// 处理字段未设置错误
	if m := reFieldNotSet.FindStringSubmatch(msg); len(m) == 2 {
		return xerr.BadRequestWith("common.paramRequired", map[string]any{"Field": m[1]})
	}
	// 处理类型不匹配错误
	if m := reTypeMismatch.FindStringSubmatch(msg); len(m) == 2 {
		return xerr.BadRequestWith("common.paramTypeMismatch", map[string]any{"Field": m[1]})
	}
	// 处理值选项错误
	if m := reValueOptions.FindStringSubmatch(msg); len(m) == 2 {
		return xerr.BadRequestWith("common.paramValueInvalidField", map[string]any{"Field": m[1]})
	}
	// 处理字段选项错误
	if m := reFieldOptions.FindStringSubmatch(msg); len(m) == 2 {
		return xerr.BadRequestWith("common.paramValueInvalidField", map[string]any{"Field": m[1]})
	}
	// 处理数字范围错误
	if msg == errNumberRange || strings.Contains(msg, errNumberRange) ||
		strings.Contains(msg, "is not defined in options") ||
		strings.Contains(msg, "value out of range") {
		return xerr.BadRequest("common.paramValueInvalid")
	}
	// 处理类型不匹配错误
	if msg == "unsupported type on setting field value" || msg == "type mismatch" {
		return xerr.BadRequest("common.invalidParam")
	}
	return nil
}
