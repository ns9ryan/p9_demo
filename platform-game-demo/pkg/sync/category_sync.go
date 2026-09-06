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

// ===== 分类同步服务 =====

// CategorySyncService 分类同步服务
type CategorySyncService struct {
	db *gorm.DB
}

// NewCategorySyncService 创建分类同步服务
func NewCategorySyncService(db *gorm.DB) *CategorySyncService {
	return &CategorySyncService{db: db}
}

// Preview 预检查分类同步
func (s *CategorySyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncPreviewResp, error) {
	logger.Info("[分类预检查] 开始执行")

	// 获取远程分类数据
	logger.Info("[分类预检查] 正在调用 GetGameCategory RPC")
	// remoteResp, err := client.GetGameCategory(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("❌ [分类预检查] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameCategory()
	if err != nil {
		logger.Errorf("❌ [分类预检查] 获取测试数据失败: %v", err)
		return nil, err
	}

	logger.Infof("[分类预检查] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.Data))
	remote := remoteResp.Data

	// 获取本地分类数据
	logger.Info("[分类预检查] 正在获取本地数据")
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		logger.Errorf("❌ [分类预检查] 获取本地数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类预检查] ✓ 获取本地数据成功, 共 %d 条", len(local))

	// 初始化预检查结果
	result := &vendors.SyncPreviewResp{
		Stats: &vendors.SyncStats{
			RemoteTotal: int64(len(remote)),
			LocalTotal:  int64(len(local)),
		},
		Diffs: make([]*vendors.SyncDiff, 0),
	}

	// 对每条远程分类进行比较
	for _, remoteItem := range remote {
		diff := s.compareOne(remoteItem, localIndex)
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

// processCategoryData 处理分类数据的创建和更新（从 Run 提取的业务逻辑）
func (s *CategorySyncService) processCategoryData(ctx context.Context, tx *gorm.DB, remoteData []*vendors.GameCategoryInfo, applyResult *vendors.SyncApplyResult) error {
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
		var localCat model.Category
		exists := tx.Where("source_category_code = ?", remoteCat.Code).
			Where("deleted_at IS NULL").
			First(&localCat).Error == nil

		if !exists {
			// 新增分类 - 生成 P9 内部的分类编码
			categoryCode := remoteCat.Code + "_" + fmt.Sprintf("%d", i)
			nameI18n := model.JSONMap{"default": remoteCat.Name}
			sourceNameI18n := model.JSONMap{"default": remoteCat.Name}

			newCat := model.Category{
				SourceID:           remoteCat.Id,
				CategoryCode:       categoryCode,
				SourceCategoryCode: remoteCat.Code,
				NameI18n:           nameI18n,
				SourceNameI18n:     sourceNameI18n,
				SortNo:             remoteCat.Weight,
				SourceSortNo:       remoteCat.Weight,
				Status:             remoteCat.Status,
				SourceStatus:       remoteCat.Status,
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
			categoryCode := remoteCat.Code + "_" + fmt.Sprintf("%d", i)
			sourceNameI18n := model.JSONMap{"default": remoteCat.Name}
			if err := tx.Model(&localCat).
				Updates(map[string]interface{}{
					"category_code":        categoryCode,
					"source_category_code": remoteCat.Code,
					"source_name_i18n":     sourceNameI18n,
					"source_sort_no":       remoteCat.Weight,
					"source_status":        remoteCat.Status,
					"updated_at":           time.Now(),
				}).Error; err != nil {
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
func (s *CategorySyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	logger.Infof("[分类同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[分类同步] 执行预检查")
	previewResp, err := s.Preview(ctx, client)
	if err != nil {
		logger.Errorf("[分类同步] 预检查失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: previewResp,
		Apply:   &vendors.SyncApplyResult{},
	}

	// 获取远程分类数据
	logger.Infof("[分类同步] 获取远程分类数据")
	// remoteResp, err := client.GetGameCategory(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("[分类同步] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameCategory()
	if err != nil {
		logger.Errorf("[分类同步] 获取远程数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类同步] 获取到远程数据: count=%d", len(remoteResp.Data))

	// 在事务中执行数据库操作
	logger.Infof("[分类同步] 开始处理分类数据")
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.processCategoryData(ctx, tx, remoteResp.Data, result.Apply)
	})
	if err != nil {
		logger.Errorf("[分类同步] 处理分类数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类同步] 处理完成: created=%d, updated=%d, deleted=%d, failed=%d",
		result.Apply.Created, result.Apply.Updated, result.Apply.Deleted, result.Apply.Failed)

	// 执行逆向同步
	logger.Infof("[分类同步] 开始执行逆向同步")
	reverseResp, err := s.ReverseSync(ctx, client)
	if err != nil {
		logger.Errorf("❌ [分类同步] 逆向同步失败: %v", err)
		return nil, err
	}
	result.Apply.Deleted = reverseResp.Apply.Deleted
	logger.Infof("[分类同步] 逆向同步完成: deleted=%d", result.Apply.Deleted)

	// 保存同步检查点
	logger.Infof("[分类同步] 准备保存同步检查点")
	checkpointMgr := NewCheckpointManager(s.db)
	checkpointValue := previewResp.Stats.RemoteTotal // 使用远程总数作为检查点值
	if err := checkpointMgr.UpdateCheckpoint(ctx, "CATEGORY",
		fmt.Sprintf("%d", checkpointValue), result, nil); err != nil {
		logger.Errorf("[分类同步] 保存检查点失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类同步] ✓ 检查点已保存")

	logger.Infof("[分类同步] ===== 同步完成 =====")
	return result, nil
}

// fetchLocal 获取本地所有分类
func (s *CategorySyncService) fetchLocal(ctx context.Context) ([]*model.Category, map[string]*model.Category, error) {
	var categories []*model.Category
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&categories).Error; err != nil {
		return nil, nil, err
	}

	// 构建本地分类索引（按上游分类编码）
	index := make(map[string]*model.Category)
	for _, cat := range categories {
		index[cat.SourceCategoryCode] = cat
	}

	return categories, index, nil
}

// compareOne 比较单个分类的本地和远程数据
func (s *CategorySyncService) compareOne(remote *vendors.GameCategoryInfo, localIndex map[string]*model.Category) *vendors.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &vendors.SyncDiff{
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

	// 对比分类名称是否变更
	localName := ""
	if local.SourceNameI18n != nil {
		if nameVal, ok := local.SourceNameI18n["default"]; ok {
			localName = nameVal.(string)
		}
	}

	if localName != remote.Name {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("分类名称变更: %s -> %s", localName, remote.Name)
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// ReverseSync 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *CategorySyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	logger.Info("[分类逆向同步] 开始执行逆向同步")

	// 获取远程分类数据
	// remoteResp, err := client.GetGameCategory(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("❌ [分类逆向同步] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGameCategory()
	if err != nil {
		logger.Errorf("❌ [分类逆向同步] 获取测试数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类逆向同步] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.Data))

	// 构建远程分类编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteCat := range remoteResp.Data {
		remoteIndex[remoteCat.Code] = true
		logger.Debugf("[分类逆向同步] 远程分类编码: %s", remoteCat.Code)
	}
	logger.Infof("[分类逆向同步] 远程分类编码总数: %d", len(remoteIndex))

	// 获取本地分类数据
	var localCategories []*model.Category
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localCategories).Error; err != nil {
		logger.Errorf("❌ [分类逆向同步] 获取本地数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[分类逆向同步] ✓ 获取本地数据成功, 共 %d 条", len(localCategories))

	// 打印本地所有分类编码
	for _, localCat := range localCategories {
		logger.Debugf("[分类逆向同步] 本地分类编码: %s (ID: %d, deleted_at: %v)",
			localCat.SourceCategoryCode, localCat.ID, localCat.DeletedAt)
	}

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: &vendors.SyncPreviewResp{
			Stats: &vendors.SyncStats{
				LocalTotal: int64(len(localCategories)),
			},
		},
		Apply: &vendors.SyncApplyResult{},
	}

	// 在事务中执行软删除
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localCat := range localCategories {
			// 如果本地分类在远程不存在，则软删除
			if !remoteIndex[localCat.SourceCategoryCode] {
				logger.Infof("[分类逆向同步] 检测到本地孤立分类: %s (ID: %d)，准备软删除",
					localCat.SourceCategoryCode, localCat.ID)

				if err := tx.Model(localCat).Update("deleted_at", time.Now()).Error; err != nil {
					logger.Errorf("❌ [分类逆向同步] 软删除失败: %v", err)
					result.Apply.Failed++
					continue
				}
				result.Apply.Deleted++
				logger.Infof("[分类逆向同步] ✓ 已软删除分类: %s (ID: %d)", localCat.SourceCategoryCode, localCat.ID)
			} else {
				logger.Debugf("[分类逆向同步] 本地分类在远程存在，无需删除: %s", localCat.SourceCategoryCode)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	logger.Infof("[分类逆向同步] ✓ 完成，删除: %d", result.Apply.Deleted)
	return result, nil
}

func getTestGameCategory() (*vendors.GetGameCategoryResp, error) {
	var list []*vendors.GameCategoryInfo
	if err := loadTestJSON("category.json", &list); err != nil {
		logger.Errorf("[分类测试数据] 读取 category.json 失败: %v", err)
		return &vendors.GetGameCategoryResp{Data: []*vendors.GameCategoryInfo{}}, err
	}

	return &vendors.GetGameCategoryResp{Data: list}, nil
}
