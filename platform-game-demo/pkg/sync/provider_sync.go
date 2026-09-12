package sync

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 提供商同步服务 =====

// ProviderSyncService 提供商同步服务
type ProviderSyncService struct {
	db            *gorm.DB
	progressQueue *ProgressQueue
}

// NewProviderSyncService 创建提供商同步服务
func NewProviderSyncService(db *gorm.DB) *ProviderSyncService {
	return &ProviderSyncService{
		db:            db,
		progressQueue: NewProgressQueue(db),
	}
}

// Preview 预检查提供商同步
func (s *ProviderSyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteProviderResponse, error) {
	// 获取远程提供商数据
	remoteResp, err := getRemoteGameProviders(ctx, client, isGetRemoteClient)
	if err != nil {
		logger.Errorf("[提供商预检查] 获取测试数据失败: %v", err)
		return nil, nil, err
	}
	remote := remoteResp.Data

	// 获取本地提供商数据
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 转换为接口切片
	remoteItems := make([]interface{}, len(remote))
	for i, v := range remote {
		remoteItems[i] = v
	}

	// 转换本地索引为接口类型
	localIndexInterface := make(map[string]interface{})
	for k, v := range localIndex {
		localIndexInterface[k] = v
	}

	// 当 req 为 nil 时，返回全量结果不做分页
	if req == nil {
		handler := &SyncPageHandler{
			RemoteData: remoteItems,
			LocalIndex: localIndexInterface,
			CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
				remoteProvider := item.(*vendors.VendorInfo)
				localIdx := make(map[string]*ent.GameProvider)
				for k, v := range localIndex {
					localIdx[k] = v.(*ent.GameProvider)
				}
				return s.compareOne(remoteProvider, localIdx)
			},
			PageSize: 0,
			Page:     0,
			SkipNoop: false,
		}
		result := handler.Handle(int64(len(remote)), int64(len(local)))
		return result, remoteResp, nil
	}

	// 创建分页处理器
	handler := &SyncPageHandler{
		RemoteData: remoteItems,
		LocalIndex: localIndexInterface,
		CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
			remoteProvider := item.(*vendors.VendorInfo)
			localIdx := make(map[string]*ent.GameProvider)
			for k, v := range localIndex {
				localIdx[k] = v.(*ent.GameProvider)
			}
			return s.compareOne(remoteProvider, localIdx)
		},
		PageSize: req.PageSize,
		Page:     req.Page,
		SkipNoop: req.IsSkip,
	}

	// 处理分页逻辑
	result := handler.Handle(int64(len(remote)), int64(len(local)))

	return result, remoteResp, nil
}

