package config

// GrpcClientConfig gRPC客户端配置（用于调用RPC服务）
type GrpcClientConfig struct {
	Target  string `yaml:"Target" json:"Target"`   // gRPC服务地址，如 "127.0.0.1:9001"
	Timeout int    `yaml:"Timeout" json:"Timeout"` // 超时时间（秒）
}

type SwaggerConfig struct {
	Protocol string `yaml:"Protocol" json:"Protocol"`
}
