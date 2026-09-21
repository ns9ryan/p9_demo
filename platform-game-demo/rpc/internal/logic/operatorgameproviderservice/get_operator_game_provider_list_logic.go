package operatorgameproviderservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameProviderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameProviderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameProviderListLogic {
	return &GetOperatorGameProviderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏提供商列表
func (l *GetOperatorGameProviderListLogic) GetOperatorGameProviderList(in *platform_game.GetOperatorGameProviderListRequest) (*platform_game.GetOperatorGameProviderListResp, error) {
	l.Infof("[RPC GetOperatorGameProviderList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameProviderList] DAO Manager not available")
		return &platform_game.GetOperatorGameProviderListResp{}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	if in.GetCheckStatus() == 1 {
		operatorGameProviderRecords, _ := l.svcCtx.DAOManager.OperatorGameProvider.FindAll(
			l.ctx,
			in.GetOpCode(),
			in.GetProviderCode(),
			offset,
			pageSize,
		)
		return &platform_game.GetOperatorGameProviderListResp{
			Items:    logic.OperatorGameProviderModelToProtoListWithoutMap(operatorGameProviderRecords),
			Total:    int64(len(operatorGameProviderRecords)),
			Page:     int32(page),
			PageSize: int32(pageSize),
		}, nil
	}

	// 查询分站游戏提供商列表
	records, total, err := l.svcCtx.DAOManager.GameProvider.GetGameProviderList(
		l.ctx,
		0,
		1,
		in.GetProviderCode(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameProviderList] query failed: %v", err)
		return &platform_game.GetOperatorGameProviderListResp{}, nil
	}
	// records对应OperatorGameProvider表中是否有数据,存在则设置check_status字段为1，不存在为2
	recordsInOperatorGameProvider, err := l.svcCtx.DAOManager.OperatorGameProvider.FindAllByOpCodeAndProviderCodes(
		l.ctx,
		in.GetOpCode(),
		func() []string {
			codes := make([]string, 0, len(records))
			for _, record := range records {
				codes = append(codes, record.SourceProviderCode)
			}
			return codes
		}(),
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameProviderList] query OperatorGameProvider failed: %v", err)
		return &platform_game.GetOperatorGameProviderListResp{}, nil
	}
	mapProviderCodeToRecord := make(map[string]*ent.OperatorGameProvider)
	for _, record := range recordsInOperatorGameProvider {
		mapProviderCodeToRecord[record.ProviderCode] = record
	}

	l.Infof("[RPC GetOperatorGameProviderList] success: total=%d", total)
	return &platform_game.GetOperatorGameProviderListResp{
		Items:    logic.OperatorGameProviderModelToProtoList(in.GetOpCode(), mapProviderCodeToRecord, records),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
