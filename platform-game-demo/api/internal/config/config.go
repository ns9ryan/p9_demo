// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package config

import (
	"github.com/zeromicro/go-zero/rest"
	"oa.98ent.com/p9/platform-game/common/kafka"
)

type Config struct {
	rest.RestConf

	Swagger    SwaggerConfig    `yaml:"Swagger" json:"Swagger"`
	GrpcClient GrpcClientConfig `yaml:"GrpcClient" json:"GrpcClient"`
	Kafka      kafka.Config     `yaml:"Kafka" json:"Kafka"`
	// Core 服务 RPC 配置，用于鉴权
	CoreRpc GrpcClientConfig `yaml:"CoreRpc" json:"CoreRpc"`
}
