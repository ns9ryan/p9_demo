// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.2

package language_allocation

import (
	"context"
	"strings"

	"oa.98ent.com/p9/core/rpc/coreclient"
	"oa.98ent.com/p9/platform-operator/api/internal/i18nkey"
	"oa.98ent.com/p9/platform-operator/api/internal/svc"
	"oa.98ent.com/p9/platform-operator/api/internal/types"
	"oa.98ent.com/p9/platform-operator/pkg/rpc/grpcerror"
	"oa.98ent.com/p9/platform-operator/rpc/pb/platformoperatorrpc/languageallocationpb"

	"github.com/zeromicro/go-zero/core/logx"
)

const coreLanguagePageSize int32 = 100

type SaveLanguageAllocationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSaveLanguageAllocationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SaveLanguageAllocationsLogic {
	return &SaveLanguageAllocationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// SaveLanguageAllocations 保存语言分配
func (l *SaveLanguageAllocationsLogic) SaveLanguageAllocations(req *types.SaveLanguageAllocationsRequest) (resp *types.SaveLanguageAllocationsResponse, err error) {
	languageCodes := make([]string, 0, len(req.LanguageCodes))

	if len(req.LanguageCodes) > 0 {
		// 获取Core全部语言主数据
		languages, err := l.getCoreLanguages()
		if err != nil {
			return nil, err
		}

		// 建立可用语言编码索引
		languageMap := make(map[string]string, len(languages))
		for _, language := range languages {
			// 已停用语言不能继续分配
			if language.Disabled != 0 {
				continue
			}

			languageMap[strings.ToLower(language.Lang)] = language.Lang
		}

		// 校验并转换为Core标准语言编码
		for _, languageCode := range req.LanguageCodes {
			code := strings.ToLower(strings.TrimSpace(languageCode))

			standardCode, ok := languageMap[code]
			if !ok {
				return nil, grpcerror.InvalidArgument(i18nkey.LanguageUnavailable)
			}

			languageCodes = append(languageCodes, standardCode)
		}
	}

	// 保存分站语言分配
	_, err = l.svcCtx.LanguageAllocationRpc.Save(
		l.ctx,
		&languageallocationpb.SaveLanguageAllocationsRequest{
			OperatorId:    req.OperatorId, // 分站ID
			LanguageCodes: languageCodes,  // 当前分配的语言编码
		},
	)
	if err != nil {
		return nil, err
	}

	// 返回保存结果
	return &types.SaveLanguageAllocationsResponse{}, nil
}

// getCoreLanguages 获取Core全部语言主数据 todo：这个要改，需要core提供获取所有语言
func (l *SaveLanguageAllocationsLogic) getCoreLanguages() ([]*coreclient.I18NLangInfo, error) {
	page := int32(1)
	languages := make([]*coreclient.I18NLangInfo, 0)

	for {
		// 分页获取Core语言主数据
		result, err := l.svcCtx.Core.GetI18NLangList(
			l.ctx,
			&coreclient.I18NLangListReq{
				Page:     page,                 // 页码, 从1开始
				PageSize: coreLanguagePageSize, // 每页数量
			},
		)
		if err != nil {
			return nil, err
		}

		languages = append(languages, result.List...)

		// 已获取全部语言
		if int64(len(languages)) >= result.Total || len(result.List) == 0 {
			break
		}

		page++
	}

	return languages, nil
}
