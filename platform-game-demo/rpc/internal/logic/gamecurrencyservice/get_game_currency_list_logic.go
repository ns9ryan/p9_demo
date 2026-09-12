package gamecurrencyservicelogic

import (
	"context"
	"fmt"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/logic"
	"oa.98ent.com/p9/platform-game/rpc/internal/svc"
	"oa.98ent.com/p9/platform-game/rpc/internal/utils"
	platformgame "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
)

type GetGameCurrencyListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetGameCurrencyListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetGameCurrencyListLogic {
	return &GetGameCurrencyListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取游戏货币列表
func (l *GetGameCurrencyListLogic) GetGameCurrencyList(in *platformgame.GetGameCurrencyListRequest) (*platformgame.GetGameCurrencyListResp, error) {
	l.Infof("[RPC GetGameCurrencyList] received req: page=%d, page_size=%d", in.Page, in.PageSize)
	if l.svcCtx == nil || l.svcCtx.DB == nil {
		l.Errorf("[RPC GetGameCurrencyList] Database not available")
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: "Database not available",
		}, nil
	}

	page, pageSize := utils.HandlePage(int64(in.Page), int64(in.PageSize))

	query := l.svcCtx.DB.WithContext(l.ctx)

	// 处理软删除条件
	if in.GetIsDeleted() == 0 {
		query = query.Where("deleted_at IS NULL")
	} else if in.GetIsDeleted() == 1 {
		query = query.Where("deleted_at IS NOT NULL")
	}

	if in.GetStatus() > 0 {
		query = query.Where("status = ?", in.GetStatus())
	}
	if in.GetGameId() > 0 {
		query = query.Where("game_id = ?", in.GetGameId())
	}
	if in.GetCurrencyId() > 0 {
		query = query.Where("currency_id = ?", in.GetCurrencyId())
	}

	sortBy := strings.TrimSpace(in.SortBy)
	sortOrder := strings.TrimSpace(in.SortOrder)
	if sortBy == "" {
		sortBy = "id"
	}
	if sortOrder == "" {
		sortOrder = "asc"
	}
	sortBy = strings.ToLower(sortBy)
	sortOrder = strings.ToLower(sortOrder)

	switch sortBy {
	case "id", "created_at":
	default:
		sortBy = "id"
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}

	var total int64
	if err := query.Model(&ent.GameCurrency{}).Count(&total).Error; err != nil {
		l.Errorf("[RPC GetGameCurrencyList] count failed: %v", err)
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("count failed: %v", err),
		}, nil
	}

	var currencies []*ent.GameCurrency
	if err := query.
		Offset(int((page - 1) * pageSize)).
		Limit(int(pageSize)).
		Order(fmt.Sprintf("%s %s", sortBy, sortOrder)).
		Find(&currencies).Error; err != nil {
		l.Errorf("[RPC GetGameCurrencyList] query failed: %v", err)
		return &platformgame.GetGameCurrencyListResp{
			Code:    constant.CodeInternalError,
			Message: fmt.Sprintf("query failed: %v", err),
		}, nil
	}
	sysCurrencyMap := GetSysCurrencyMap(l.ctx, l.svcCtx)
	currenciesProto := make([]*platformgame.GameCurrencyInfo, 0, len(currencies))
	for _, currency := range currencies {
		if currency == nil {
			continue
		}
		gameRecord, err := GetGameRecord(l.ctx, l.svcCtx, currency.GameId)
		if err != nil {
			l.Errorf("[RPC GetGameCurrencyList] query game failed: %v", err)
		}
		currenciesProto = append(currenciesProto, logic.CurrencyModelToProto(currency, gameRecord, sysCurrencyMap))
	}
	l.Infof("[RPC GetGameCurrencyList] success: total=%d", total)
	return &platformgame.GetGameCurrencyListResp{
		Code:     constant.CodeSuccess,
		Message:  "success",
		Data:     currenciesProto,
		Total:    total,
		Page:     int32(page),
		PageSize: int32(pageSize),
	}, nil
}
