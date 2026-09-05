package service

import "testing"

func TestPageReqNormalize(t *testing.T) {
	t.Run("defaults", func(t *testing.T) {
		p := PageReq{}
		p.normalize(20)
		if p.Page != 1 || p.PageSize != 20 {
			t.Fatalf("got %+v", p)
		}
	})
	t.Run("cap 100", func(t *testing.T) {
		p := PageReq{Page: 2, PageSize: 500}
		p.normalize(20)
		if p.Page != 2 || p.PageSize != 100 {
			t.Fatalf("got %+v", p)
		}
	})
	t.Run("keep valid", func(t *testing.T) {
		p := PageReq{Page: 3, PageSize: 50}
		p.normalize(20)
		if p.Page != 3 || p.PageSize != 50 {
			t.Fatalf("got %+v", p)
		}
	})
}
