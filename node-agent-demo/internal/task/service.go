package task

import (
	operatorbaseservice "oa.98ent.com/p9/operator-base/rpc/client/operatorservice"
	platformoperatorservice "oa.98ent.com/p9/platform-operator/rpc/client/operatorservice"
)

// Service 调度任务执行服务
type Service struct {
	platformOperatorRpc platformoperatorservice.OperatorService // 总网分站RPC
	operatorBaseRpc     operatorbaseservice.OperatorService     // 当前节点分站基础RPC
}

// NewService 创建调度任务执行服务
func NewService(
	platformOperatorRpc platformoperatorservice.OperatorService,
	operatorBaseRpc operatorbaseservice.OperatorService,
) *Service {
	return &Service{
		platformOperatorRpc: platformOperatorRpc,
		operatorBaseRpc:     operatorBaseRpc,
	}
}
