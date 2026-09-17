package config

// Config 节点代理配置
type Config struct {
	NodeCode   string       // 节点业务编码
	AuthSecret string       // 节点认证密钥
	Dispatch   DispatchConf // 调度中心连接配置
}

// DispatchConf 调度中心连接配置
type DispatchConf struct {
	WebSocketURL string // WebSocket连接地址
}
