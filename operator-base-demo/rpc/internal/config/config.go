package config

import (
	"github.com/zeromicro/go-zero/zrpc"

	"oa.98ent.com/p9/operator-base/pkg/database"
)

// Config RPC服务配置
type Config struct {
	zrpc.RpcServerConf

	// 数据库配置
	DatabaseConf database.DatabaseConf
}
