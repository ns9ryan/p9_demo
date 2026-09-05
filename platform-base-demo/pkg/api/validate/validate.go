package validate

import (
	"errors"
	"net/http"
	"reflect"
	"strings"

	enLang "github.com/go-playground/locales/en"
	zhLang "github.com/go-playground/locales/zh_Hans"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

// Validator 参数校验器
type Validator struct {
	validator       *validator.Validate      // 校验器
	trans           map[string]ut.Translator // 翻译器
	defaultLanguage string                   // 默认语言
}

// New 创建参数校验器
func New(defaultLanguage string) (*Validator, error) {
	en := enLang.New()
	zh := zhLang.New()

	// 创建多语言翻译器
	uni := ut.New(zh, en, zh)
	enTrans, _ := uni.GetTranslator(langEn)
	zhTrans, _ := uni.GetTranslator(langZh)

	v := &Validator{
		validator: validator.New(),
		trans: map[string]ut.Translator{
			langEn: enTrans,
			langZh: zhTrans,
		},
		defaultLanguage: defaultLanguage,
	}

	// 使用JSON字段名生成校验错误信息
	v.validator.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		if name == "" {
			return field.Name
		}

		return name
	})

	// 注册英文默认校验翻译
	if err := enTranslations.RegisterDefaultTranslations(v.validator, enTrans); err != nil {
		return nil, err
	}

	// 注册中文默认校验翻译
	if err := zhTranslations.RegisterDefaultTranslations(v.validator, zhTrans); err != nil {
		return nil, err
	}

	return v, nil
}

// Validate 校验请求参数
func (v *Validator) Validate(r *http.Request, data any) error {
	// 校验请求参数
	if err := v.validator.Struct(data); err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			return err
		}

		// 根据请求语言选择校验翻译器
		lang := resolveLanguage(r.Header.Get("X-Lang"), v.defaultLanguage)
		trans := v.trans[lang]

		// 生成参数校验错误信息
		messages := make([]string, 0, len(validationErrors))
		for _, validationError := range validationErrors {
			messages = append(messages, validationError.Translate(trans))
		}

		return NewError(strings.Join(messages, " "))
	}

	return nil
}

// =====================================================================================================================

// RegisterValidation 注册自定义校验规则
func (v *Validator) RegisterValidation(tag string, fn validator.Func) error {
	return v.validator.RegisterValidation(tag, fn)
}

// RegisterValidationTranslation 注册自定义校验规则翻译
func (v *Validator) RegisterValidationTranslation(
	tag string,
	lang string,
	registerFn validator.RegisterTranslationsFunc,
	translationFn validator.TranslationFunc,
) error {
	lang, ok := lookupLanguage(lang)
	if !ok {
		return errors.New("不支持的校验语言")
	}

	trans := v.trans[lang]
	return v.validator.RegisterTranslation(
		tag,
		trans,
		registerFn,
		translationFn,
	)
}
