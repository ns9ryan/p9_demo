package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/platform-base/pkg/database"
)

type Config struct {
	zrpc.RpcServerConf

	// 数据库配置
	DatabaseConf database.DatabaseConf
}
