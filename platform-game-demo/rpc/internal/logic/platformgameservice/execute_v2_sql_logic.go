package platformgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	sync "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type ExecuteV2SqlLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewExecuteV2SqlLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ExecuteV2SqlLogic {
	return &ExecuteV2SqlLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 执行 v2sql 下的 SQL 初始化脚本（供运维/初始化使用）
func (l *ExecuteV2SqlLogic) ExecuteV2Sql(in *sync.ExecuteV2SqlRequest) (*sync.ExecuteV2SqlResp, error) {
	// todo: add your logic here and delete this line

	return &sync.ExecuteV2SqlResp{}, nil
}
