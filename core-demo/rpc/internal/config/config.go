package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	DB struct {
		Driver string
		DSN    string
	}
	DataRedis struct {
		Host string
		Pass string
		DB   int
	}
	Jwt struct {
		AccessSecret  string
		AccessExpire  int64
		RefreshSecret string
		RefreshExpire int64
	}
	PartnerMode string
	APIPrefix   string
	InitToken   string
}

// IsDebug 是否为调试模式
func (c *Config) IsDebug() bool {
	return c.Mode == service.DevMode || c.Mode == service.TestMode
}

// IsPro 是否为生产模式
func (c *Config) IsPro() bool {
	return c.Mode == service.ProMode
}

// IsTest 是否为测试模式
func (c *Config) IsTest() bool {
	return c.Mode == service.TestMode
}

// IsDev 是否为开发模式
func (c *Config) IsDev() bool {
	return c.Mode == service.DevMode
}

// IsRt 是否为回归测试模式
func (c *Config) IsRt() bool {
	return c.Mode == service.RtMode
}

// IsPre 是否为预发布模式
func (c *Config) IsPre() bool {
	return c.Mode == service.PreMode
}
