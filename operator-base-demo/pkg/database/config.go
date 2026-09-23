package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	defaultMaxOpenConns    = 100             // 最大打开连接数
	defaultMaxIdleConns    = 20              // 最大空闲连接数
	defaultConnMaxIdleTime = 5 * time.Minute // 连接最大空闲时间
	defaultConnMaxLifetime = time.Hour       // 连接最大生命周期
	defaultPingTimeout     = 5 * time.Second // 数据库连接检查超时
)

// DatabaseConf 数据库配置
type DatabaseConf struct {
	Host     string // 数据库地址
	Port     int    // 数据库端口
	DBName   string // 数据库名称
	Username string // 数据库用户名
	Password string // 数据库密码
	SSLMode  string // SSL模式
}

// NewDriver 创建Ent数据库驱动
func (c DatabaseConf) NewDriver() (*entsql.Driver, error) {
	// 检查数据库配置
	if err := c.validate(); err != nil {
		return nil, err
	}

	// 创建数据库连接池
	db, err := sql.Open("pgx", c.dsn())
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接池失败: %w", err)
	}

	// 设置数据库连接池
	configurePool(db)

	// 检查数据库连接
	ctx, cancel := context.WithTimeout(context.Background(), defaultPingTimeout)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	// 创建Ent数据库驱动
	return entsql.OpenDB(dialect.Postgres, db), nil
}

// validate 检查数据库配置
func (c DatabaseConf) validate() error {
	if c.Host == "" {
		return errors.New("数据库地址不能为空")
	}

	if c.Port <= 0 || c.Port > 65535 {
		return fmt.Errorf("数据库端口无效: %d", c.Port)
	}

	if c.DBName == "" {
		return errors.New("数据库名称不能为空")
	}

	if c.Username == "" {
		return errors.New("数据库用户名不能为空")
	}

	if c.SSLMode == "" {
		return errors.New("数据库SSL模式不能为空")
	}

	return nil
}

// dsn 生成PostgreSQL连接地址
func (c DatabaseConf) dsn() string {
	user := url.User(c.Username)

	// 有密码时加入数据库密码
	if c.Password != "" {
		user = url.UserPassword(c.Username, c.Password)
	}

	dsn := &url.URL{
		Scheme: "postgres",
		User:   user,
		Host:   net.JoinHostPort(c.Host, strconv.Itoa(c.Port)),
		Path:   c.DBName,
	}

	// 设置SSL模式
	query := dsn.Query()
	query.Set("sslmode", c.SSLMode)
	dsn.RawQuery = query.Encode()

	return dsn.String()
}

// configurePool 设置数据库连接池
func configurePool(db *sql.DB) {
	db.SetMaxOpenConns(defaultMaxOpenConns)
	db.SetMaxIdleConns(defaultMaxIdleConns)
	db.SetConnMaxIdleTime(defaultConnMaxIdleTime)
	db.SetConnMaxLifetime(defaultConnMaxLifetime)
}
