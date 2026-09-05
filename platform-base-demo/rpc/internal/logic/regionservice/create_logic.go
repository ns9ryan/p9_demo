package regionservicelogic

import (
	"context"
	"regexp"
	"strings"

	entsql "entgo.io/ent/dialect/sql"
	"oa.98ent.com/p9/platform-base/pkg/i18nkey"
	"oa.98ent.com/p9/platform-base/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-base/rpc/ent"
	entlanguage "oa.98ent.com/p9/platform-base/rpc/ent/language"
	entregion "oa.98ent.com/p9/platform-base/rpc/ent/region"
	"oa.98ent.com/p9/platform-base/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-base/rpc/internal/svc"
	"oa.98ent.com/p9/platform-base/rpc/pb/base/region"

	"github.com/zeromicro/go-zero/core/logx"
)

var (
	regionCodeRegexp  = regexp.MustCompile(`^[A-Z]{2}$`)
	callingCodeRegexp = regexp.MustCompile(`^[0-9]{1,3}$`)
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

// Create 创建国家地区
func (l *CreateLogic) Create(in *region.CreateRegionRequest) (*region.CreateRegionResponse, error) {
	// 校验国家地区编码
	if !regionCodeRegexp.MatchString(in.Code) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验国际电话区号
	if !callingCodeRegexp.MatchString(in.CallingCode) {
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
	last, err := l.svcCtx.DB.Region.
		Query().
		Order(
			entregion.BySortNo(entsql.OrderDesc()), // 按排序值降序
			entregion.ByID(entsql.OrderDesc()),     // 排序值相同时按ID降序
		).
		First(l.ctx)
	if err != nil && !ent.IsNotFound(err) {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 已有国家地区时追加到最后
	if last != nil {
		sortNo = last.SortNo + 1
	}

	// 创建国家地区
	result, err := l.svcCtx.DB.Region.
		Create().
		SetCode(in.Code).               // 国家或地区编码
		SetCallingCode(in.CallingCode). // 国际电话区号
		SetNameI18n(nameI18N).          // 多语言名称
		SetNillableStatus(in.Status).   // 状态: 1启用, 2停用
		SetSortNo(sortNo).              // 排序值
		Save(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	return &region.CreateRegionResponse{
		Id: result.ID,
	}, nil
}
