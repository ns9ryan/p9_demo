package operatorgamecategoryservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOperatorGameCategoryListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOperatorGameCategoryListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOperatorGameCategoryListLogic {
	return &GetOperatorGameCategoryListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取分站游戏分类列表
func (l *GetOperatorGameCategoryListLogic) GetOperatorGameCategoryList(in *platform_game.GetOperatorGameCategoryListRequest) (*platform_game.GetOperatorGameCategoryListResp, error) {
	l.Infof("[RPC GetOperatorGameCategoryList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DAOManager == nil {
		l.Errorf("[RPC GetOperatorGameCategoryList] DAO Manager not available")
		return &platform_game.GetOperatorGameCategoryListResp{}, nil
	}
	if in.GetOpCode() == "" {
		l.Errorf("[RPC GetOperatorGameCategoryList] op_code is required")
		return &platform_game.GetOperatorGameCategoryListResp{}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))
	offset := (page - 1) * pageSize

	// check_status状态为1时为已分配，为2时为未分配。check_status=1时直接查询operator_game_category表中已分配的记录
	if in.GetCheckStatus() == 1 {
		operatorGameCategoryRecords, _ := l.svcCtx.DAOManager.OperatorGameCategory.FindAll(
			l.ctx,
			in.GetOpCode(),
			in.GetCategoryCode(),
			offset,
			pageSize,
		)
		return &platform_game.GetOperatorGameCategoryListResp{
			Items:    logic.OperatorGameCategoryModelToProtoListWithoutMap(operatorGameCategoryRecords),
			Total:    int64(len(operatorGameCategoryRecords)),
			Page:     int32(page),
			PageSize: int32(pageSize),
		}, nil
	}
	// 查询总网游戏分类列表
	records, total, err := l.svcCtx.DAOManager.GameCategory.GetGameCategoryList(
		l.ctx,
		0,
		1,
		in.GetCategoryCode(),
		offset,
		pageSize,
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameCategoryList] query failed: %v", err)
		return &platform_game.GetOperatorGameCategoryListResp{}, nil
	}
	// records对应OperatorGameCategory表中是否有数据,存在则设置check_status字段为1，不存在为2
	recordsInOperatorGameCategory, err := l.svcCtx.DAOManager.OperatorGameCategory.FindAllByOpCodeAndCategoryCodes(
		l.ctx,
		in.GetOpCode(),
		func() []string {
			codes := make([]string, 0, len(records))
			for _, record := range records {
				codes = append(codes, record.SourceCategoryCode)
			}
			return codes
		}(),
	)
	if err != nil {
		l.Errorf("[RPC GetOperatorGameCategoryList] query OperatorGameCategory failed: %v", err)
		return &platform_game.GetOperatorGameCategoryListResp{}, nil
	}
	mapCategoryCodeToRecord := make(map[string]*ent.OperatorGameCategory)
	for _, record := range recordsInOperatorGameCategory {
		mapCategoryCodeToRecord[record.CategoryCode] = record
	}

	l.Infof("[RPC GetOperatorGameCategoryList] success: total=%d", total)
	return &platform_game.GetOperatorGameCategoryListResp{
		Items:    logic.OperatorGameCategoryModelToProtoList(in.GetOpCode(), mapCategoryCodeToRecord, records),
		Total:    int64(total),
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
