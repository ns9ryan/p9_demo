package regionallocationservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorregionallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/regionallocation"
)

type SaveLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSaveLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLogic {
	return &SaveLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// Save 保存经营地区分配
func (l *SaveLogic) Save(in *regionallocation.SaveRegionAllocationsRequest) (*regionallocation.SaveRegionAllocationsResponse, error) {
	// 分站ID必须大于0
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 整理经营地区编码并去重
	regionCodes := make([]string, 0, len(in.RegionCodes))
	regionCodeSet := make(map[string]struct{}, len(in.RegionCodes))

	for _, code := range in.RegionCodes {
		regionCode := strings.ToUpper(strings.TrimSpace(code))
		if regionCode == "" {
			return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
		}

		if _, exists := regionCodeSet[regionCode]; exists {
			continue
		}

		regionCodeSet[regionCode] = struct{}{}
		regionCodes = append(regionCodes, regionCode)
	}

	// 确认分站存在
	_, err := l.svcCtx.DB.Operator.Get(l.ctx, in.OperatorId)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
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

	// 获取当前经营地区编码
	currentCodes, err := tx.OperatorRegionAllocation.
		Query().
		Where(operatorregionallocation.OperatorIDEQ(in.OperatorId)).
		Select(operatorregionallocation.FieldRegionCode).
		Strings(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 整理当前经营地区编码
	currentCodeSet := make(map[string]struct{}, len(currentCodes))
	for _, regionCode := range currentCodes {
		currentCodeSet[regionCode] = struct{}{}
	}

	// 计算需要删除的经营地区编码
	deleteCodes := make([]string, 0)
	for _, regionCode := range currentCodes {
		if _, exists := regionCodeSet[regionCode]; !exists {
			deleteCodes = append(deleteCodes, regionCode)
		}
	}

	// 计算需要新增的经营地区编码
	createCodes := make([]string, 0)
	for _, regionCode := range regionCodes {
		if _, exists := currentCodeSet[regionCode]; !exists {
			createCodes = append(createCodes, regionCode)
		}
	}

	// 批量删除已经取消的经营地区分配
	if len(deleteCodes) > 0 {
		_, err = tx.OperatorRegionAllocation.
			Delete().
			Where(
				operatorregionallocation.OperatorIDEQ(in.OperatorId),
				operatorregionallocation.RegionCodeIn(deleteCodes...),
			).
			Exec(l.ctx)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}
	}

	// 批量创建新增的经营地区分配
	if len(createCodes) > 0 {
		builders := make([]*ent.OperatorRegionAllocationCreate, 0, len(createCodes))

		for _, regionCode := range createCodes {
			builders = append(
				builders,
				tx.OperatorRegionAllocation.
					Create().
					SetOperatorID(in.OperatorId). // 分站ID
					SetRegionCode(regionCode),    // 经营地区编码
			)
		}

		_, err = tx.OperatorRegionAllocation.
			CreateBulk(builders...).
			Save(l.ctx)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}
	}

	// 提交数据库事务
	if err = tx.Commit(); err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 返回保存结果
	return &regionallocation.SaveRegionAllocationsResponse{}, nil
}
