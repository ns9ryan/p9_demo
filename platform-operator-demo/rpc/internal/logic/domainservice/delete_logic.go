package domainservicelogic

import (
	"context"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/domain"
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

// Delete 删除分站域名
func (l *DeleteLogic) Delete(in *domain.DeleteDomainRequest) (*domain.DeleteDomainResponse, error) {
	// 域名ID必须大于0
	if in.Id <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 删除分站域名
	err := l.svcCtx.DB.OperatorDomain.
		DeleteOneID(in.Id).
		Exec(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回删除结果
	return &domain.DeleteDomainResponse{}, nil
}
