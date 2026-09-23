package validate

import "strings"

const (
	langZh = "zh"
	langEn = "en"
)

var supportedLanguages = map[string]string{
	"zh":    langZh,
	"zh-cn": langZh,
	"en":    langEn,
	"en-us": langEn,
}

// resolveLanguage 获取支持的校验语言
func resolveLanguage(value, fallback string) string {
	// 优先使用请求语言
	if lang, ok := lookupLanguage(value); ok {
		return lang
	}

	// 请求语言不支持时使用默认语言
	if lang, ok := lookupLanguage(fallback); ok {
		return lang
	}

	// 默认语言也不支持时使用中文
	return langZh
}

// lookupLanguage 查找支持的校验语言
func lookupLanguage(value string) (string, bool) {
	value = strings.ToLower(strings.TrimSpace(value))
	lang, ok := supportedLanguages[value]
	return lang, ok
}
