package utils

import (
	"net"
	"net/http"
	"strings"
)

// ClientIP 获取客户端IP
func ClientIP(r *http.Request) string {
	if r == nil {
		return ""
	}
	if x := r.Header.Get("X-Forwarded-For"); x != "" {
		return strings.TrimSpace(strings.Split(x, ",")[0])
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

// UserAgent 获取用户代理
func UserAgent(r *http.Request) string {
	if r == nil {
		return ""
	}
	return r.UserAgent()
}

// ResolveIDs 优先使用切片 ids，如果切片为空则
// 尝试将单个 id 包装成切片返回，均无有效数据则返回 nil
func ResolveIDs(id int64, ids []int64) []int64 {
	if len(ids) > 0 {
		return ids
	}
	if id > 0 {
		return []int64{id}
	}
	return nil
}

// GetEffectiveId 优先返回有效的 firstID，否则返回 secondID
func GetEffectiveId(firstID, secondID int64) int64 {
	if firstID > 0 {
		return firstID
	}
	return secondID
}
