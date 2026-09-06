package catalog

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/core/rpc/coreclient"
)

// Register 注册Platform Base菜单和API目录
func Register(cli coreclient.Core) error {
	req := PlatformBaseReq()

	var last error

	// Core启动可能稍晚，失败时短暂重试
	for range 20 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

		_, last = cli.RegisterCatalog(ctx, req)
		cancel()

		if last == nil {
			return nil
		}

		time.Sleep(200 * time.Millisecond)
	}

	return fmt.Errorf("注册Platform Base目录失败: %w", last)
}
