package operatorservicelogic

import (
	"context"

	"oa.98ent.com/p9/common/xerr"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatoradmin"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatoragentlineallocation"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatordomain"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorlanguageallocation"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorprofile"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorregionallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/operatorpb"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteLogic {
	return &DeleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Delete 删除分站
func (l *DeleteLogic) Delete(in *operatorpb.DeleteOperatorRequest) (*operatorpb.DeleteOperatorResponse, error) {
	// 分站ID必须大于0
	if in.Id <= 0 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ValidationError))
	}

	// 获取当前分站
	current, err := l.svcCtx.DB.Operator.Get(l.ctx, in.Id)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 仅未发布或发布失败的分站允许删除
	if current.PublishStatus != 1 && current.PublishStatus != 4 {
		return nil, xerr.RpcErr(xerr.BadRequest(i18nkey.ConstraintError))
	}

	// 开启数据库事务
	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	// 删除语言分配
	_, err = tx.OperatorLanguageAllocation.
		Delete().
		Where(operatorlanguageallocation.OperatorIDEQ(in.Id)).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 删除经营地区分配
	_, err = tx.OperatorRegionAllocation.
		Delete().
		Where(operatorregionallocation.OperatorIDEQ(in.Id)).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 删除代理子线路分配
	_, err = tx.OperatorAgentLineAllocation.
		Delete().
		Where(operatoragentlineallocation.OperatorIDEQ(in.Id)).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 删除分站域名
	_, err = tx.OperatorDomain.
		Delete().
		Where(operatordomain.OperatorIDEQ(in.Id)).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 删除分站档案
	_, err = tx.OperatorProfile.
		Delete().
		Where(operatorprofile.OperatorIDEQ(in.Id)).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 按分站清理管理员账号
	_, err = tx.OperatorAdmin.
		Delete().
		Where(operatoradmin.OperatorIDEQ(in.Id)).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 删除分站
	err = tx.Operator.
		DeleteOneID(in.Id).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 提交数据库事务
	if err = tx.Commit(); err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回删除结果
	return &operatorpb.DeleteOperatorResponse{}, nil
}
