package operatorgamechannelservicelogic

import (
	"context"
	"time"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type BatchCreateOperatorGameChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchCreateOperatorGameChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchCreateOperatorGameChannelLogic {
	return &BatchCreateOperatorGameChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 批量创建分站游戏渠道
func (l *BatchCreateOperatorGameChannelLogic) BatchCreateOperatorGameChannel(in *platform_game.BatchCreateOperatorGameChannelRequest) (*platform_game.BatchCreateOperatorGameChannelResp, error) {
	l.Infof("[RPC BatchCreateOperatorGameChannel] received req: %d items", len(in.Items))
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC BatchCreateOperatorGameChannel] DAO Manager not available")
		return &platform_game.BatchCreateOperatorGameChannelResp{}, nil
	}

	if len(in.Items) == 0 {
		return &platform_game.BatchCreateOperatorGameChannelResp{
			Total:   0,
			Success: 0,
			Failed:  0,
		}, nil
	}

	// 构建创建对象列表
	createList := make([]*ent.OperatorGameChannelCreate, 0, len(in.Items))
	for _, item := range in.Items {
		// 判断对应code的分站是否存在
		exists, err := l.svcCtx.DAOManager.Operator.ExistByCode(l.ctx, item.OpCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameChannel] 检查分站 %s 是否存在失败: %v", item.OpCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGameChannel] 分站 %s 不存在", item.OpCode)
			continue
		}

		// 判断对应code的游戏是否存在
		exists, err = l.svcCtx.DAOManager.GameChannel.ExistByCode(l.ctx, item.ChannelCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameChannel] 检查游戏渠道 %s 是否存在失败: %v", item.ChannelCode, err)
			continue
		}
		if !exists {
			l.Errorf("[RPC BatchCreateOperatorGameChannel] 游戏渠道 %s 不存在", item.ChannelCode)
			continue
		}
		// 判断是否已经存在相同的分站游戏渠道记录
		exists, err = l.svcCtx.DAOManager.OperatorGameChannel.ExistByOpCodeAndChannelCode(l.ctx, item.OpCode, item.ChannelCode)
		if err != nil {
			l.Errorf("[RPC BatchCreateOperatorGameChannel] 检查分站游戏渠道 %s-%s 是否存在失败: %v", item.OpCode, item.ChannelCode, err)
			continue
		}
		if exists {
			l.Errorf("[RPC BatchCreateOperatorGameChannel] 分站游戏渠道 %s-%s 已存在", item.OpCode, item.ChannelCode)
			continue
		}
		createList = append(createList, l.svcCtx.DB.OperatorGameChannel.Create().
			SetOpCode(item.OpCode).
			SetChannelCode(item.ChannelCode).
			SetStatus(int16(item.Status)).
			SetCreatedAt(time.Now()).
			SetUpdatedAt(time.Now()))
	}
	if len(createList) == 0 {
		l.Errorf("[RPC BatchCreateOperatorGameChannel] no valid items to create")
		return &platform_game.BatchCreateOperatorGameChannelResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)) - int64(len(createList)),
		}, nil
	}
	records, err := l.svcCtx.DAOManager.OperatorGameChannel.BatchCreateOperatorGameChannel(l.ctx, createList)
	if err != nil {
		l.Errorf("[RPC BatchCreateOperatorGameChannel] batch create failed: %v", err)
		return &platform_game.BatchCreateOperatorGameChannelResp{
			Total:   int64(len(in.Items)),
			Success: 0,
			Failed:  int64(len(in.Items)),
		}, nil
	}

	l.Infof("[RPC BatchCreateOperatorGameChannel] success: created %d records", len(records))
	return &platform_game.BatchCreateOperatorGameChannelResp{
		Total:   int64(len(in.Items)),
		Success: int64(len(records)),
		Failed:  int64(len(in.Items)) - int64(len(records)),
	}, nil
}
