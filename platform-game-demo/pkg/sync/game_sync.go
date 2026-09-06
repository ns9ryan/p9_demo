package sync

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 游戏同步服务 =====

// GameSyncService 游戏同步服务
type GameSyncService struct {
	db *gorm.DB
}

// NewGameSyncService 创建游戏同步服务
func NewGameSyncService(db *gorm.DB) *GameSyncService {
	return &GameSyncService{db: db}
}

// Preview 预检查游戏同步
func (s *GameSyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncPreviewResp, error) {
	// 获取远程游戏数据
	// remoteResp, err := client.GetGame(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	log.Printf("[游戏预检查] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGame()
	if err != nil {
		log.Printf("[游戏预检查] 获取测试数据失败: %v", err)
		return nil, err
	}
	remote := remoteResp.Data

	// 获取本地游戏数据
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

	// 对每条远程游戏进行比较
	for _, remoteItem := range remote {
		diff := s.compareOne(remoteItem, localIndex, nil)
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

// validateForeignKeys 验证游戏的外键（分类、厂商、渠道是否存在）
// 返回对应的本地ID，如果外键不存在则返回错误
func (s *GameSyncService) validateRequiredTables(ctx context.Context) error {
	var count int64

	if err := s.db.WithContext(ctx).Model(&model.Category{}).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		logger.Error("[游戏同步] 本地分类表为空，无法执行游戏同步")
		return fmt.Errorf("本地分类表为空，无法执行游戏同步")
	}

	if err := s.db.WithContext(ctx).Model(&model.Provider{}).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		logger.Error("[游戏同步] 本地厂商表为空，无法执行游戏同步")
		return fmt.Errorf("本地厂商表为空，无法执行游戏同步")
	}

	if err := s.db.WithContext(ctx).Model(&model.Channel{}).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		logger.Error("[游戏同步] 本地渠道表为空，无法执行游戏同步")
		return fmt.Errorf("本地渠道表为空，无法执行游戏同步")
	}

	return nil
}

func (s *GameSyncService) validateForeignKeys(ctx context.Context, tx *gorm.DB, remoteGame *vendors.GameInfo, depIndex *DependencyIndex) (categoryID, providerID, channelID int64, err error) {
	// 验证分类是否存在
	if remoteGame.CatId == 0 {
		logger.Errorf("[游戏同步-外键验证] 游戏 %s 分类ID为0（无效）", remoteGame.Code)
		return 0, 0, 0, fmt.Errorf("分类ID为0")
	}
	var category model.Category
	if err := tx.Where("source_id = ?", remoteGame.CatId).
		Where("deleted_at IS NULL").
		First(&category).Error; err != nil {
		logger.Errorf("[游戏同步-外键验证] 游戏 %s 分类不存在: source_id=%d, 错误: %v", remoteGame.Code, remoteGame.CatId, err)
		return 0, 0, 0, fmt.Errorf("分类不存在: source_id=%d", remoteGame.CatId)
	}
	categoryID = category.ID

	// 验证厂商是否存在
	if remoteGame.VenId == 0 {
		logger.Errorf("[游戏同步-外键验证] 游戏 %s 厂商ID为0（无效）", remoteGame.Code)
		return 0, 0, 0, fmt.Errorf("厂商ID为0")
	}
	var provider model.Provider
	if err := tx.Where("source_id = ?", remoteGame.VenId).
		Where("deleted_at IS NULL").
		First(&provider).Error; err != nil {
		logger.Errorf("[游戏同步-外键验证] 游戏 %s 厂商不存在: source_id=%d, 错误: %v", remoteGame.Code, remoteGame.VenId, err)
		return 0, 0, 0, fmt.Errorf("厂商不存在: source_id=%d", remoteGame.VenId)
	}
	providerID = provider.ID

	// 验证渠道是否存在（渠道可以为0，表示厂商直连）
	if remoteGame.ChanId == 0 {
		logger.Debugf("[游戏同步-外键验证] 游戏 %s 渠道ID为0（厂商直连）", remoteGame.Code)
		channelID = 0
	} else {
		var channel model.Channel
		if err := tx.Where("source_id = ?", remoteGame.ChanId).
			Where("deleted_at IS NULL").
			First(&channel).Error; err != nil {
			logger.Errorf("[游戏同步-外键验证] 游戏 %s 渠道不存在: source_id=%d, 错误: %v", remoteGame.Code, remoteGame.ChanId, err)
			return 0, 0, 0, fmt.Errorf("渠道不存在: source_id=%d", remoteGame.ChanId)
		}
		channelID = channel.ID
	}

	return categoryID, providerID, channelID, nil
}

