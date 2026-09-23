package operatoradminservicelogic

import (
	"context"
	"strings"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatoradmin"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/adminpb"

	"entgo.io/ent/dialect/sql"
	"github.com/zeromicro/go-zero/core/logx"
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

// List 获取分站管理员列表
func (l *ListLogic) List(in *adminpb.ListAdminsRequest) (*adminpb.ListAdminsResponse, error) {
	// 校验分页参数
	if in.Page < 1 || in.PageSize < 1 || in.PageSize > 100 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 校验分站ID
	if in.OperatorId != nil && *in.OperatorId <= 0 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 校验状态
	if in.Status != nil && (*in.Status < 1 || *in.Status > 2) {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	query := l.svcCtx.DB.OperatorAdmin.Query()

	// 按分站筛选
	if in.OperatorId != nil {
		query = query.Where(operatoradmin.OperatorIDEQ(*in.OperatorId))
	}

	// 按关键字筛选, 匹配账号或显示名称
	if in.Keyword != nil {
		keyword := strings.TrimSpace(*in.Keyword)
		if keyword != "" {
			query = query.Where(
				operatoradmin.Or(
					operatoradmin.UsernameContainsFold(keyword),
					operatoradmin.DisplayNameContainsFold(keyword),
				),
			)
		}
	}

	// 按状态筛选
	if in.Status != nil {
		query = query.Where(operatoradmin.StatusEQ(*in.Status))
	}

	total, err := query.Clone().Count(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	offset := (in.Page - 1) * in.PageSize

	results, err := query.
		WithOperator().
		Order(
			operatoradmin.ByCreatedAt(sql.OrderDesc()),
			operatoradmin.ByID(sql.OrderDesc()),
		).
		Offset(int(offset)).
		Limit(int(in.PageSize)).
		All(l.ctx)
	if err != nil {
		return nil, enterror.Handle(l.Logger, err)
	}

	list := make([]*adminpb.AdminInfo, 0, len(results))
	for _, result := range results {
		list = append(list, toAdminInfo(result))
	}

	return &adminpb.ListAdminsResponse{
		Total: int64(total),
		List:  list,
	}, nil
}
