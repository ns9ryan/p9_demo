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

	Swagger         SwaggerConfig `yaml:"Swagger" json:"Swagger"`
	PlatformGameRpc zrpc.RpcClientConf
	CoreRpc         zrpc.RpcClientConf
	I18n            i18n.Config `yaml:"I18n" json:"I18n"` // 国际化配置
	PartnerMode     string      `yaml:"PartnerMode" json:"PartnerMode,optional"`
}

// GetI18nCode 获取i18n code码
func (c *Config) GetI18nCode() string {
	if c.PartnerMode == "on" {
		return i18n.CodeOperator
	}
	return i18n.CodePlatform
}

// IsDebug 是否为调试模式
func (c *Config) IsDebug() bool {
	return c.Mode == service.DevMode || c.Mode == service.TestMode
}

// IsLocal 是否为本地开发环境
func (c *Config) IsLocal() bool {
	return c.Mode == service.DevMode
}