// processGameData 处理游戏数据的创建和更新（从 Run 提取的业务逻辑）
func (s *GameSyncService) processGameData(ctx context.Context, tx *gorm.DB, remoteData []*vendors.GameInfo, applyResult *vendors.SyncApplyResult) error {
	i := 0
	counterMap := make(map[int64]int) // 用于记录每个 source_id 的计数，确保唯一性
	for _, remoteGame := range remoteData {
		if counterMap[remoteGame.Id] > 0 {
			logger.Errorf("[游戏同步] 检测到重复的远程游戏ID: %d, 计数器: %d, 跳过处理", remoteGame.Id, counterMap[remoteGame.Id])
			applyResult.Failed++
			continue
		}
		// 更新计数器
		counterMap[remoteGame.Id]++
		i++

		// 添加日志记录远程数据
		logger.Infof("[游戏同步] 处理游戏: Code=%s, Id=%d, CatId=%d, VenId=%d, ChanId=%d",
			remoteGame.Code, remoteGame.Id, remoteGame.CatId, remoteGame.VenId, remoteGame.ChanId)

		// 检查本地是否存在该游戏（使用上游游戏ID）
		var localGame model.Game
		exists := tx.Where("source_id = ?", remoteGame.Id).
			Where("deleted_at IS NULL").
			First(&localGame).Error == nil

		logger.Infof("[游戏同步] 本地游戏存在: %v", exists)

		if !exists {
			// 新增游戏前验证外键
			categoryID, providerID, channelID, err := s.validateForeignKeys(ctx, tx, remoteGame, nil)
			if err != nil {
				logger.Errorf("[游戏同步] 游戏 %s 外键验证失败: %v", remoteGame.Code, err)
				applyResult.Failed++
				continue
			}

			// 新增游戏 - 生成 P9 内部的游戏编码
			gameCode := remoteGame.Code + "_" + fmt.Sprintf("%d", i)
			nameI18n := model.JSONMap{"default": remoteGame.Name}
			sourceNameI18n := model.JSONMap{"default": remoteGame.Name}

			logger.Infof("[游戏新增] 即将创建: Code=%s, ChannelID=%d, CategoryID=%d, ProviderID=%d",
				gameCode, channelID, categoryID, providerID)

			newGame := model.Game{
				SourceID:          remoteGame.Id,
				CategoryID:        categoryID,
				ProviderID:        providerID,
				ChannelID:         channelID,
				GameCode:          gameCode,
				SourceGameCode:    remoteGame.Code,
				NameI18n:          nameI18n,
				SourceNameI18n:    sourceNameI18n,
				ImageURL:          remoteGame.Image,
				SourceImageURL:    remoteGame.Image,
				ProviderKey:       remoteGame.VenKey,
				SourceProviderKey: remoteGame.VenKey,
				SourceStatus:      1,
				Status:            1,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			}
			if err := tx.Create(&newGame).Error; err != nil {
				logger.Errorf("[游戏新增] 失败: %v", err)
				return err
			}
			applyResult.Created++
		} else {
			// 更新游戏前验证外键
			categoryID, providerID, channelID, err := s.validateForeignKeys(ctx, tx, remoteGame, nil)
			if err != nil {
				logger.Errorf("[游戏同步] 游戏 %s 外键验证失败: %v", remoteGame.Code, err)
				applyResult.Failed++
				continue
			}

			gameCode := remoteGame.Code + "_" + fmt.Sprintf("%d", i)
			sourceNameI18n := model.JSONMap{"default": remoteGame.Name}

			logger.Infof("[游戏更新] 即将更新: Code=%s, ChannelID=%d, CategoryID=%d, ProviderID=%d",
				gameCode, channelID, categoryID, providerID)

			if err := tx.Model(&localGame).
				Updates(map[string]interface{}{
					"source_id":           remoteGame.Id,
					"category_id":         categoryID,
					"provider_id":         providerID,
					"channel_id":          channelID,
					"game_code":           gameCode,
					"source_game_code":    remoteGame.Code,
					"source_name_i18n":    sourceNameI18n,
					"image_url":           remoteGame.Image,
					"source_image_url":    remoteGame.Image,
					"provider_key":        remoteGame.VenKey,
					"source_provider_key": remoteGame.VenKey,
					"source_status":       1,
					"updated_at":          time.Now(),
				}).Error; err != nil {
				logger.Errorf("[游戏更新] 失败: %v", err)
				return err
			}
			applyResult.Updated++
		}
	}
	return nil
}

