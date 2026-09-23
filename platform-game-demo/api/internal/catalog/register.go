package catalog

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/common/i18n"
	"oa.98ent.com/p9/platform-game/api/internal/svc"
)

// 注册菜单、API目录、多语言数据
func Register(serverCtx *svc.ServiceContext) error {
	req := AdminReq(i18n.CodePlatform)
	var last error
	for range 5 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, last = serverCtx.Core.RegisterCatalog(ctx, req)
		cancel()
		if last == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("register platform-game-api catalog: %w", last)
}
