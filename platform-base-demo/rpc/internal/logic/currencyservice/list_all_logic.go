package currencyservicelogic

import (
	"context"

	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	entcurrency "oa.98ent.com/p9/platform-base/rpc/ent/currency"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/currency"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListAllLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListAllLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListAllLogic {
	return &ListAllLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// ListAll 获取全部货币
func (l *ListAllLogic) ListAll(in *currency.ListAllCurrenciesRequest) (*currency.ListAllCurrenciesResponse, error) {
	// 校验状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 2) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 创建货币查询
	query := l.svcCtx.DB.Currency.Query()

	// 按状态筛选
	if in.Status != nil {
		query = query.Where(entcurrency.StatusEQ(*in.Status))
	}

	// 获取全部货币
	results, err := query.
		Order(
			entcurrency.BySortNo(), // 按排序值升序
			entcurrency.ByID(),     // 排序值相同时按ID升序
		).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换货币列表
	list := make([]*currency.CurrencyInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toCurrencyInfo(result))
	}

	return &currency.ListAllCurrenciesResponse{
		List: list,
	}, nil
}