// Run 执行游戏同步
func (s *GameSyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	logger.Infof("[游戏同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[游戏同步] 执行预检查")
	previewResp, err := s.Preview(ctx, client)
	if err != nil {
		logger.Errorf("[游戏同步] 预检查失败: %v", err)
		return nil, err
	}
	logger.Infof("[游戏同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: previewResp,
		Apply:   &vendors.SyncApplyResult{},
	}

	// 获取远程游戏数据
	logger.Infof("[游戏同步] 获取远程游戏数据")
	// remoteResp, err := client.GetGame(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	logger.Errorf("[游戏同步] 获取远程数据失败: %v", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGame()
	if err != nil {
		logger.Errorf("[游戏同步] 获取测试数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[游戏同步] 获取到远程数据: count=%d", len(remoteResp.Data))

	// 先校验本地依赖表是否已准备好
	logger.Infof("[游戏同步] 校验依赖表")
	if err := s.validateRequiredTables(ctx); err != nil {
		logger.Errorf("[游戏同步] 依赖表校验失败: %v", err)
		return nil, err
	}
	logger.Infof("[游戏同步] 依赖表校验完成")

	// 在事务中执行数据库操作
	logger.Infof("[游戏同步] 开始处理游戏数据")
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.processGameData(ctx, tx, remoteResp.Data, result.Apply)
	})
	if err != nil {
		logger.Errorf("[游戏同步] 处理游戏数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[游戏同步] 处理完成: created=%d, updated=%d, deleted=%d, failed=%d",
		result.Apply.Created, result.Apply.Updated, result.Apply.Deleted, result.Apply.Failed)

	// 执行逆向同步
	logger.Infof("[游戏同步] 开始执行逆向同步")
	reverseResp, err := s.ReverseSync(ctx, client)
	if err != nil {
		logger.Errorf("❌ [游戏同步] 逆向同步失败: %v", err)
		return nil, err
	}
	result.Apply.Deleted = reverseResp.Apply.Deleted
	logger.Infof("[游戏同步] 逆向同步完成: deleted=%d", result.Apply.Deleted)

	// 保存同步检查点
	logger.Infof("[游戏同步] 准备保存同步检查点")
	checkpointMgr := NewCheckpointManager(s.db)
	checkpointValue := previewResp.Stats.RemoteTotal // 使用远程总数作为检查点值
	if err := checkpointMgr.UpdateCheckpoint(ctx, "GAME",
		fmt.Sprintf("%d", checkpointValue), result, nil); err != nil {
		logger.Errorf("[游戏同步] 保存检查点失败: %v", err)
		return nil, err
	}
	logger.Infof("[游戏同步] ✓ 检查点已保存")

	logger.Infof("[游戏同步] ===== 同步完成 =====")
	return result, nil
}

// fetchLocal 获取本地所有游戏
func (s *GameSyncService) fetchLocal(ctx context.Context) ([]*model.Game, map[string]*model.Game, error) {
	var games []*model.Game
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&games).Error; err != nil {
		return nil, nil, err
	}

	// 构建本地游戏索引（使用SourceGameCode作为键）
	index := make(map[string]*model.Game)
	for _, game := range games {
		index[game.SourceGameCode] = game
	}

	return games, index, nil
}

// DependencyIndex 依赖索引结构
type DependencyIndex struct {
	// 远程ID -> 本地ID的映射
	Categories map[int64]int64
	Providers  map[int64]int64
	Channels   map[int64]int64
}

