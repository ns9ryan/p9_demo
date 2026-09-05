package config

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	CoreRpc  zrpc.RpcClientConf
	CROSConf struct {
		Address string
	}
}
