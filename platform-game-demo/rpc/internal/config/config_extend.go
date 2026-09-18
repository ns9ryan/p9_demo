package config

import (
	"github.com/zeromicro/go-zero/zrpc"
	"oa.98ent.com/p9/platform-base/pkg/database"
)

// Config 应用配置结构
type Config struct {
	zrpc.RpcServerConf

	DatabaseConf database.DatabaseConf

	// 游戏供应商 gRPC 服务器地址（游戏供应商服务）
	VendorGrpcServerAddr string `json:"vendorGrpcServerAddr,optional" yaml:"VendorGrpcServerAddr"`

	SyncBatchSize int `json:"syncBatchSize,optional,default=100" yaml:"SyncBatchSize"`

	Kafka struct {
		// Kafka brokers 列表
		Brokers []string `json:"brokers,optional" yaml:"Brokers"`
	} `json:"kafka,optional" yaml:"Kafka"`
}
