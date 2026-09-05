package config

import "github.com/zeromicro/go-zero/zrpc"

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
