package sync

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 提供商同步服务 =====

// ProviderSyncService 提供商同步服务
type ProviderSyncService struct {
	db *gorm.DB
}

// NewProviderSyncService 创建提供商同步服务
func NewProviderSyncService(db *gorm.DB) *ProviderSyncService {
	return &ProviderSyncService{db: db}
}

// Preview 预检查提供商同步
func (s *ProviderSyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncPreviewResp, error) {
	// 获取远程提供商数据
	// remoteResp, err := client.GetVendor(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("[提供商预检查] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameProviders()
	if err != nil {
		logger.Errorf("[提供商预检查] 获取测试数据失败: %v", err)
		return nil, err
	}
	remote := remoteResp.Data

	// 获取本地提供商数据
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		return nil, err
	}

	// 初始化预检查结果
	result := &vendors.SyncPreviewResp{
		Stats: &vendors.SyncStats{
			RemoteTotal: int64(len(remote)),
			LocalTotal:  int64(len(local)),
		},
		Diffs: make([]*vendors.SyncDiff, 0),
	}

	// 对每条远程提供商进行比较
	for _, remoteItem := range remote {
		diff := s.compareOneByCode(remoteItem, localIndex)
		result.Diffs = append(result.Diffs, diff)

		// 按操作类型统计
		switch diff.Action {
		case "create":
			result.Stats.CreateTotal++
		case "update":
			result.Stats.UpdateTotal++
		case "noop":
			result.Stats.NoopTotal++
		case "conflict":
			result.Stats.ConflictTotal++
		}
	}

	return result, nil
}

