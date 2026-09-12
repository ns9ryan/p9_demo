package utils

import (
	"net/http"
	"testing"
)

func TestClientIPNormalizesMappedIPv4(t *testing.T) {
	req := &http.Request{Header: http.Header{}}
	req.Header.Set("X-Forwarded-For", "::ffff:192.168.0.6")
	if got := ClientIP(req); got != "192.168.0.6" {
		t.Fatalf("xff mapped=%q", got)
	}

	req.Header.Set("X-Forwarded-For", "[::ffff:192.168.0.6], 10.0.0.1")
	if got := ClientIP(req); got != "192.168.0.6" {
		t.Fatalf("xff first hop=%q", got)
	}

	req.Header.Del("X-Forwarded-For")
	req.RemoteAddr = "[::ffff:192.168.0.6]:52344"
	if got := ClientIP(req); got != "192.168.0.6" {
		t.Fatalf("remote mapped=%q", got)
	}

	req.RemoteAddr = "127.0.0.1:8080"
	if got := ClientIP(req); got != "127.0.0.1" {
		t.Fatalf("remote v4=%q", got)
	}

	req.RemoteAddr = "[::1]:8080"
	if got := ClientIP(req); got != "::1" {
		t.Fatalf("remote v6=%q", got)
	}

	if ClientIP(nil) != "" {
		t.Fatal("nil request")
	}
}

func TestIPAllowed(t *testing.T) {
	if IPAllowed("192.168.0.6", nil) {
		t.Fatal("empty list")
	}
	if IPAllowed("192.168.0.6", []string{}) {
		t.Fatal("empty slice")
	}
	if !IPAllowed("192.168.0.6", []string{"192.168.0.6"}) {
		t.Fatal("exact")
	}
	if !IPAllowed("::ffff:192.168.0.6", []string{"192.168.0.6"}) {
		t.Fatal("mapped client")
	}
	if !IPAllowed("192.168.0.6", []string{"::ffff:192.168.0.6"}) {
		t.Fatal("mapped entry")
	}
	if IPAllowed("192.168.0.7", []string{"192.168.0.6"}) {
		t.Fatal("mismatch")
	}
	if !IPAllowed("10.1.2.3", []string{"10.0.0.0/8"}) {
		t.Fatal("cidr")
	}
	if IPAllowed("11.1.2.3", []string{"10.0.0.0/8"}) {
		t.Fatal("cidr miss")
	}
	if IPAllowed("not-an-ip", []string{"192.168.0.6"}) {
		t.Fatal("invalid client")
	}
}

func TestNormalizeIPWhitelist(t *testing.T) {
	got, ok := NormalizeIPWhitelist([]string{" ::ffff:192.168.0.6 ", "10.0.0.0/8", "192.168.0.6"})
	if !ok {
		t.Fatal("expected ok")
	}
	if len(got) != 2 || got[0] != "192.168.0.6" || got[1] != "10.0.0.0/8" {
		t.Fatalf("got=%v", got)
	}
	if _, ok := NormalizeIPWhitelist([]string{"not-an-ip"}); ok {
		t.Fatal("expected invalid")
	}
	if _, ok := NormalizeIPWhitelist([]string{"10.0.0.0/99"}); ok {
		t.Fatal("expected invalid cidr")
	}
	empty, ok := NormalizeIPWhitelist([]string{"", "  "})
	if !ok || len(empty) != 0 {
		t.Fatalf("blank entries: %v ok=%v", empty, ok)
	}
}
