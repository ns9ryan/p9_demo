package gamecurrencyservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateGameCurrencyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateGameCurrencyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateGameCurrencyLogic {
	return &UpdateGameCurrencyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新游戏货币
func (l *UpdateGameCurrencyLogic) UpdateGameCurrency(in *platformgame.UpdateGameCurrencyRequest) (*platformgame.UpdateGameCurrencyResp, error) {
	l.Infof("[RPC UpdateGameCurrency] received request: id=%d", in.Id)

	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		return &platformgame.UpdateGameCurrencyResp{
			Code:    constant.CodeInternalError,
			Message: "DAO Manager not available",
		}, nil
	}

	// 准备更新数据
	updates := make(map[string]interface{})
	if in.Status > 0 {
		updates["status"] = int64(in.Status)
	}
	updates["updated_at"] = time.Now()

	// 执行更新
	_, err := l.svcCtx.DAOManager.GameCurrency.UpdateGameCurrency(l.ctx, in.Id, updates)
	if err != nil {
		return &platformgame.UpdateGameCurrencyResp{
			Code:    constant.CodeInternalError,
			Message: err.Error(),
		}, nil
	}

	return &platformgame.UpdateGameCurrencyResp{
		Code:    constant.CodeSuccess,
		Message: "ok",
	}, nil
}
