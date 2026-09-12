package sync

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	platform_game_sync "oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 分类同步服务 =====

// CategorySyncService 分类同步服务
type CategorySyncService struct {
	db            *gorm.DB
	progressQueue *ProgressQueue
	mu            sync.Mutex
}

// NewCategorySyncService 创建分类同步服务
func NewCategorySyncService(db *gorm.DB) *CategorySyncService {
	return &CategorySyncService{
		db:            db,
		progressQueue: NewProgressQueue(db),
	}
}

// 预检查分类同步
func (s *CategorySyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game_sync.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteCategoryResponse, error) {
	logger.Info("[分类预检查] 开始执行")

	// 获取远程分类数据
	logger.Info("[分类预检查] 正在调用 GetGameCategory RPC")
	remoteResp, err := getRemoteGameCategory(ctx, client, isGetRemoteClient)
	if err != nil {
		logger.Errorf("❌ [分类预检查] 获取测试数据失败: %v", err)
		return nil, nil, err
	}

	logger.Infof("[分类预检查] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.Data))

	// 获取本地分类数据
	logger.Info("[分类预检查] 正在获取本地数据")
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		logger.Errorf("❌ [分类预检查] 获取本地数据失败: %v", err)
		return nil, nil, err
	}
	logger.Infof("[分类预检查] ✓ 获取本地数据成功, 共 %d 条", len(local))

	// 转换为接口切片
	remoteItems := make([]interface{}, len(remoteResp.Data))
	for i, v := range remoteResp.Data {
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
				remoteCategory := item.(*vendors.GameCategoryInfo)
				localIdx := make(map[string]*ent.GameCategory)
				for k, v := range localIndex {
					localIdx[k] = v.(*ent.GameCategory)
				}
				return s.compareOne(remoteCategory, localIdx)
			},
			PageSize: 0,
			Page:     0,
			SkipNoop: false,
		}
		result := handler.Handle(int64(len(remoteResp.Data)), int64(len(local)))
		return result, remoteResp, nil
	}

	// 创建分页处理器
	handler := &SyncPageHandler{
		RemoteData: remoteItems,
		LocalIndex: localIndexInterface,
		CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
			remoteCategory := item.(*vendors.GameCategoryInfo)
			localIdx := make(map[string]*ent.GameCategory)
			for k, v := range localIndex {
				localIdx[k] = v.(*ent.GameCategory)
			}
			return s.compareOne(remoteCategory, localIdx)
		},
		PageSize: req.PageSize,
		Page:     req.Page,
		SkipNoop: req.IsSkip,
	}

	// 处理分页逻辑
	result := handler.Handle(int64(len(remoteResp.Data)), int64(len(local)))
	return result, remoteResp, nil
}

