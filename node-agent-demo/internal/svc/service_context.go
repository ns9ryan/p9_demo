package svc

import (
	"oa.98ent.com/p9/node-agent/internal/config"
	"oa.98ent.com/p9/node-agent/internal/websocket"
)

// ServiceContext 服务上下文
type ServiceContext struct {
	Config    config.Config     // 服务配置
	WebSocket *websocket.Client // 节点WebSocket客户端
}

// NewServiceContext 创建服务上下文
func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:    c,
		WebSocket: websocket.NewClient(c), // 节点WebSocket客户端
	}
}
