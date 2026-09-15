package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/platform-base/pkg/database"
)

// Config 应用配置结构
type Config struct {
	zrpc.RpcServerConf

	DatabaseConf database.DatabaseConf

	// gRPC 服务器地址（游戏供应商服务）
	GrpcServerAddr string `json:"grpcServerAddr,optional" yaml:"GrpcServerAddr"`

	Kafka struct {
		// Kafka brokers 列表
		Brokers []string `json:"brokers,optional" yaml:"Brokers"`
	} `json:"kafka,optional" yaml:"Kafka"`
}