// processCategoryData 处理分类数据的创建和更新（从 Run 提取的业务逻辑）
func (s *CategorySyncService) processCategoryDataWithProgress(ctx context.Context, tx *gorm.DB, remoteData []*vendors.GameCategoryInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 source_category_code 的计数，确保唯一性

	for _, remoteCat := range remoteData {
		if counterMap[remoteCat.Code] > 0 {
			logger.Errorf("[分类同步] 检测到重复的远程分类编码: %s, 计数器: %d, 跳过处理", remoteCat.Code, counterMap[remoteCat.Code])
			applyResult.Failed++
			continue
		}
		// 更新计数器
		counterMap[remoteCat.Code]++
		i++
		// 检查本地是否存在该分类（使用上游分类编码）
		var localCat ent.GameCategory
		exists := tx.Where("source_category_code = ?", remoteCat.Code).
			Where("deleted_at IS NULL").
			First(&localCat).Error == nil
		categoryCode := remoteCat.Code + "_" + fmt.Sprintf("%d", i)
		nameI18n := map[string]interface{}{"default": remoteCat.Name}
		sourceNameI18n := map[string]interface{}{"default": remoteCat.Name}
		if !exists {
			// 新增分类 - 生成 P9 内部的分类编码
			newCat := ent.GameCategory{
				SourceId:           remoteCat.Id,
				CategoryCode:       categoryCode,
				SourceCategoryCode: remoteCat.Code,
				NameI18n:           JSONToString(nameI18n),
				SourceNameI18n:     JSONToString(sourceNameI18n),
				SortNo:             remoteCat.Id,
				SourceSortNo:       sql.NullInt64{Int64: 0, Valid: true},
				Status:             int64(remoteCat.Status),
				SourceStatus:       int64(remoteCat.Status),
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			}
			if err := tx.Create(&newCat).Error; err != nil {
				logger.Errorf("[分类新增] 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Created++
			logger.Infof("[分类新增] ✓ 已创建新分类: %s", remoteCat.Code)
		} else {
			// 更新分类信息
			updateData := map[string]interface{}{
				"category_code":        categoryCode,
				"source_category_code": remoteCat.Code,
				"source_name_i18n":     JSONToString(sourceNameI18n),
				"source_sort_no":       0,
				"source_status":        int64(remoteCat.Status),
			}
			if len(syncCols) > 0 {
				updateData = GetUpdatedData(map[string]interface{}{
					"category_code":        categoryCode,
					"source_category_code": remoteCat.Code,
					"name_i18n":            JSONToString(nameI18n),
					"source_name_i18n":     JSONToString(sourceNameI18n),
					"sort_no":              remoteCat.Id,
					"source_sort_no":       0,
					"status":               int64(remoteCat.Status),
					"source_status":        int64(remoteCat.Status),
				}, syncCols)
				updateData["updated_at"] = time.Now()
			}
			if err := tx.Model(&localCat).
				Updates(updateData).Error; err != nil {
				logger.Errorf("[分类更新] 失败: %v", err)
				applyResult.Failed++
				continue
			}
			applyResult.Updated++
			logger.Infof("[分类更新] ✓ 已更新分类: %s", remoteCat.Code)
		}

	}

	return nil
}

// Run 执行分类同步
func (s *CategorySyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	logger.Infof("[分类同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[分类同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		logger.Errorf("[分类同步] 预检查失败: %v", err)
		return err
	}
	s.progressQueue.Start(ctx)
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal {
		defer s.progressQueue.Stop()
		logger.Infof("[分类同步] 无需更新，直接返回")
		logger.Infof("[分类同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
			previewResp.Stats.UpdateTotal, previewResp.Stats.CreateTotal, previewResp.Stats.DeleteTotal, previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
		// 更新checkpointID进度为100
		s.progressQueue.Send(&ProgressMessage{
			TableName:      "category",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
		})
		return nil
	}
	logger.Infof("[分类同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 获取远程分类数据
	logger.Infof("[分类同步] 获取远程分类数据")
	// remoteResp 已经从 Preview 获取，无需再次获取
	logger.Infof("[分类同步] 获取到远程数据: count=%d", len(remoteResp.Data))

	// 创建初始检查点记录（此时progress=0）
	logger.Infof("[分类同步] 创建初始检查点 (syncCols=%v)", syncCols)
	logger.Infof("[分类同步] ✓ 检查点已创建: checkpointID=%d", checkpointID)

	go func() {
		defer s.progressQueue.Stop()
		// 在事务中执行数据库操作
		logger.Infof("[分类同步] 开始处理分类数据")
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
			logger.Infof("[分类同步] 处理批次: start=%d, end=%d", i, end)
			// 模拟等待
			logger.Debugf("[进度队列] 模拟处理延迟: 表=category, checkpointID=%d", checkpointID)
			time.Sleep(1000 * time.Millisecond)
			logger.Debugf("[进度队列] 开始处理消息: 表=category, checkpointID=%d", checkpointID)
			err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				return s.processCategoryDataWithProgress(ctx, tx, batch, apply, syncCols, checkpointID)
			})
			if err != nil {
				logger.Errorf("[分类同步] 处理分类数据失败: %v", err)
				return
			}
			progress := CalculateProgress(int64(i+1), int64(totalCount))
			msg := &ProgressMessage{
				TableName:      "category",
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
			logger.Infof("[分类同步] 发送进度消息: 表=category, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
			s.progressQueue.Send(msg)
		}
		if err != nil {
			logger.Errorf("[分类同步] 处理分类数据失败: %v", err)
			return
		}

		// 执行逆向同步
		logger.Infof("[分类同步] 开始执行逆向同步")
		_, err := s.ReverseSync(ctx, client, remoteResp, apply)
		if err != nil {
			logger.Errorf("❌ [分类同步] 逆向同步失败: %v", err)
			return
		}
		logger.Infof("[分类同步] 逆向同步完成: deleted=%d", apply.Deleted)

		// 最后一次进度更新（100%）
		logger.Infof("[分类同步] 发送最终进度消息（100%%）")
		finalMsg := &ProgressMessage{
			TableName:      "category",
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

		logger.Infof("[分类同步] ===== 同步完成 =====")
	}()

	return nil
}

// fetchLocal 获取本地所有分类
func (s *CategorySyncService) fetchLocal(ctx context.Context) ([]*ent.GameCategory, map[string]*ent.GameCategory, error) {
	var categories []*ent.GameCategory
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	// 构建本地分类索引（按上游分类编码）
	index := make(map[string]*ent.GameCategory)
	for _, cat := range categories {
		index[cat.SourceCategoryCode] = cat
	}

	return categories, index, nil
}

// compareOne 比较单个分类的本地和远程数据
func (s *CategorySyncService) compareOne(remote *vendors.GameCategoryInfo, localIndex map[string]*ent.GameCategory) *platform_game.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &platform_game.SyncDiff{
		ObjectType: "category",
		ObjectId:   remote.Id,
		ObjectCode: remote.Code,
		RemoteId:   remote.Id,
		RemoteCode: remote.Code,
	}

	// 如果本地不存在该分类
	if !exists {
		diff.Action = "create"
		diff.Reason = "本地不存在该分类"
		return diff
	}

	if CompareName(local.SourceNameI18n, remote.Name) {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("分类名称变更: %s -> %s", local.SourceNameI18n, remote.Name)
		diff.ConflictType = "name_changed"
		return diff
	}

	if local.SourceCategoryCode != remote.Code {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("分类编码变更: %s -> %s", local.SourceCategoryCode, remote.Code)
		diff.ConflictType = "code_changed"
		return diff
	}

	if local.SourceStatus != int64(remote.Status) {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("分类状态变更: %d -> %d", local.SourceStatus, remote.Status)
		diff.ConflictType = "status_changed"
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *CategorySyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteCategoryResponse, apply *Apply) (*Apply, error) {
	logger.Info("[分类逆向同步] 开始执行逆向同步")

	// 获取远程分类数据
	logger.Infof("[分类逆向同步] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.Data))

	// 构建远程分类编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteCat := range remoteResp.Data {
		remoteIndex[remoteCat.Code] = true
		logger.Debugf("[分类逆向同步] 远程分类编码: %s", remoteCat.Code)
	}
	logger.Infof("[分类逆向同步] 远程分类编码总数: %d", len(remoteIndex))

	// 获取本地分类数据
	var localCategories []*ent.GameCategory
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localCategories).Error; err != nil {
		logger.Errorf("❌ [分类逆向同步] 获取本地数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类逆向同步] ✓ 获取本地数据成功, 共 %d 条", len(localCategories))

	// 打印本地所有分类编码
	for _, localCat := range localCategories {
		logger.Debugf("[分类逆向同步] 本地分类编码: %s (ID: %d, deleted_at: %v)",
			localCat.SourceCategoryCode, localCat.Id, localCat.DeletedAt)
	}

	// 在事务中执行软删除
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localCat := range localCategories {
			// 如果本地分类在远程不存在，则软删除
			if !remoteIndex[localCat.SourceCategoryCode] {
				if err := tx.Model(localCat).Update("deleted_at", time.Now()).Error; err != nil {
					logger.Errorf("❌ [分类逆向同步] 软删除失败: %v", err)
					apply.Failed++
					continue
				}
				apply.Deleted++
				logger.Infof("[分类逆向同步] ✓ 已软删除分类: %s (ID: %d)", localCat.SourceCategoryCode, localCat.Id)
			} else {
				logger.Debugf("[分类逆向同步] 本地分类在远程存在，无需删除: %s", localCat.SourceCategoryCode)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Infof("[分类逆向同步] ✓ 完成，删除: %d", apply.Deleted)
	return apply, nil
}

func getRemoteGameCategory(ctx context.Context, client vendors.VendorGameServiceClient, isGetRemoteClient bool) (*RemoteCategoryResponse, error) {
	if isGetRemoteClient {
		logger.Infof("[同步数据源] 正在获取远程分类数据...")
		remoteResp, err := client.GetGameCategory(ctx, &vendors.Empty{})
		if err != nil {
			logger.Errorf("❌ [同步数据源] 获取远程数据失败: %v", err)
			return nil, err
		}
		logger.Infof("[同步数据源] ✓ 获取远程分类数据成功, 共 %d 条", len(remoteResp.CategoryList))
		return &RemoteCategoryResponse{Data: remoteResp.CategoryList}, nil
	} else {
		logger.Infof("[同步数据源] 正在获取本地分类数据...")
		var list []*vendors.GameCategoryInfo
		if err := loadLocalJSON("category.json", &list); err != nil {
			logger.Errorf("[同步数据源] 读取 category.json 失败: %v", err)
			return &RemoteCategoryResponse{Data: []*vendors.GameCategoryInfo{}}, err
		}

		// 转换 vendors.GameCategoryInfo 为 platformgame.CategoryInfo
		categoryList := make([]*vendors.GameCategoryInfo, 0, len(list))
		for _, item := range list {
			categoryList = append(categoryList, &vendors.GameCategoryInfo{
				Id:      item.Id,
				Code:    item.Code,
				Name:    item.Name,
				Status:  item.Status,
				Channel: item.Channel,
			})
		}
		logger.Infof("[同步数据源] ✓ 获取本地分类数据成功, 共 %d 条", len(categoryList))
		return &RemoteCategoryResponse{Data: categoryList}, nil
	}
}

type RemoteCategoryResponse struct {
	Data []*vendors.GameCategoryInfo `json:"data"`
}
