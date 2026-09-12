package operatorservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/operator"
)

type CompleteLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCompleteLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteLogic {
	return &CompleteLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Complete 完成分站创建
func (l *CompleteLogic) Complete(in *operator.CompleteOperatorRequest) (*operator.CompleteOperatorResponse, error) {
	// 分站ID必须大于0
	if in.Id <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 获取当前分站
	current, err := l.svcCtx.DB.Operator.Get(l.ctx, in.Id)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 已完成时直接返回
	if current.CreationStatus == 2 {
		return &operator.CompleteOperatorResponse{}, nil
	}

	// 完成分站创建
	err = current.
		Update().
		SetCreationStatus(2). // 创建状态: 2已完成
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回完成结果
	return &operator.CompleteOperatorResponse{}, nil
}
