package entdb

import (
	"testing"

	"entgo.io/ent/dialect"
)

func TestNormalize(t *testing.T) {
	name, dialectName, err := Normalize("postgres")
	if err != nil || name != "pgx" || dialectName != dialect.Postgres {
		t.Fatalf("postgres: %s %s %v", name, dialectName, err)
	}
	name, dialectName, err = Normalize("PostgreSQL")
	if err != nil || name != "pgx" || dialectName != dialect.Postgres {
		t.Fatalf("postgresql: %s %s %v", name, dialectName, err)
	}
	name, dialectName, err = Normalize("mysql")
	if err != nil || name != "mysql" || dialectName != dialect.MySQL {
		t.Fatalf("mysql: %s %s %v", name, dialectName, err)
	}
	if _, _, err := Normalize("sqlite"); err == nil {
		t.Fatal("expected unsupported driver")
	}
}
