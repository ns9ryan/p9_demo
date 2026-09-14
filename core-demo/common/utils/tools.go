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
		return NormalizeIP(strings.TrimSpace(strings.Split(x, ",")[0]))
	}
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return NormalizeIP(r.RemoteAddr)
	}
	return NormalizeIP(host)
}

// NormalizeIP 规范化 IP：去空白、去方括号，并把 IPv4-mapped 收成 IPv4。
func NormalizeIP(raw string) string {
	raw = strings.TrimSpace(raw)
	raw = strings.TrimPrefix(raw, "[")
	raw = strings.TrimSuffix(raw, "]")
	ip := net.ParseIP(raw)
	if ip == nil {
		return raw
	}
	if v4 := ip.To4(); v4 != nil {
		return v4.String()
	}
	return ip.String()
}

// NormalizeIPWhitelist 规范化白名单：单 IP 精确规范化，CIDR 校验后保留。非法条目返回 false。
func NormalizeIPWhitelist(list []string) ([]string, bool) {
	out := make([]string, 0, len(list))
	seen := make(map[string]struct{}, len(list))
	for _, raw := range list {
		s := strings.TrimSpace(raw)
		if s == "" {
			continue
		}
		norm, ok := normalizeWhitelistEntry(s)
		if !ok {
			return nil, false
		}
		if _, dup := seen[norm]; dup {
			continue
		}
		seen[norm] = struct{}{}
		out = append(out, norm)
	}
	return out, true
}

func normalizeWhitelistEntry(s string) (string, bool) {
	if strings.Contains(s, "/") {
		_, cidr, err := net.ParseCIDR(s)
		if err != nil {
			return "", false
		}
		return cidr.String(), true
	}
	n := NormalizeIP(s)
	if net.ParseIP(n) == nil {
		return "", false
	}
	return n, true
}

// IPAllowed 判断客户端 IP 是否命中白名单（精确 IP 或 CIDR）。空列表一律不通过。
func IPAllowed(clientIP string, whitelist []string) bool {
	if len(whitelist) == 0 {
		return false
	}
	ip := net.ParseIP(NormalizeIP(clientIP))
	if ip == nil {
		return false
	}
	for _, entry := range whitelist {
		entry = strings.TrimSpace(entry)
		if entry == "" {
			continue
		}
		if strings.Contains(entry, "/") {
			_, cidr, err := net.ParseCIDR(entry)
			if err != nil {
				continue
			}
			if cidr.Contains(ip) {
				return true
			}
			continue
		}
		if NormalizeIP(entry) == ip.String() {
			return true
		}
	}
	return false
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
