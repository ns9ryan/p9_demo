package validate

import (
	"reflect"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
)

const notBlankTag = "notblank"

// registerNotBlank 注册非空白字符串校验规则
func registerNotBlank(v *Validator) error {
	// 注册非空白字符串校验
	if err := v.RegisterValidation(notBlankTag, validateNotBlank); err != nil {
		return err
	}

	// 注册中文校验提示
	if err := v.RegisterValidationTranslation(
		notBlankTag,
		langZh,
		registerNotBlankZh,
		translateNotBlank,
	); err != nil {
		return err
	}

	// 注册英文校验提示
	if err := v.RegisterValidationTranslation(
		notBlankTag,
		langEn,
		registerNotBlankEn,
		translateNotBlank,
	); err != nil {
		return err
	}

	return nil
}

// validateNotBlank 校验字符串不能只包含空白字符
func validateNotBlank(fl validator.FieldLevel) bool {
	field := fl.Field()
	if field.Kind() != reflect.String {
		return false
	}

	return strings.TrimSpace(field.String()) != ""
}

// registerNotBlankZh 注册中文提示
func registerNotBlankZh(trans ut.Translator) error {
	return trans.Add(notBlankTag, "{0}不能为空或只包含空格", true)
}

// registerNotBlankEn 注册英文提示
func registerNotBlankEn(trans ut.Translator) error {
	return trans.Add(notBlankTag, "{0} cannot be blank", true)
}

// translateNotBlank 翻译非空白字符串校验错误
func translateNotBlank(trans ut.Translator, fieldError validator.FieldError) string {
	message, _ := trans.T(notBlankTag, fieldError.Field())
	return message
}
