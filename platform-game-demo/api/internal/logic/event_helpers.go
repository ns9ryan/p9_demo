// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package logic

import (
	"context"

	"oa.98ent.com/p9/platform-game/common/kafka"
)

// sendForceQuitEvent 发送游戏强制踢线事件
func sendForceQuitEvent(ctx context.Context, resourceID int64, scope string) error {
	return kafka.SendForceQuitEvent(ctx, resourceID, scope)
}
