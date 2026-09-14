// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package operator

import (
	"context"
	"strings"

	"oa.98ent.com/p9/platform-base/rpc/pb/platformbaserpc/currencypb"
	"oa.98ent.com/p9/platform-base/rpc/pb/platformbaserpc/timezonepb"
	"oa.98ent.com/p9/platform-operator/api/internal/i18nkey"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type UpdateOperatorLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateOperatorLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateOperatorLogic {
	return &UpdateOperatorLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UpdateOperator 修改分站
func (l *UpdateOperatorLogic) UpdateOperator(req *types.UpdateOperatorRequest) (resp *types.UpdateOperatorResponse, err error) {
	// 校验时区
	var timezoneCode *string
	if req.TimezoneCode != nil {
		code := strings.TrimSpace(*req.TimezoneCode)

		timezoneResult, err := l.svcCtx.TimezoneRpc.GetByCode(
			l.ctx,
			&timezonepb.GetTimezoneByCodeRequest{
				Code: code, // IANA时区编码
			},
		)
		if err != nil {
			// 时区不存在时统一按不可用处理
			if status.Code(err) == codes.NotFound {
				return nil, grpcerror.InvalidArgument(i18nkey.TimezoneUnavailable)
			}

			return nil, err
		}

		// 停用的时区不可用于分站
		if timezoneResult.Timezone.Status != 1 {
			return nil, grpcerror.InvalidArgument(i18nkey.TimezoneUnavailable)
		}

		// 使用查询返回的标准时区编码
		timezoneCode = &timezoneResult.Timezone.Code
	}

	// 校验结算币种
	var settlementCurrencyCode *string
	if req.SettlementCurrencyCode != nil {
		code := strings.ToUpper(strings.TrimSpace(*req.SettlementCurrencyCode))

		currencyResult, err := l.svcCtx.CurrencyRpc.GetByCode(
			l.ctx,
			&currencypb.GetCurrencyByCodeRequest{
				Code: code, // 货币编码
			},
		)
		if err != nil {
			// 结算币种不存在时统一按不可用处理
			if status.Code(err) == codes.NotFound {
				return nil, grpcerror.InvalidArgument(i18nkey.SettlementCurrencyUnavailable)
			}

			return nil, err
		}

		// 停用的结算币种不可用于分站
		if currencyResult.Currency.Status != 1 {
			return nil, grpcerror.InvalidArgument(i18nkey.SettlementCurrencyUnavailable)
		}

		// 使用查询返回的标准结算币种编码
		settlementCurrencyCode = &currencyResult.Currency.Code
	}

	// 修改分站
	_, err = l.svcCtx.OperatorRpc.Update(
		l.ctx,
		&operatorpb.UpdateOperatorRequest{
			Id:                     req.Id,                 // 分站ID
			Name:                   req.Name,               // 分站名称
			TimezoneCode:           timezoneCode,           // 时区编码
			SettlementCurrencyCode: settlementCurrencyCode, // 结算币种编码
			Status:                 req.Status,             // 分站状态: 1正常, 2暂停, 3关闭
			Remark:                 req.Remark,             // 内部备注
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回修改结果
	return &types.UpdateOperatorResponse{}, nil
}
