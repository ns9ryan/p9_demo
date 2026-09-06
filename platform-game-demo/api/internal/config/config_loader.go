package config

import (
	"log"

	"github.com/zeromicro/go-zero/core/conf"
)

// LoadConfig 加载配置，从本地YAML文件读取
func LoadConfig(configFile string) (*Config, error) {
	log.Printf("[ConfigLoader] 开始加载配置，配置文件: %s\n", configFile)

	var config Config
	// 使用 conf.MustLoad 加载配置，要求所有必需字段都有效
	conf.MustLoad(configFile, &config)

	// 调试：打印完整的 CoreRpc 配置
	log.Printf("[ConfigLoader] 调试信息: CoreRpc=%+v\n", config.CoreRpc)
	log.Printf("[ConfigLoader] 调试信息: CoreRpc.Target='%s' (长度=%d)\n", config.CoreRpc.Target, len(config.CoreRpc.Target))

	if config.CoreRpc.Target != "" {
		log.Printf("[ConfigLoader] ✓ Core RPC 配置已加载: Target=%s, Timeout=%d\n",
			config.CoreRpc.Target, config.CoreRpc.Timeout)
	} else {
		log.Printf("[ConfigLoader] ⚠️  Core RPC 未配置，鉴权功能不可用\n")
	}

	log.Printf("[ConfigLoader] ✓ 配置加载完成: Name=%s, Host=%s, Port=%d\n",
		config.Name, config.Host, config.Port)

	return &config, nil
}
