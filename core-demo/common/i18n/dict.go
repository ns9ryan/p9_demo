package i18n

import (
	"context"
	"strings"
	"sync"

	"oa.98ent.com/p9/core/common/modules/cache"
)

const (
	GroupMenu = "menu"
	GroupAPI  = "api"
)

type DictLoader func(ctx context.Context, group, lang string) (map[string]string, error)

var (
	dictLoader DictLoader
	dictMu     sync.RWMutex
)

// SetDictLoader 设置数据加载器
func SetDictLoader(l DictLoader) {
	dictMu.Lock()
	defer dictMu.Unlock()
	dictLoader = l
}

// cacheKey 缓存key
func cacheKey(group, lang string) string {
	return "i18n:" + group + ":" + lang
}

// Dict 获取group组的数据(缓存)
func Dict(ctx context.Context, group string) map[string]string {
	dictMu.RLock()
	loader := dictLoader
	dictMu.RUnlock()
	if loader == nil {
		return nil
	}
	lang := Lang(ctx)
	data, err := cache.TwoMinuteCache.GetC(cacheKey(group, lang), func(map[string]any) (any, error) {
		return loader(context.Background(), group, lang)
	}, nil, false)
	if err != nil || data == nil {
		return nil
	}
	m, ok := data.(map[string]string)
	if !ok {
		return nil
	}
	return m
}

// TG 获取翻译(数据库)。先查拼接后的完整 key（如 menu.route.dashboard），再查短 key。
func TG(ctx context.Context, group, key string) string {
	if key == "" {
		return key
	}
	if d := Dict(ctx, group); d != nil {
		if group != "" {
			if v := d[group+"."+key]; v != "" {
				return v
			}
		}
		if v := d[key]; v != "" {
			return v
		}
		return key
	}
	return T(ctx, key)
}

// Invalidate 清除缓存
func Invalidate(group, lang string) {
	cache.TwoMinuteCache.Cache.Delete(cacheKey(group, lang))
}

// InvalidateAll 清除所有缓存
func InvalidateAll() {
	for k := range cache.TwoMinuteCache.Cache.Items() {
		if strings.HasPrefix(k, "i18n:") {
			cache.TwoMinuteCache.Cache.Delete(k)
		}
	}
}
