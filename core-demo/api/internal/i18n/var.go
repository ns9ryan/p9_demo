package i18n

import "embed"

// FS API多语言资源
//
//go:embed locale/*.json
var LocaleFS embed.FS
