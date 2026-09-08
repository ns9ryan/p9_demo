package database

import (
	"os"
	"strconv"
	"testing"
)

// TestDatabaseConf_dsn 测试 PostgreSQL 连接地址生成
func TestDatabaseConf_dsn(t *testing.T) {
	tests := []struct {
		name string
		conf DatabaseConf
		want string
	}{
		{
			name: "普通配置",
			conf: DatabaseConf{
				Host:     "127.0.0.1",
				Port:     5432,
				DBName:   "platform_base",
				Username: "postgres",
				Password: "123456",
				SSLMode:  "disable",
			},
			want: "postgres://postgres:123456@127.0.0.1:5432/platform_base?sslmode=disable",
		},
		{
			name: "密码包含特殊字符",
			conf: DatabaseConf{
				Host:     "localhost",
				Port:     5432,
				DBName:   "platform_base",
				Username: "postgres",
				Password: "123@456",
				SSLMode:  "require",
			},
			want: "postgres://postgres:123%40456@localhost:5432/platform_base?sslmode=require",
		},
		{
			name: "IPv6地址",
			conf: DatabaseConf{
				Host:     "::1",
				Port:     5432,
				DBName:   "platform_base",
				Username: "postgres",
				Password: "123456",
				SSLMode:  "disable",
			},
			want: "postgres://postgres:123456@[::1]:5432/platform_base?sslmode=disable",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.conf.dsn()

			if got != tt.want {
				t.Errorf("dsn() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestDatabaseConf_NewDriver 测试数据库是否可以正常连接
//
// 该测试需要真实 PostgreSQL 环境
// 未配置 TEST_DATABASE_HOST 时自动跳过
func TestDatabaseConf_NewDriver(t *testing.T) {
	host := os.Getenv("TEST_DATABASE_HOST")
	if host == "" {
		t.Skip("未配置测试数据库，跳过数据库连接测试")
	}

	port := 5432
	if value := os.Getenv("TEST_DATABASE_PORT"); value != "" {
		parsed, err := strconv.Atoi(value)
		if err != nil {
			t.Fatalf("TEST_DATABASE_PORT 配置错误: %v", err)
		}

		port = parsed
	}

	conf := DatabaseConf{
		Host:     host,
		Port:     port,
		DBName:   os.Getenv("TEST_DATABASE_NAME"),
		Username: os.Getenv("TEST_DATABASE_USERNAME"),
		Password: os.Getenv("TEST_DATABASE_PASSWORD"),
		SSLMode:  "disable",
	}

	driver, err := conf.NewDriver()
	if err != nil {
		t.Fatalf("NewDriver() error = %v", err)
	}

	// 测试结束后关闭数据库连接
	if err := driver.Close(); err != nil {
		t.Errorf("关闭数据库连接失败: %v", err)
	}
}
