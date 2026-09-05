package languageservicelogic

import (
	"context"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-base/rpc/ent"
	entlanguage "oa.98ent.com/p9/platform-base/rpc/ent/language"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/language"

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

// Create 创建语言
func (l *CreateLogic) Create(in *language.CreateLanguageRequest) (*language.CreateLanguageResponse, error) {
	// 多语言名称不能为空
	if len(in.NameI18N) == 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 必须包含当前语言自己的名称
	if _, ok := in.NameI18N[in.Code]; !ok {
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

		// 当前正在创建的语言不需要查询数据库
		if code != in.Code {
			codes = append(codes, code)
		}
	}

	// 校验其他语言编码是否已经存在
	if len(codes) > 0 {
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
	}

	// 默认从第一个排序位置开始
	sortNo := int64(1)

	// 获取当前最后一个排序值
	last, err := l.svcCtx.DB.Language.
		Query().
		Order(
			entlanguage.BySortNo(entsql.OrderDesc()), // 按排序值降序
			entlanguage.ByID(entsql.OrderDesc()),     // 排序值相同时按ID降序
		).
		First(l.ctx)
	if err != nil && !ent.IsNotFound(err) {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 已有语言时追加到最后
	if last != nil {
		sortNo = last.SortNo + 1
	}

	// 创建语言
	result, err := l.svcCtx.DB.Language.
		Create().
		SetCode(in.Code).             // 语言编码
		SetNameI18n(nameI18N).        // 多语言名称
		SetNillableStatus(in.Status). // 状态: 1启用, 2停用
		SetSortNo(sortNo).            // 排序值
		Save(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	return &language.CreateLanguageResponse{
		Id: result.ID,
	}, nil
}
