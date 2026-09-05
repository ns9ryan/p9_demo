package i18n

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"path"
	"strings"

	goi18n "github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

// Translator 翻译器
type Translator struct {
	bundle          *goi18n.Bundle // 语言资源包
	defaultLanguage string         // 默认语言
}

// New 创建翻译器
func New(config Config, localeFS fs.FS) (*Translator, error) {
	// 解析默认语言
	defaultTag, err := language.Parse(config.DefaultLanguage)
	if err != nil {
		return nil, fmt.Errorf("无效的默认语言 %q: %w", config.DefaultLanguage, err)
	}

	// 创建语言资源包
	bundle := goi18n.NewBundle(defaultTag)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	// 加载所有JSON语言文件
	err = fs.WalkDir(localeFS, ".", func(filePath string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// 只加载JSON文件
		if entry.IsDir() || !strings.EqualFold(path.Ext(filePath), ".json") {
			return nil
		}

		if _, err := bundle.LoadMessageFileFS(localeFS, filePath); err != nil {
			return fmt.Errorf("加载语言文件 %q 失败: %w", filePath, err)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return &Translator{
		bundle:          bundle,
		defaultLanguage: config.DefaultLanguage,
	}, nil
}

// Translate 根据指定语言翻译消息
func (t *Translator) Translate(languageCode, key string) string {
	// 未指定语言时使用默认语言
	if languageCode == "" {
		languageCode = t.defaultLanguage
	}

	// 创建当前语言的本地化器
	localizer := goi18n.NewLocalizer(
		t.bundle,
		languageCode,
		t.defaultLanguage,
	)

	// 获取对应的翻译消息
	message, err := localizer.Localize(&goi18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil || message == "" {
		return key
	}

	return message
}

// Trans 根据上下文中的语言翻译消息
func (t *Translator) Trans(ctx context.Context, key string) string {
	return t.Translate(LanguageFromContext(ctx), key)
}
