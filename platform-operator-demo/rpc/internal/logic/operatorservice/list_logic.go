package operatorservicelogic

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	entoperator "oa.98ent.com/p9/platform-operator/rpc/ent/operator"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/operator"
)

type ListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListLogic {
	return &ListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// List 获取分站管理列表
func (l *ListLogic) List(in *operator.ListOperatorsRequest) (*operator.ListOperatorsResponse, error) {
	// 校验分页参数
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验创建状态
	if in.CreationStatus != nil && (*in.CreationStatus < 1 || *in.CreationStatus > 2) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验发布状态
	if in.PublishStatus != nil && (*in.PublishStatus < 1 || *in.PublishStatus > 4) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 校验分站状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 3) {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 创建分站查询
	query := l.svcCtx.DB.Operator.Query()

	// 按关键字筛选
	if in.Keyword != nil {
		keyword := strings.TrimSpace(*in.Keyword)
		if keyword != "" {
			query = query.Where(
				entoperator.Or(
					entoperator.CodeContainsFold(keyword), // 匹配分站编码
					entoperator.NameContainsFold(keyword), // 匹配分站名称
				),
			)
		}
	}

	// 按创建状态筛选
	if in.CreationStatus != nil {
		query = query.Where(entoperator.CreationStatusEQ(*in.CreationStatus))
	}

	// 按发布状态筛选
	if in.PublishStatus != nil {
		query = query.Where(entoperator.PublishStatusEQ(*in.PublishStatus))
	}

	// 按分站状态筛选
	if in.Status != nil {
		query = query.Where(entoperator.StatusEQ(*in.Status))
	}

	// 获取符合条件的数据总数
	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 计算分页偏移量
	offset := (in.Page - 1) * in.PageSize

	// 获取当前页分站数据
	results, err := query.
		Order(
			entoperator.ByCreatedAt(sql.OrderDesc()), // 按创建时间倒序
			entoperator.ByID(sql.OrderDesc()),        // 创建时间相同时按ID倒序
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 转换分站列表
	list := make([]*operator.OperatorInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toOperatorInfo(result))
	}

	// 返回分站列表
	return &operator.ListOperatorsResponse{
		Total: int64(total), // 数据总数
		List:  list,         // 分站列表
	}, nil
}
