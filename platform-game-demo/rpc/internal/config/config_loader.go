package config

import (
	"log"

	"github.com/zeromicro/go-zero/core/conf"
)

// LoadConfig 加载配置，从本地YAML文件读取
func LoadConfig(configFile string) (*Config, error) {
	log.Printf("[RpcConfigLoader] 开始加载配置，配置文件: %s\n", configFile)

	// 直接从本地YAML文件读取
	log.Printf("[RpcConfigLoader] 从本地配置文件读取: %s\n", configFile)
	var config Config
	if err := conf.Load(configFile, &config); err != nil {
		log.Printf("[RpcConfigLoader] ✗ 本地配置加载失败: %v\n", err)
		return nil, err
	}
	log.Printf("[RpcConfigLoader] ✓ 本地配置加载成功\n")

	log.Printf("[RpcConfigLoader] ✓ 配置加载完成: Name=%s, ListenOn=%s\n",
		config.Name, config.ListenOn)

	return &config, nil
}
