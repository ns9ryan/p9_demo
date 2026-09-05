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
	id := int64(12)
	if Domain(&id) != "12" {
		t.Fatalf("got %q", Domain(&id))
	}
	if DomainID(0) != "" || DomainID(7) != "7" {
		t.Fatal("DomainID")
	}
}
