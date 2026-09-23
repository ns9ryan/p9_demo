// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator

import (
	"context"
	"strings"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-base/rpc/pb/platformbaserpc/currencypb"
	"oa.98ent.com/p9/platform-base/rpc/pb/platformbaserpc/timezonepb"
	"oa.98ent.com/p9/platform-operator/api/internal/i18nkey"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOperatorLogic {
	return &CreateOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// CreateOperator 创建分站
func (l *CreateOperatorLogic) CreateOperator(req *types.CreateOperatorRequest) (resp *types.CreateOperatorResponse, err error) {
	// 整理时区编码
	timezoneCode := strings.TrimSpace(req.TimezoneCode)

	// 整理结算币种编码
	settlementCurrencyCode := strings.ToUpper(strings.TrimSpace(req.SettlementCurrencyCode))

	// 获取并校验时区
	timezoneResult, err := l.svcCtx.TimezoneRpc.GetByCode(
		l.ctx,
		&timezonepb.GetTimezoneByCodeRequest{
			Code: timezoneCode, // IANA时区编码
		},
	)
	if err != nil {
		return nil, err
	}

	// 停用的时区不可用于创建分站
	if timezoneResult.Timezone.Status != 1 {
		return nil, xerr.BadRequest(i18nkey.TimezoneUnavailable)
	}

	// 获取并校验结算币种
	currencyResult, err := l.svcCtx.CurrencyRpc.GetByCode(
		l.ctx,
		&currencypb.GetCurrencyByCodeRequest{
			Code: settlementCurrencyCode, // 货币编码
		},
	)
	if err != nil {
		return nil, err
	}

	// 停用的结算币种不可用于创建分站
	if currencyResult.Currency.Status != 1 {
		return nil, xerr.BadRequest(i18nkey.SettlementCurrencyUnavailable)
	}

	// 创建分站
	result, err := l.svcCtx.OperatorRpc.Create(
		l.ctx,
		&operatorpb.CreateOperatorRequest{
			Name:                   req.Name,                     // 分站名称
			TimezoneCode:           timezoneResult.Timezone.Code, // 时区编码
			SettlementCurrencyCode: currencyResult.Currency.Code, // 结算币种编码
			Status:                 req.Status,                   // 分站状态: 1正常, 2暂停, 3关闭
			Remark:                 req.Remark,                   // 内部备注
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回创建结果
	return &types.CreateOperatorResponse{
		Id:   result.Id,   // 分站ID
		Code: result.Code, // 分站业务编码
	}, nil
}
