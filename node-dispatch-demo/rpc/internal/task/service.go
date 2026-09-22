package task

import (
	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/connection"
)

// Service 调度任务服务
type Service struct {
	db          *ent.Client         // Ent数据库客户端
	connections *connection.Manager // 节点连接管理器
}

// NewService 创建调度任务服务
func NewService(db *ent.Client, connections *connection.Manager) *Service {
	return &Service{
		db:          db,
		connections: connections,
	}
}
