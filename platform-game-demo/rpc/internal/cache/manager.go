package cache

import (
	"encoding/json"
	"sync"
	"time"
)

type CacheItem struct {
	Data      interface{}
	ExpiresAt time.Time
}

type Manager struct {
	store sync.Map
}

func NewManager() *Manager {
	return &Manager{}
}

// Set 设置缓存，ttl为过期时间（秒）
func (m *Manager) Set(key string, data interface{}, ttlSeconds int64) {
	expiresAt := time.Now().Add(time.Duration(ttlSeconds) * time.Second)
	m.store.Store(key, &CacheItem{
		Data:      data,
		ExpiresAt: expiresAt,
	})
}

// Get 获取缓存
func (m *Manager) Get(key string) (interface{}, bool) {
	val, ok := m.store.Load(key)
	if !ok {
		return nil, false
	}

	item := val.(*CacheItem)
	// 检查是否过期
	if time.Now().After(item.ExpiresAt) {
		m.store.Delete(key)
		return nil, false
	}

	return item.Data, true
}

// GetJSON 获取缓存并解析为 JSON
func (m *Manager) GetJSON(key string, v interface{}) bool {
	data, ok := m.Get(key)
	if !ok {
		return false
	}

	// 尝试 JSON 编解码
	bytes, err := json.Marshal(data)
	if err != nil {
		return false
	}
	err = json.Unmarshal(bytes, v)
	return err == nil
}

// Delete 删除缓存
func (m *Manager) Delete(key string) {
	m.store.Delete(key)
}

// Clear 清空所有缓存
func (m *Manager) Clear() {
	m.store.Range(func(key, _ interface{}) bool {
		m.store.Delete(key)
		return true
	})
}
