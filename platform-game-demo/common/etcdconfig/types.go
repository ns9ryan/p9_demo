package etcdconfig

// EtcdConfig ETCD连接配置
type EtcdConfig struct {
	Hosts []string `json:"hosts,optional" yaml:"hosts"`
	Key   string   `json:"key,optional" yaml:"key"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver   string `json:"driver,default=mysql"`
	DSN      string `json:"dsn"`
	MaxConns int    `json:"maxConns,default=10"`
	MaxIdle  int    `json:"maxIdle,default=5"`
}

// GrpcServerConfig gRPC服务器配置
type GrpcServerConfig struct {
	Address string `json:"address,optional"`
}

// SwaggerConfig Swagger配置
type SwaggerConfig struct {
	Protocol string `json:"protocol,default=https"`
}

// LoaderOptions 配置加载器选项
type LoaderOptions struct {
	// 默认的ETCD Hosts
	DefaultHosts []string
	// 默认的ETCD Key
	DefaultKey string
	// 本地配置文件路径
	ConfigFile string
	// 日志输出函数（可选）
	Logger func(format string, args ...interface{})
}

// DefaultLoaderOptions 返回默认选项
func DefaultLoaderOptions(configFile string) LoaderOptions {
	return LoaderOptions{
		DefaultHosts: []string{"127.0.0.1:2379"},
		DefaultKey:   "platform-game",
		ConfigFile:   configFile,
		Logger: func(format string, args ...interface{}) {
			// 默认不输出日志
		},
	}
}
