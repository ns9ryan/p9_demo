package service

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestEmptySlicesJSONArray(t *testing.T) {
	t.Run("menu tree", func(t *testing.T) {
		b, err := json.Marshal(buildMenuTree(nil, 0))
		if err != nil || string(b) != "[]" {
			t.Fatalf("tree=%s err=%v", b, err)
		}
	})
	t.Run("leaf children", func(t *testing.T) {
		n := MenuNode{ID: 1, Children: buildMenuTree(nil, 1)}
		b, err := json.Marshal(n)
		if err != nil || !strings.Contains(string(b), `"children":[]`) {
			t.Fatalf("node=%s err=%v", b, err)
		}
	})
	t.Run("role codes", func(t *testing.T) {
		u := UserPublic{ID: 1, RoleCodes: make([]string, 0), IPWhitelist: make([]string, 0)}
		b, err := json.Marshal(u)
		if err != nil || !strings.Contains(string(b), `"role_codes":[]`) {
			t.Fatalf("user=%s err=%v", b, err)
		}
		if !strings.Contains(string(b), `"ip_whitelist":[]`) {
			t.Fatalf("whitelist=%s", b)
		}
	})
	t.Run("public users", func(t *testing.T) {
		b, err := json.Marshal(PublicUsers(nil))
		if err != nil || string(b) != "[]" {
			t.Fatalf("users=%s err=%v", b, err)
		}
	})
}
