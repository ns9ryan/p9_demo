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

	// Platform Base RPC配置
	PlatformBaseRpc zrpc.RpcClientConf

	// 国际化配置
	I18n i18n.Config
}
