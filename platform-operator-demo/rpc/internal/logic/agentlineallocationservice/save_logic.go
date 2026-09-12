package agentlineallocationservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/agentline"
	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatoragentlineallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/agentlineallocation"
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

// Save 保存代理子线路分配
func (l *SaveLogic) Save(in *agentlineallocation.SaveAgentLineAllocationsRequest) (*agentlineallocation.SaveAgentLineAllocationsResponse, error) {
	// 分站ID必须大于0
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 整理代理子线路编码并去重
	agentLineCodes := make([]string, 0, len(in.AgentLineCodes))
	agentLineCodeSet := make(map[string]struct{}, len(in.AgentLineCodes))

	for _, code := range in.AgentLineCodes {
		agentLineCode := strings.TrimSpace(code)

		// 校验代理子线路编码
		if !agentline.IsValid(agentLineCode) {
			return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
		}

		if _, exists := agentLineCodeSet[agentLineCode]; exists {
			continue
		}

		agentLineCodeSet[agentLineCode] = struct{}{}
		agentLineCodes = append(agentLineCodes, agentLineCode)
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

	// 获取当前代理子线路编码
	currentCodes, err := tx.OperatorAgentLineAllocation.
		Query().
		Where(operatoragentlineallocation.OperatorIDEQ(in.OperatorId)).
		Select(operatoragentlineallocation.FieldAgentLineCode).
		Strings(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 整理当前代理子线路编码
	currentCodeSet := make(map[string]struct{}, len(currentCodes))
	for _, agentLineCode := range currentCodes {
		currentCodeSet[agentLineCode] = struct{}{}
	}

	// 计算需要删除的代理子线路编码
	deleteCodes := make([]string, 0)
	for _, agentLineCode := range currentCodes {
		if _, exists := agentLineCodeSet[agentLineCode]; !exists {
			deleteCodes = append(deleteCodes, agentLineCode)
		}
	}

	// 计算需要新增的代理子线路编码
	createCodes := make([]string, 0)
	for _, agentLineCode := range agentLineCodes {
		if _, exists := currentCodeSet[agentLineCode]; !exists {
			createCodes = append(createCodes, agentLineCode)
		}
	}

	// 批量删除已经取消的代理子线路分配
	if len(deleteCodes) > 0 {
		_, err = tx.OperatorAgentLineAllocation.
			Delete().
			Where(
				operatoragentlineallocation.OperatorIDEQ(in.OperatorId),
				operatoragentlineallocation.AgentLineCodeIn(deleteCodes...),
			).
			Exec(l.ctx)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}
	}

	// 批量创建新增的代理子线路分配
	if len(createCodes) > 0 {
		builders := make([]*ent.OperatorAgentLineAllocationCreate, 0, len(createCodes))

		for _, agentLineCode := range createCodes {
			builders = append(
				builders,
				tx.OperatorAgentLineAllocation.
					Create().
					SetOperatorID(in.OperatorId).    // 分站ID
					SetAgentLineCode(agentLineCode), // 代理子线路编码
			)
		}

		_, err = tx.OperatorAgentLineAllocation.
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
	return &agentlineallocation.SaveAgentLineAllocationsResponse{}, nil
}
