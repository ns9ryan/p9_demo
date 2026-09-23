package svc

import (
	"context"

	"oa.98ent.com/p9/operator-base/rpc/ent/migrate"

	"github.com/zeromicro/go-zero/core/logx"
)

// MustMigrate 执行数据库自动迁移
func (s *ServiceContext) MustMigrate() {
	ctx := context.Background()

	// 根据 Ent Schema 自动创建或更新数据库结构
	logx.Must(
		s.DB.Schema.Create(
			ctx,
			migrate.WithForeignKeys(false), // 不创建数据库外键
			migrate.WithDropIndex(true),    // 允许删除废弃索引
		),
	)
}
