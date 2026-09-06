package kafka

import (
	"sync"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	config *Config
	mu     sync.Mutex
)

// SetConfig 设置 kafka 配置
func SetConfig(cfg *Config) {
	mu.Lock()
	defer mu.Unlock()
	config = cfg
	if config != nil {
		initProducer()
	}
}

// GetConfig 获取 kafka 配置
func GetConfig() *Config {
	mu.Lock()
	defer mu.Unlock()
	return config
}

// initProducer 初始化生产者
func initProducer() {
	if config == nil || len(config.Brokers) == 0 {
		logx.Infof("kafka brokers not configured, skipping producer pool initialization")
		return
	}
	if err := buildPool(); err != nil {
		logx.Errorf("failed to initialize kafka producer pool: %v", err)
		return
	}
	logx.Infof("kafka producer pool initialized")
}
