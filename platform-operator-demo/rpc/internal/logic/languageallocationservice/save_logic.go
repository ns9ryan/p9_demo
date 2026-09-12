package languageallocationservicelogic

import (
	"context"
	"strings"

	"github.com/zeromicro/go-zero/core/logx"

	"oa.98ent.com/p9/platform-operator/pkg/i18nkey"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/ent"
	"oa.98ent.com/p9/platform-operator/rpc/ent/operatorlanguageallocation"
	"oa.98ent.com/p9/platform-operator/rpc/internal/enterror"
	"oa.98ent.com/p9/platform-operator/rpc/internal/svc"
	"oa.98ent.com/p9/platform-operator/rpc/pb/operator/languageallocation"
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

// Save 保存语言分配
func (l *SaveLogic) Save(in *languageallocation.SaveLanguageAllocationsRequest) (*languageallocation.SaveLanguageAllocationsResponse, error) {
	// 分站ID必须大于0
	if in.OperatorId <= 0 {
		return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
	}

	// 整理语言编码并去重
	languageCodes := make([]string, 0, len(in.LanguageCodes))
	languageCodeSet := make(map[string]struct{}, len(in.LanguageCodes))

	for _, code := range in.LanguageCodes {
		languageCode := strings.TrimSpace(code)
		if languageCode == "" {
			return nil, grpcerror.InvalidArgument(i18nkey.ValidationError)
		}

		if _, exists := languageCodeSet[languageCode]; exists {
			continue
		}

		languageCodeSet[languageCode] = struct{}{}
		languageCodes = append(languageCodes, languageCode)
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

	// 获取当前语言编码
	currentCodes, err := tx.OperatorLanguageAllocation.
		Query().
		Where(operatorlanguageallocation.OperatorIDEQ(in.OperatorId)).
		Select(operatorlanguageallocation.FieldLanguageCode).
		Strings(l.ctx)
	if err != nil {
		// 转换Ent错误为gRPC错误
		return nil, enterror.Handle(l.Logger, err)
	}

	// 整理当前语言编码
	currentCodeSet := make(map[string]struct{}, len(currentCodes))
	for _, languageCode := range currentCodes {
		currentCodeSet[languageCode] = struct{}{}
	}

	// 计算需要删除的语言编码
	deleteCodes := make([]string, 0)
	for _, languageCode := range currentCodes {
		if _, exists := languageCodeSet[languageCode]; !exists {
			deleteCodes = append(deleteCodes, languageCode)
		}
	}

	// 计算需要新增的语言编码
	createCodes := make([]string, 0)
	for _, languageCode := range languageCodes {
		if _, exists := currentCodeSet[languageCode]; !exists {
			createCodes = append(createCodes, languageCode)
		}
	}

	// 批量删除已经取消的语言分配
	if len(deleteCodes) > 0 {
		_, err = tx.OperatorLanguageAllocation.
			Delete().
			Where(
				operatorlanguageallocation.OperatorIDEQ(in.OperatorId),
				operatorlanguageallocation.LanguageCodeIn(deleteCodes...),
			).
			Exec(l.ctx)
		if err != nil {
			// 转换Ent错误为gRPC错误
			return nil, enterror.Handle(l.Logger, err)
		}
	}

	// 批量创建新增的语言分配
	if len(createCodes) > 0 {
		builders := make([]*ent.OperatorLanguageAllocationCreate, 0, len(createCodes))

		for _, languageCode := range createCodes {
			builders = append(
				builders,
				tx.OperatorLanguageAllocation.
					Create().
					SetOperatorID(in.OperatorId).  // 分站ID
					SetLanguageCode(languageCode), // 语言编码
			)
		}

		_, err = tx.OperatorLanguageAllocation.
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
	return &languageallocation.SaveLanguageAllocationsResponse{}, nil
}
