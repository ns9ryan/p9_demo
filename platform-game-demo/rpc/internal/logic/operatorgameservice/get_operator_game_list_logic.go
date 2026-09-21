package operatorgameservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameListLogic {
	return &GetOperatorGameListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏列表
func (l *GetOperatorGameListLogic) GetOperatorGameList(in *platform_game.GetOperatorGameListRequest) (*platform_game.GetOperatorGameListResp, error) {
	l.Infof("[RPC GetOperatorGameList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameList] DAO Manager not available")
		return &platform_game.GetOperatorGameListResp{}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	if in.GetCheckStatus() == 1 {
		operatorGameRecords, _ := l.svcCtx.DAOManager.OperatorGame.FindAll(
			l.ctx,
			in.GetOpCode(),
			in.GetGameCode(),
			offset,
			pageSize,
		)
		return &platform_game.GetOperatorGameListResp{
			Items:    logic.OperatorGameModelToProtoListWithoutMap(operatorGameRecords),
			Total:    int64(len(operatorGameRecords)),
			Page:     int32(page),
			PageSize: int32(pageSize),
		}, nil
	}

	// 查询分站游戏列表
	records, total, err := l.svcCtx.DAOManager.Game.GetGameList(
		l.ctx,
		0,
		"",
		1,
		in.GetGameCode(),
		0,
		0,
		0,
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameList] query failed: %v", err)
		return &platform_game.GetOperatorGameListResp{}, nil
	}
	// records对应OperatorGame表中是否有数据,存在则设置check_status字段为1，不存在为2
	recordsInOperatorGame, err := l.svcCtx.DAOManager.OperatorGame.FindAllByOpCodeAndGameCodes(
		l.ctx,
		in.GetOpCode(),
		func() []string {
			codes := make([]string, 0, len(records))
			for _, record := range records {
				codes = append(codes, record.SourceGameCode)
			}
			return codes
		}(),
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameList] query OperatorGame failed: %v", err)
		return &platform_game.GetOperatorGameListResp{}, nil
	}
	mapGameCodeToRecord := make(map[string]*ent.OperatorGame)
	for _, record := range recordsInOperatorGame {
		mapGameCodeToRecord[record.GameCode] = record
	}
	l.Infof("[RPC GetOperatorGameList] success: total=%d, recordsInOperatorGame=%d", total, len(recordsInOperatorGame))
	return &platform_game.GetOperatorGameListResp{
		Items:    logic.OperatorGameModelToProtoList(in.GetOpCode(), mapGameCodeToRecord, records),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