// 处理提供商数据的创建和更新（从 Run 提取的业务逻辑）
func (s *ProviderSyncService) processProviderDataWithProgress(ctx context.Context, tx *gorm.DB, remoteData []*vendors.VendorInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
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
		var localProvider ent.GameProvider
		exists := tx.Where("source_provider_code = ?", remoteProvider.Code).
			Where("deleted_at IS NULL").
			First(&localProvider).Error == nil
		nameI18n := map[string]interface{}{"default": remoteProvider.Name}
		sourceNameI18n := map[string]interface{}{"default": remoteProvider.Name}
		// 生成 P9 内部的提供商编码
		providerCode := remoteProvider.Code + "_" + fmt.Sprintf("%d", i)
		if !exists {
			newProvider := ent.GameProvider{
				SourceId:           remoteProvider.Id,
				ProviderCode:       providerCode,
				SourceProviderCode: remoteProvider.Code,
				NameI18n:           JSONToString(nameI18n),
				SourceNameI18n:     JSONToString(sourceNameI18n),
				SortNo:             int64(remoteProvider.Id),
				SourceSortNo:       sql.NullInt64{Int64: int64(0), Valid: true},
				SourceStatus:       int64(remoteProvider.Status),
				Status:             int64(remoteProvider.Status),
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
			updateData := map[string]interface{}{
				"provider_code":        providerCode,
				"source_provider_code": remoteProvider.Code,
				"source_name_i18n":     sourceNameI18n,
				"source_status":        remoteProvider.Status,
				"source_sort_no":       0,
				"updated_at":           time.Now(),
			}
			if len(syncCols) > 0 {
				updateData = GetUpdatedData(map[string]interface{}{
					"provider_code":        providerCode,
					"source_provider_code": remoteProvider.Code,
					"name_i18n":            nameI18n,
					"source_name_i18n":     sourceNameI18n,
					"sort_no":              int64(remoteProvider.Id),
					"source_sort_no":       0,
					"status":               int64(remoteProvider.Status),
					"source_status":        int64(remoteProvider.Status),
				}, syncCols)
				updateData["updated_at"] = time.Now()
			}
			if err := tx.Model(&localProvider).
				Updates(updateData).Error; err != nil {
				logger.Errorf("[提供商更新] 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Updated++
		}
	}

	return nil
}

func (s *ProviderSyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	logger.Infof("[提供商同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[提供商同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		logger.Errorf("[提供商同步] 预检查失败: %v", err)
		return err
	}
	// 创建并启动进度队列
	s.progressQueue.Start(ctx)
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal {
		defer s.progressQueue.Stop()
		logger.Infof("[提供商同步] 无需更新，直接返回")
		logger.Infof("[提供商同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
			previewResp.Stats.UpdateTotal, previewResp.Stats.CreateTotal, previewResp.Stats.DeleteTotal, previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
		// 更新checkpointID进度为100
		s.progressQueue.Send(&ProgressMessage{
			TableName:      "provider",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
		})
		return nil
	}
	logger.Infof("[提供商同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
	// 创建初始检查点记录（此时progress=0）
	logger.Infof("[提供商同步] 创建初始检查点")
	logger.Infof("[提供商同步] ✓ 检查点已创建: checkpointID=%d", checkpointID)

	go func() {
		defer s.progressQueue.Stop()
		// 在事务中执行数据库操作
		logger.Infof("[提供商同步] 开始处理提供商数据")
		// 把remoteResp.Data拆分为100条一批进行处理（可根据实际情况调整批次大小）
		batchSize := BatchSize
		totalCount := len(remoteResp.Data)
		apply := &Apply{}
		for i := 0; i < len(remoteResp.Data); i += batchSize {
			end := i + batchSize
			if end > len(remoteResp.Data) {
				end = len(remoteResp.Data)
			}
			batch := remoteResp.Data[i:end]
			logger.Infof("[提供商同步] 处理批次: start=%d, end=%d", i, end)
			// 模拟等待
			logger.Debugf("[进度队列] 模拟处理延迟: 表=provider, checkpointID=%d", checkpointID)
			time.Sleep(1000 * time.Millisecond)
			logger.Debugf("[进度队列] 开始处理消息: 表=provider, checkpointID=%d", checkpointID)
			err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				return s.processProviderDataWithProgress(ctx, tx, batch, apply, syncCols, checkpointID)
			})
			if err != nil {
				logger.Errorf("[提供商同步] 处理提供商数据失败: %v", err)
				return
			}
			progress := CalculateProgress(int64(i+1), int64(totalCount))
			msg := &ProgressMessage{
				TableName:      "provider",
				ProcessedCount: int32(i + 1),
				RemoteTotal:    previewResp.Stats.RemoteTotal,
				LocalTotal:     previewResp.Stats.LocalTotal,
				Progress:       progress,
				CheckpointID:   checkpointID,
				Created:        apply.Created,
				Updated:        apply.Updated,
				Deleted:        apply.Deleted,
				Failed:         apply.Failed,
				Skipped:        apply.Skipped,
			}
			logger.Infof("[提供商同步] 发送进度消息: 表=provider, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
			s.progressQueue.Send(msg)
		}
		if err != nil {
			logger.Errorf("[提供商同步] 处理提供商数据失败: %v", err)
			return
		}

		// 执行逆向同步
		logger.Infof("[提供商同步] 开始执行逆向同步")
		err := s.ReverseSync(ctx, client, remoteResp, apply)
		if err != nil {
			logger.Errorf("❌ [提供商同步] 逆向同步失败: %v", err)
			return
		}
		logger.Infof("[提供商同步] 逆向同步完成: deleted=%d", apply.Deleted)

		// 最后一次进度更新（100%）
		logger.Infof("[提供商同步] 发送最终进度消息（100%%）")
		finalMsg := &ProgressMessage{
			TableName:      "provider",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
			Created:        apply.Created,
			Updated:        apply.Updated,
			Deleted:        apply.Deleted,
			Failed:         apply.Failed,
			Skipped:        apply.Skipped,
		}
		s.progressQueue.Send(finalMsg)

		logger.Infof("[提供商同步] ===== 同步完成 =====")
	}()
	// 构建执行结果
	return nil
}

// fetchLocal 获取本地所有提供商
func (s *ProviderSyncService) fetchLocal(ctx context.Context) ([]*ent.GameProvider, map[string]*ent.GameProvider, error) {
	var providers []*ent.GameProvider
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&providers).Error; err != nil {
		return nil, nil, err
	}

	// 构建本地提供商索引（按上游提供商编码）
	index := make(map[string]*ent.GameProvider)
	for _, provider := range providers {
		if provider.SourceProviderCode != "" {
			index[provider.SourceProviderCode] = provider
		}
	}

	return providers, index, nil
}

// 比较单个提供商的本地和远程数据
func (s *ProviderSyncService) compareOne(remote *vendors.VendorInfo, localIndex map[string]*ent.GameProvider) *platform_game.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &platform_game.SyncDiff{
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

	if CompareName(local.SourceNameI18n, remote.Name) {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("提供商名称变更: %s -> %s", local.SourceNameI18n, remote.Name)
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *ProviderSyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteProviderResponse, apply *Apply) error {
	logger.Info("[提供商逆向同步] 开始执行逆向同步")

	// 构建远程提供商编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteProvider := range remoteResp.Data {
		remoteIndex[remoteProvider.Code] = true
	}

	// 获取本地提供商数据
	var localProviders []*ent.GameProvider
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localProviders).Error; err != nil {
		logger.Errorf("[提供商逆向同步] ✗ 获取本地数据失败: %v", err)
		return err
	}
	logger.Infof("[提供商逆向同步] ✓ 获取本地数据成功, 共 %d 条", len(localProviders))

	// 在事务中执行软删除
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localProvider := range localProviders {
			// 如果本地渠道在远程不存在，则软删除
			if !remoteIndex[localProvider.SourceProviderCode] {
				if err := tx.Model(localProvider).Update("deleted_at", time.Now()).Error; err != nil {
					logger.Errorf("[提供商逆向同步] 软删除失败: %v", err)
					apply.Failed++
					continue
				}
				apply.Deleted++
				logger.Infof("[提供商逆向同步] ✓ 已软删除提供商: %s", localProvider.SourceProviderCode)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	logger.Infof("[提供商逆向同步] ✓ 完成，删除: %d", apply.Deleted)
	return nil
}

func getRemoteGameProviders(ctx context.Context, client vendors.VendorGameServiceClient, isGetRemoteClient bool) (*RemoteProviderResponse, error) {
	if isGetRemoteClient {
		logger.Infof("[同步数据源] 使用远程客户端获取数据")
		remoteResp, err := client.GetVendor(ctx, &vendors.Empty{})
		if err != nil {
			logger.Errorf("❌ [同步数据源] 获取远程数据失败: %v", err)
			return nil, err
		}
		logger.Infof("[同步数据源] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.VendorList))
		return &RemoteProviderResponse{Data: remoteResp.VendorList}, nil
	} else {
		logger.Infof("[同步数据源] 使用本地 JSON 数据获取提供商信息")
		var list []*vendors.VendorInfo
		if err := loadLocalJSON("provider.json", &list); err != nil {
			logger.Errorf("[同步数据源] 读取 provider.json 失败: %v", err)
			return nil, err
		}
		logger.Infof("[同步数据源] ✓ 读取本地 JSON 数据成功, 共 %d 条", len(list))
		return &RemoteProviderResponse{Data: list}, nil
	}
}

type RemoteProviderResponse struct {
	Data []*vendors.VendorInfo `json:"data"`
}
