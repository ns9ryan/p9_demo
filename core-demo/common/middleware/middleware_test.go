package middleware

import (
	"net/http"
	"testing"
)

func TestPreviewWriteDenied(t *testing.T) {
	allow := [][2]string{
		{http.MethodGet, "/admin/user/info"},
		{http.MethodGet, "/admin/user/perm"},
		{http.MethodGet, "/admin/operator/self"},
		{http.MethodGet, "/admin/menu/role"},
		{http.MethodHead, "/admin/user/detail"},
		{http.MethodOptions, "/admin/role/detail"},
		{http.MethodPost, "/admin/user/list"},
		{http.MethodPost, "/admin/role/list"},
		{http.MethodPost, "/admin/menu/list"},
		{http.MethodPost, "/admin/api/list"},
		{http.MethodPost, "/admin/authority/menu/role"},
		{http.MethodPost, "/admin/authority/api/role"},
	}
	deny := [][2]string{
		{http.MethodPost, "/admin/user/create"},
		{http.MethodPost, "/admin/user/update"},
		{http.MethodPost, "/admin/user/delete"},
		{http.MethodPost, "/admin/user/password"},
		{http.MethodPost, "/admin/user/password/self"},
		{http.MethodPost, "/admin/user/roles"},
		{http.MethodPost, "/admin/logout"},
		{http.MethodPost, "/admin/logout/all"},
		{http.MethodPost, "/admin/operator/update"},
		{http.MethodPost, "/admin/role/create"},
		{http.MethodPost, "/admin/menu/update"},
		{http.MethodPost, "/admin/authority/menu/update"},
		{http.MethodPost, "/admin/authority/api/update"},
		{http.MethodPut, "/admin/user/list"},
		{http.MethodDelete, "/admin/user/info"},
	}
	for _, c := range allow {
		if previewWriteDenied(c[0], c[1]) {
			t.Fatalf("allow %s %s", c[0], c[1])
		}
	}
	for _, c := range deny {
		if !previewWriteDenied(c[0], c[1]) {
			t.Fatalf("deny %s %s", c[0], c[1])
		}
	}
}
