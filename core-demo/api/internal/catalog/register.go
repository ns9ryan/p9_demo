package catalog

import (
	"context"
	"fmt"
	"time"

	"github.com/duke-git/lancet/v2/retry"
	"oa.98ent.com/p9/core/api/internal/svc"
	"oa.98ent.com/p9/core/common/i18n"
)

// 注册菜单、API目录、多语言数据
func Register(serverCtx *svc.ServiceContext) error {
	// 根据模式获取i18n代码
	i18nCode := i18n.CodeByPartnerMode(serverCtx.Config.PartnerMode)
	req := catalogReq(i18nCode)
	// Core Rpc可能还未就绪，失败时短暂重试5次，每次间隔200毫秒
	err := retry.Retry(
		func() error {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()

			_, err := serverCtx.Core.RegisterCatalog(ctx, req)
			return err
		},
		retry.RetryTimes(5),
		retry.RetryWithLinearBackoff(200*time.Millisecond),
	)
	if err != nil {
		return fmt.Errorf("register core-api catalog: %w", err)
	}

	return nil
}
