package currencyservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	entcurrency "oa.98ent.com/p9/platform-base/rpc/ent/currency"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/currency"
)

type GetByCodeLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetByCodeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetByCodeLogic {
	return &GetByCodeLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// GetByCode 按编码获取货币
func (l *GetByCodeLogic) GetByCode(in *currency.GetCurrencyByCodeRequest) (*currency.GetCurrencyByCodeResponse, error) {
	// 整理货币编码
	code := strings.ToUpper(strings.TrimSpace(in.Code))
	if code == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 按编码获取货币
	result, err := l.svcCtx.DB.Currency.
		Query().
		Where(entcurrency.CodeEQ(code)).
		Only(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回货币信息
	return &currency.GetCurrencyByCodeResponse{
		Currency: toCurrencyInfo(result), // 货币信息
	}, nil
}
