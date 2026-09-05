package entdb

import (
	"database/sql"
	"fmt"
	"strings"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Open 打开数据库
func Open(driver, dsn string) (dialect.Driver, error) {
	db, dialectName, err := OpenDB(driver, dsn)
	if err != nil {
		return nil, err
	}
	return entsql.OpenDB(dialectName, db), nil
}

// OpenDB 打开数据库
func OpenDB(driver, dsn string) (*sql.DB, string, error) {
	name, dialectName, err := Normalize(driver)
	if err != nil {
		return nil, "", err
	}
	db, err := sql.Open(name, dsn)
	if err != nil {
		return nil, "", err
	}
	if err := db.Ping(); err != nil {
		_ = db.Close()
		return nil, "", err
	}
	return db, dialectName, nil
}

// Normalize 规范化数据库驱动程序名称
func Normalize(driver string) (drvName, dialectName string, err error) {
	switch strings.ToLower(strings.TrimSpace(driver)) {
	case "postgres", "postgresql", "pg":
		return "pgx", dialect.Postgres, nil
	case "mysql":
		return "mysql", dialect.MySQL, nil
	default:
		return "", "", fmt.Errorf("demo-core: unsupported db driver %q", driver)
	}
}
