package etcdconfig

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"
)

// LoadEtcdConfig 加载ETCD配置信息
// 从本地配置文件中读取ETCD的连接信息（Hosts和Key）
func LoadEtcdConfig(configFile string, logger func(string, ...interface{})) (*EtcdConfig, error) {
	if logger == nil {
		logger = func(string, ...interface{}) {}
	}

	etcdConfig := &EtcdConfig{
		Hosts: []string{"127.0.0.1:2379"}, // 默认值
		Key:   "platform-game",            // 默认值
	}

	if configFile == "" {
		logger("[EtcdConfig] 配置文件为空，使用默认值")
		return etcdConfig, nil
	}

	data, err := os.ReadFile(configFile)
	if err != nil {
		logger("[EtcdConfig] 文件读取失败: %v，使用默认值", err)
		return etcdConfig, nil
	}

	logger("[EtcdConfig] 配置文件读取成功，大小: %d 字节", len(data))

	var tmpConfig struct {
		Etcd EtcdConfig `yaml:"etcd"`
	}
	if err := yaml.Unmarshal(data, &tmpConfig); err != nil {
		logger("[EtcdConfig] YAML解析失败: %v，使用默认值", err)
		return etcdConfig, nil
	}

	logger("[EtcdConfig] YAML解析成功")

	if len(tmpConfig.Etcd.Hosts) > 0 {
		etcdConfig.Hosts = tmpConfig.Etcd.Hosts
		logger("[EtcdConfig] ✓ 从配置读取ETCD Hosts: %v", etcdConfig.Hosts)
	} else {
		logger("[EtcdConfig] ⚠️  ETCD Hosts为空，使用默认值: %v", etcdConfig.Hosts)
	}

	if tmpConfig.Etcd.Key != "" {
		etcdConfig.Key = tmpConfig.Etcd.Key
		logger("[EtcdConfig] ✓ 从配置读取ETCD Key: %s", etcdConfig.Key)
	} else {
		logger("[EtcdConfig] ⚠️  ETCD Key为空，使用默认值: %s", etcdConfig.Key)
	}

	return etcdConfig, nil
}

