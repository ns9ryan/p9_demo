package casbinx

import (
	"testing"

	"github.com/casbin/casbin/v2/persist"
)

func TestAdapterIsBatchAdapter(t *testing.T) {
	var a persist.Adapter = NewAdapter(nil)
	if _, ok := a.(persist.BatchAdapter); !ok {
		t.Fatal("adapter must implement persist.BatchAdapter")
	}
}

func TestDomain(t *testing.T) {
	if Domain(nil) != "" {
		t.Fatal("nil should be empty domain")
	}
	code := "demo"
	if Domain(&code) != "demo" {
		t.Fatalf("got %q", Domain(&code))
	}
	empty := ""
	if Domain(&empty) != "" {
		t.Fatal("empty code should be empty domain")
	}
}
