package config

import "github.com/zeromicro/go-zero/zrpc"

// Config 应用配置结构
type Config struct {
	zrpc.RpcServerConf

	// 数据库配置
	Database DatabaseConfig `json:"database,optional" yaml:"Database"`

	// gRPC 服务器地址（游戏供应商服务）
	GrpcServerAddr string `json:"grpcServerAddr,optional" yaml:"GrpcServerAddr"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `json:"driver,default=mysql,optional" yaml:"Driver"`
	DSN      string `json:"dsn,optional" yaml:"DSN"`
	MaxConns int    `json:"maxConns,default=10" yaml:"MaxConns"`
	MaxIdle  int    `json:"maxIdle,default=5" yaml:"MaxIdle"`
}
