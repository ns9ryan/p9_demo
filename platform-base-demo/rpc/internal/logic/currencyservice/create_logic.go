package currencyservicelogic

import (
	"context"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-base/rpc/ent"
	entcurrency "oa.98ent.com/p9/platform-base/rpc/ent/currency"
	entlanguage "oa.98ent.com/p9/platform-base/rpc/ent/language"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/currency"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateLogic {
	return &CreateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Create 创建货币
func (l *CreateLogic) Create(in *currency.CreateCurrencyRequest) (*currency.CreateCurrencyResponse, error) {
	// 校验货币类型
	if in.CurrencyType < 1 || in.CurrencyType > 2 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 金额换算倍率必须大于0
	if in.AmountFactor <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 货币符号不能为空
	symbol := strings.TrimSpace(in.Symbol)
	if symbol == "" {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 多语言名称不能为空
	if len(in.NameI18N) == 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验并整理多语言名称
	nameI18N := make(map[string]string, len(in.NameI18N))
	codes := make([]string, 0, len(in.NameI18N))

	for code, name := range in.NameI18N {
		name = strings.TrimSpace(name)
		if code == "" || name == "" {
			return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
		}

		nameI18N[code] = name
		codes = append(codes, code)
	}

	// 校验多语言名称中的语言编码是否已经存在
	count, err := l.svcCtx.DB.Language.
		Query().
		Where(entlanguage.CodeIn(codes...)).
		Count(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	if count != len(codes) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 默认从第一个排序位置开始
	sortNo := int64(1)

	// 获取当前最后一个排序值
	last, err := l.svcCtx.DB.Currency.
		Query().
		Order(
			entcurrency.BySortNo(entsql.OrderDesc()), // 按排序值降序
			entcurrency.ByID(entsql.OrderDesc()),     // 排序值相同时按ID降序
		).
		First(l.ctx)
	if err != nil && !ent.IsNotFound(err) {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 已有货币时追加到最后
	if last != nil {
		sortNo = last.SortNo + 1
	}

	// 创建货币
	result, err := l.svcCtx.DB.Currency.
		Create().
		SetCode(in.Code).                 // 货币编码
		SetNameI18n(nameI18N).            // 多语言名称
		SetCurrencyType(in.CurrencyType). // 货币类型: 1法定货币, 2虚拟货币
		SetSymbol(symbol).                // 货币符号
		SetAmountFactor(in.AmountFactor). // 金额换算倍率
		SetNillableStatus(in.Status).     // 状态: 1启用, 2停用
		SetSortNo(sortNo).                // 排序值
		Save(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	return &currency.CreateCurrencyResponse{
		Id: result.ID,
	}, nil
}