// Run 执行提供商同步
// processProviderData 处理提供商数据的创建和更新（从 Run 提取的业务逻辑）
func (s *ProviderSyncService) processProviderData(ctx context.Context, tx *gorm.DB, remoteData []*vendors.VendorInfo, applyResult *vendors.SyncApplyResult) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 source_provider_code 的计数，确保唯一性
	for _, remoteProvider := range remoteData {
		if counterMap[remoteProvider.Code] > 0 {
			logger.Errorf("[提供商同步] 检测到重复的远程提供商编码: %s, 计数器: %d, 跳过处理", remoteProvider.Code, counterMap[remoteProvider.Code])
			applyResult.Failed++
			continue
		}
		// 更新计数器
		counterMap[remoteProvider.Code]++
		i++
		// 检查本地是否存在该提供商（使用上游提供商编码）
		var localProvider model.Provider
		exists := tx.Where("source_provider_code = ?", remoteProvider.Code).
			Where("deleted_at IS NULL").
			First(&localProvider).Error == nil

		if !exists {
			// 新增提供商 - 生成 P9 内部的提供商编码
			providerCode := remoteProvider.Code + "_" + fmt.Sprintf("%d", i)
			nameI18n := model.JSONMap{"default": remoteProvider.Name}
			sourceNameI18n := model.JSONMap{"default": remoteProvider.Name}

			newProvider := model.Provider{
				SourceID:           remoteProvider.Id,
				ProviderCode:       providerCode,
				SourceProviderCode: remoteProvider.Code,
				SourceNameI18n:     sourceNameI18n,
				NameI18n:           nameI18n,
				SourceStatus:       remoteProvider.Status,
				Status:             remoteProvider.Status,
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			if err := tx.Create(&newProvider).Error; err != nil {
				logger.Errorf("[提供商新增] 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Created++
		} else {
			// 更新提供商信息
			providerCode := remoteProvider.Code + "_" + fmt.Sprintf("%d", i)
			sourceNameI18n := model.JSONMap{"default": remoteProvider.Name}
			if err := tx.Model(&localProvider).
				Updates(map[string]interface{}{
					"provider_code":        providerCode,
					"source_provider_code": remoteProvider.Code,
					"source_name_i18n":     sourceNameI18n,
					"source_status":        remoteProvider.Status,
					"updated_at":           time.Now(),
				}).Error; err != nil {
				logger.Errorf("[提供商更新] 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Updated++
		}
	}
	return nil
}

func (s *ProviderSyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	logger.Infof("[提供商同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[提供商同步] 执行预检查")
	previewResp, err := s.Preview(ctx, client)
	if err != nil {
		logger.Errorf("[提供商同步] 预检查失败: %v", err)
		return nil, err
	}
	logger.Infof("[提供商同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: previewResp,
		Apply:   &vendors.SyncApplyResult{},
	}

	// 获取远程提供商数据
	logger.Infof("[提供商同步] 获取远程提供商数据")
	// remoteResp, err := client.GetVendor(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("[提供商同步] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameProviders()
	if err != nil {
		logger.Errorf("[提供商同步] 获取测试数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[提供商同步] 获取到远程数据: count=%d", len(remoteResp.Data))

	// 在事务中执行数据库操作
	logger.Infof("[提供商同步] 开始处理提供商数据")
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.processProviderData(ctx, tx, remoteResp.Data, result.Apply)
	})
	if err != nil {
		logger.Errorf("[提供商同步] 处理提供商数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[提供商同步] 处理完成: created=%d, updated=%d, deleted=%d, failed=%d",
		result.Apply.Created, result.Apply.Updated, result.Apply.Deleted, result.Apply.Failed)

	// 执行逆向同步
	logger.Infof("[提供商同步] 开始执行逆向同步")
	reverseResp, err := s.ReverseSync(ctx, client)
	if err != nil {
		logger.Errorf("[提供商同步] 逆向同步失败: %v", err)
		return nil, err
	}
	result.Apply.Deleted = reverseResp.Apply.Deleted
	logger.Infof("[提供商同步] 逆向同步完成: deleted=%d", result.Apply.Deleted)

	// 保存同步检查点
	logger.Infof("[提供商同步] 准备保存同步检查点")
	checkpointMgr := NewCheckpointManager(s.db)
	checkpointValue := previewResp.Stats.RemoteTotal // 使用远程总数作为检查点值
	if err := checkpointMgr.UpdateCheckpoint(ctx, "PROVIDER",
		fmt.Sprintf("%d", checkpointValue), result, nil); err != nil {
		logger.Errorf("[提供商同步] 保存检查点失败: %v", err)
		return nil, err
	}
	logger.Infof("[提供商同步] ✓ 检查点已保存")

	logger.Infof("[提供商同步] ===== 同步完成 =====")
	return result, nil
}

// fetchLocal 获取本地所有提供商
func (s *ProviderSyncService) fetchLocal(ctx context.Context) ([]*model.Provider, map[string]*model.Provider, error) {
	var providers []*model.Provider
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&providers).Error; err != nil {
		return nil, nil, err
	}

	// 构建本地提供商索引（按上游提供商编码）
	index := make(map[string]*model.Provider)
	for _, provider := range providers {
		if provider.SourceProviderCode != "" {
			index[provider.SourceProviderCode] = provider
		}
	}

	return providers, index, nil
}

// compareOneByCode 比较单个提供商的本地和远程数据
func (s *ProviderSyncService) compareOneByCode(remote *vendors.VendorInfo, localIndex map[string]*model.Provider) *vendors.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &vendors.SyncDiff{
		ObjectType: "provider",
		ObjectId:   remote.Id,
		ObjectCode: remote.Code,
		RemoteId:   remote.Id,
		RemoteCode: remote.Code,
	}

	// 如果本地不存在该提供商
	if !exists {
		diff.Action = "create"
		diff.Reason = "本地不存在该提供商"
		return diff
	}

	// 对比提供商名称是否变更
	var localSourceName string
	if sourceNameI18n, ok := local.SourceNameI18n["default"]; ok {
		localSourceName = sourceNameI18n.(string)
	}
	if localSourceName != remote.Name {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("提供商名称变更: %s -> %s", localSourceName, remote.Name)
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// ReverseSync 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *ProviderSyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	logger.Info("[提供商逆向同步] 开始执行逆向同步")

	// 获取远程提供商数据
	// remoteResp, err := client.GetVendor(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("[提供商逆向同步] ✗ 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameProviders()
	if err != nil {
		logger.Errorf("[提供商逆向同步] 获取测试数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[提供商逆向同步] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.Data))

	// 构建远程提供商编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteProvider := range remoteResp.Data {
		remoteIndex[remoteProvider.Code] = true
	}

	// 获取本地提供商数据
	var localProviders []*model.Provider
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localProviders).Error; err != nil {
		logger.Errorf("[提供商逆向同步] ✗ 获取本地数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[提供商逆向同步] ✓ 获取本地数据成功, 共 %d 条", len(localProviders))

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: &vendors.SyncPreviewResp{
			Stats: &vendors.SyncStats{
				LocalTotal: int64(len(localProviders)),
			},
		},
		Apply: &vendors.SyncApplyResult{},
	}

	// 东前需要有本地数据按 source_code 的重数统计
	localBySource := make(map[string]int)
	for _, localProvider := range localProviders {
		localBySource[localProvider.SourceProviderCode]++
	}

	// 在事务中执行软删除
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localProvider := range localProviders {
			duplicateCount := localBySource[localProvider.SourceProviderCode]

			// 如果同一 source_code 有多个记录，进行软删除
			if duplicateCount > 1 {
				logger.Infof("[提供商逆向同步] 检测到重复记录: source_provider_code=%s, 重复数=%d, 作软删除处理",
					localProvider.SourceProviderCode, duplicateCount)

				if err := tx.Model(localProvider).Update("deleted_at", time.Now()).Error; err != nil {
					logger.Errorf("[提供商逆向同步] 软删除失败: %v", err)
					result.Apply.Failed++
					continue
				}
				result.Apply.Deleted++
				logger.Infof("[提供商逆向同步] ✓ 已软删除重复记录: %s (ID: %d)", localProvider.SourceProviderCode, localProvider.ID)
			} else if !remoteIndex[localProvider.SourceProviderCode] {
				// 此记录不重复，但本地数据在远程不存在，需要软删除
				logger.Infof("[提供商逆向同步] 检测到本地孤立记录: %s (ID: %d), 准备软删除",
					localProvider.SourceProviderCode, localProvider.ID)

				if err := tx.Model(localProvider).Update("deleted_at", time.Now()).Error; err != nil {
					logger.Errorf("[提供商逆向同步] 软删除失败: %v", err)
					result.Apply.Failed++
					continue
				}
				result.Apply.Deleted++
				logger.Infof("[提供商逆向同步] ✓ 已软删除孤立记录: %s (ID: %d)", localProvider.SourceProviderCode, localProvider.ID)
			} else {
				logger.Debugf("[提供商逆向同步] 本地记录在远程存在，且不重复，无需删除: %s", localProvider.SourceProviderCode)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Infof("[提供商逆向同步] ✓ 完成，删除: %d", result.Apply.Deleted)
	return result, nil
}

func getTestGameProviders() (*vendors.GetVendorResp, error) {
	var list []*vendors.VendorInfo
	if err := loadTestJSON("provider.json", &list); err != nil {
		logger.Errorf("[提供商测试数据] 读取 provider.json 失败: %v", err)
		return &vendors.GetVendorResp{Data: []*vendors.VendorInfo{}}, err
	}

	return &vendors.GetVendorResp{Data: list}, nil
}
