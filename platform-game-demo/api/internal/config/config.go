// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/platform-base/pkg/i18n"
)

type Config struct {
	rest.RestConf

	Swagger         SwaggerConfig `yaml:"Swagger" json:"Swagger"`
	PlatformGameRpc zrpc.RpcClientConf
	CoreRpc         zrpc.RpcClientConf
	I18n            i18n.Config `yaml:"I18n" json:"I18n"` // 国际化配置
}
