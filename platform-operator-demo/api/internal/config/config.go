// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"

	"oa.98ent.com/p9/common/i18n"
)

type Config struct {
	rest.RestConf

	// Platform Operator RPC配置
	PlatformOperatorRpc zrpc.RpcClientConf

	// Platform Base RPC配置
	PlatformBaseRpc zrpc.RpcClientConf

	// Platform Game RPC配置
	PlatformGameRpc zrpc.RpcClientConf

	// Core RPC配置
	CoreRpc zrpc.RpcClientConf

	// 国际化配置
	I18n i18n.Config
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
