package svc

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/node-dispatch/rpc/ent/migrate"
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