// compareOne 比较单个游戏的本地和远程数据
func (s *GameSyncService) compareOne(remote *vendors.GameInfo, localIndex map[string]*model.Game, depIndex *DependencyIndex) *vendors.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &vendors.SyncDiff{
		ObjectType: "game",
		ObjectId:   remote.Id,
		ObjectCode: remote.Code,
		RemoteId:   remote.Id,
		RemoteCode: remote.Code,
	}

	// 如果本地不存在该游戏
	if !exists {
		diff.Action = "create"
		diff.Reason = "本地不存在该游戏"
		return diff
	}

	// 对比游戏名称是否变更
	var localSourceName string
	if sourceNameI18n, ok := local.SourceNameI18n["default"]; ok {
		localSourceName = sourceNameI18n.(string)
	}
	if localSourceName != remote.Name {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("游戏名称变更: %s -> %s", localSourceName, remote.Name)
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// ReverseSync 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *GameSyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	log.Println("[游戏逆向同步] 开始执行逆向同步...")

	// 获取远程游戏数据
	// remoteResp, err := client.GetGame(ctx, &emptypb.Empty{})
	// if err != nil {
	// 	log.Printf("[游戏逆向同步] ✗ 获取远程数据失败: %v\n", err)
	// 	return nil, err
	// }
	remoteResp, err := getTestGame()
	if err != nil {
		log.Printf("[游戏逆向同步] ✗ 获取远程数据失败: %v\n", err)
		return nil, err
	}
	log.Printf("[游戏逆向同步] ✓ 获取远程数据成功, 共 %d 条\n", len(remoteResp.Data))

	// 构建远程游戏编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteGame := range remoteResp.Data {
		remoteIndex[remoteGame.Code] = true
	}

	// 获取本地游戏数据
	var localGames []*model.Game
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localGames).Error; err != nil {
		log.Printf("[游戏逆向同步] ✗ 获取本地数据失败: %v\n", err)
		return nil, err
	}
	log.Printf("[游戏逆向同步] ✓ 获取本地数据成功, 共 %d 条\n", len(localGames))

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: &vendors.SyncPreviewResp{
			Stats: &vendors.SyncStats{
				LocalTotal: int64(len(localGames)),
			},
		},
		Apply: &vendors.SyncApplyResult{},
	}

	// 在事务中执行软删除
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localGame := range localGames {
			// 如果本地游戏在远程不存在，则软删除
			if !remoteIndex[localGame.SourceGameCode] {
				// 检查唯一性：是否存在其他 deleted_at IS NULL 的记录有相同的 source_game_code
				var duplicateCount int64
				if err := tx.Model(&model.Game{}).
					Where("source_id = ?", localGame.SourceID).
					Where("deleted_at IS NULL").
					Count(&duplicateCount).Error; err != nil {
					log.Printf("[游戏逆向同步] 唯一性检查失败: %v", err)
					result.Apply.Failed++
					continue
				}

				if duplicateCount > 1 {
					logger.Errorf("[游戏逆向同步] 存在 %d 条记录有相同的 source_id=%d, 跳过本条删除以保护数据", duplicateCount, localGame.SourceID)
					result.Apply.Failed++
					continue
				}

				if err := tx.Model(localGame).Update("deleted_at", time.Now()).Error; err != nil {
					log.Printf("[游戏逆向同步] 软删除失败: %v", err)
					result.Apply.Failed++
					continue
				}
				result.Apply.Deleted++
				log.Printf("[游戏逆向同步] ✓ 已软删除游戏: %s\n", localGame.SourceGameCode)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("[游戏逆向同步] ✓ 完成，删除: %d\n", result.Apply.Deleted)
	return result, nil
}

func getTestGame() (*vendors.GetGameResp, error) {
	var list []*vendors.GameInfo
	if err := loadTestJSON("game.json", &list); err != nil {
		log.Printf("[游戏测试数据] 读取 game.json 失败: %v", err)
		return &vendors.GetGameResp{Data: []*vendors.GameInfo{}}, err
	}

	return &vendors.GetGameResp{Data: list}, nil
}
