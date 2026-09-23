package operator

import "oa.98ent.com/p9/operator-base/rpc/ent"

// Service operator 业务服务
type Service struct {
	db *ent.Client
}

// NewService 创建 operator 业务服务
func NewService(db *ent.Client) *Service {
	return &Service{
		db: db,
	}
}
