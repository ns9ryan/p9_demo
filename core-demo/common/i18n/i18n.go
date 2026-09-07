package i18n

import (
	"context"
	"embed"
	"encoding/json"
	"strings"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

const (
	LangZH = "zh-CN"
	LangHK = "zh-HK"
	LangEN = "en-US"
)

type langKey struct{}

//go:embed locale/*.json
var localeFS embed.FS

var bundle *i18n.Bundle

func init() {
	b := i18n.NewBundle(language.MustParse(LangZH))
	mustAdd(b, language.MustParse(LangZH), "locale/zh.json")
	mustAdd(b, language.MustParse(LangEN), "locale/en.json")
	bundle = b
}

func mustAdd(b *i18n.Bundle, tag language.Tag, name string) {
	data, err := localeFS.ReadFile(name)
	if err != nil {
		panic(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		panic(err)
	}
	flat := map[string]string{}
	flatten("", raw, flat)
	msgs := make([]*i18n.Message, 0, len(flat))
	for id, other := range flat {
		msgs = append(msgs, &i18n.Message{ID: id, Other: other})
	}
	if err := b.AddMessages(tag, msgs...); err != nil {
		panic(err)
	}
}

func flatten(prefix string, raw map[string]any, out map[string]string) {
	for k, v := range raw {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch t := v.(type) {
		case string:
			out[key] = t
		case map[string]any:
			flatten(key, t, out)
		}
	}
}

// WithLang 设置语言到上下文
func WithLang(ctx context.Context, lang string) context.Context {
	return context.WithValue(ctx, langKey{}, lang)
}

// Lang 获取语言
func Lang(ctx context.Context) string {
	if ctx == nil {
		return LangZH
	}
	if v, ok := ctx.Value(langKey{}).(string); ok && v != "" {
		return v
	}
	return LangZH
}

// ParseLang 解析语言
func ParseLang(header string) string {
	if header == "" {
		return LangZH
	}
	first := strings.TrimSpace(strings.Split(header, ",")[0])
	first = strings.TrimSpace(strings.Split(first, ";")[0])
	if first == "" {
		return LangZH
	}
	lower := strings.ToLower(first)
	if strings.HasPrefix(lower, "en") {
		return LangEN
	}
	if strings.HasPrefix(lower, "zh") {
		if strings.HasPrefix(lower, "zh-hk") {
			return LangHK
		}
		return LangZH
	}
	return first
}

// T 获取翻译(本地)
func T(ctx context.Context, id string) string {
	return Tf(ctx, id, nil)
}

func Tf(ctx context.Context, id string, data map[string]any) string {
	if id == "" {
		return id
	}
	cfg := &i18n.LocalizeConfig{MessageID: id}
	if data != nil {
		cfg.TemplateData = data
	}
	s, err := i18n.NewLocalizer(bundle, Lang(ctx)).Localize(cfg)
	if err != nil || s == "" {
		return id
	}
	return s
}