// ReadConfigFromEtcd 从ETCD读取完整配置（返回原始YAML数据）
func ReadConfigFromEtcd(etcdHosts []string, etcdKey string, logger func(string, ...interface{})) ([]byte, error) {
	if logger == nil {
		logger = func(string, ...interface{}) {}
	}

	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   etcdHosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("ETCD客户端创建失败: %w", err)
	}
	defer cli.Close()

	logger("[EtcdLoader] 正在连接ETCD: %v", etcdHosts)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 尝试读取单个完整配置键
	singleKey := "/" + etcdKey
	logger("[EtcdLoader] 尝试查询单个配置键: %s", singleKey)
	resp, err := cli.Get(ctx, singleKey)
	if err != nil {
		return nil, fmt.Errorf("ETCD读取失败: %w", err)
	}

	if len(resp.Kvs) > 0 {
		logger("[EtcdLoader] ✓ 找到完整配置键: %s", singleKey)
		logger("[EtcdLoader] 配置长度: %d 字节", len(resp.Kvs[0].Value))
		return resp.Kvs[0].Value, nil
	}

	logger("[EtcdLoader] 未找到单个配置键，尝试查询文件夹结构")

	// 尝试读取文件夹形式的多个键
	prefix := "/" + etcdKey + "/"
	logger("[EtcdLoader] 尝试查询前缀: %s", prefix)
	resp, err = cli.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("ETCD读取失败: %w", err)
	}

	if len(resp.Kvs) == 0 {
		return nil, fmt.Errorf("ETCD中找不到配置（既无单个键 %s，也无文件夹前缀 %s）", singleKey, prefix)
	}

	logger("[EtcdLoader] ✓ 从文件夹找到 %d 个配置键", len(resp.Kvs))

	// 构建配置map（使用大小写不敏感的映射）
	configMap := make(map[string]interface{})

	for _, kv := range resp.Kvs {
		keyName := string(kv.Key)
		logger("[EtcdLoader] 读取键: %s", keyName)

		var data map[string]interface{}
		if err := json.Unmarshal(kv.Value, &data); err != nil {
			logger("[EtcdLoader] ⚠️  键 %s JSON解析失败: %v", keyName, err)
			continue
		}

		// 解析键名最后一个部分作为section名称
		sectionName := string(kv.Key[len(prefix):])

		// 规范化字段名（小写 -> 大写首字母）
		normalizedData := normalizeFieldNames(data)
		configMap[sectionName] = normalizedData
		logger("[EtcdLoader] ✓ 添加配置段: %s", sectionName)
	}

	// 转换为YAML格式
	yamlData, err := yaml.Marshal(configMap)
	if err != nil {
		return nil, fmt.Errorf("YAML序列化失败: %w", err)
	}

	logger("[EtcdLoader] ✓ 配置合并成功，总大小: %d 字节", len(yamlData))
	logger("[EtcdLoader] 合并后的YAML内容:\n%s", string(yamlData))

	// 如果有 api 部分，将其提升到顶级
	if apiData, ok := configMap["api"].(map[string]interface{}); ok {
		logger("[EtcdLoader] 检测到 api 配置段，正在提升字段到顶级...")

		// 将 api 数据合并到顶级
		for key, value := range apiData {
			if _, exists := configMap[key]; !exists {
				configMap[key] = value
				logger("[EtcdLoader] ✓ 提升字段: %s = %v", key, value)
			}
		}

		// 重新序列化
		yamlData, err = yaml.Marshal(configMap)
		if err != nil {
			return nil, fmt.Errorf("YAML序列化失败: %w", err)
		}
		logger("[EtcdLoader] ✓ 重新序列化完成，总大小: %d 字节", len(yamlData))
		logger("[EtcdLoader] 最终YAML内容:\n%s", string(yamlData))
	}

	return yamlData, nil
}

// normalizeFieldNames 规范化字段名（将小写首字母转换为大写）
func normalizeFieldNames(data map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, value := range data {
		// 将第一个字母大写
		normalizedKey := capitalizeFirstLetter(key)

		// 递归处理嵌套的 map
		if nestedMap, ok := value.(map[string]interface{}); ok {
			result[normalizedKey] = normalizeFieldNames(nestedMap)
		} else {
			result[normalizedKey] = value
		}
	}

	return result
}

// capitalizeFirstLetter 将字符串的第一个字母大写
func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	if len(s) == 1 {
		return strings.ToUpper(s)
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// ParseDatabaseConfig 解析数据库配置
func ParseDatabaseConfig(data map[string]interface{}) DatabaseConfig {
	dbConfig := DatabaseConfig{}
	if v, ok := data["driver"].(string); ok && v != "" {
		dbConfig.Driver = v
	}
	if v, ok := data["dsn"].(string); ok && v != "" {
		dbConfig.DSN = v
	}
	if v, ok := data["maxConns"].(float64); ok {
		dbConfig.MaxConns = int(v)
	}
	if v, ok := data["maxIdle"].(float64); ok {
		dbConfig.MaxIdle = int(v)
	}
	return dbConfig
}

// ParseGrpcConfig 解析gRPC配置
func ParseGrpcConfig(data map[string]interface{}) GrpcServerConfig {
	grpcConfig := GrpcServerConfig{}
	if v, ok := data["address"].(string); ok && v != "" {
		grpcConfig.Address = v
	}
	return grpcConfig
}

// ParseSwaggerConfig 解析Swagger配置
func ParseSwaggerConfig(data map[string]interface{}) SwaggerConfig {
	swaggerConfig := SwaggerConfig{}
	if v, ok := data["protocol"].(string); ok && v != "" {
		swaggerConfig.Protocol = v
	}
	return swaggerConfig
}
