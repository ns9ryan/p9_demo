package catalog

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/core/rpc/coreclient"
)

func Register(cli coreclient.Core) error {
	req := AdminReq()
	var last error
	for range 20 {
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		_, last = cli.RegisterCatalog(ctx, req)
		cancel()
		if last == nil {
			return nil
		}
		time.Sleep(200 * time.Millisecond)
	}
	return fmt.Errorf("register core-api catalog: %w", last)
}
