package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/retry"
	"oa.98ent.com/p9/core/rpc/coreclient"
)

// Register 注册Platform Base菜单和API目录
func Register(cli coreclient.Core) error {
	req := registerRequest()

	// Core启动可能稍晚，失败时短暂重试
	err := retry.Retry(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := cli.RegisterCatalog(ctx, req)
			return err
		},
		retry.RetryTimes(20),
		retry.RetryWithLinearBackoff(200*time.Millisecond),
	)
	if err != nil {
		return fmt.Errorf("注册Platform Base目录失败: %w", err)
	}

	return nil
}
