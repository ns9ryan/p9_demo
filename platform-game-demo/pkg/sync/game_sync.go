package sync

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 游戏同步服务 =====

// GameSyncService 游戏同步服务
type GameSyncService struct {
	db            *gorm.DB
	progressQueue *ProgressQueue
}

// NewGameSyncService 创建游戏同步服务
func NewGameSyncService(db *gorm.DB) *GameSyncService {
	return &GameSyncService{
		db:            db,
		progressQueue: NewProgressQueue(db),
	}
}

// Preview 预检查游戏同步
func (s *GameSyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteGameResponse, error) {
	// 获取远程游戏数据
	remoteResp, err := getRemoteGame(ctx, client, isGetRemoteClient)
	if err != nil {
		log.Printf("[游戏预检查] 获取远程数据失败: %v", err)
		return nil, nil, err
	}

	// 获取本地游戏数据
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		log.Printf("[游戏预检查] 获取本地数据失败: %v", err)
		return nil, nil, err
	}

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

	pageSize := int64(0)
	page := int64(0)
	skipNoop := false
	if req != nil {
		pageSize = req.PageSize
		page = req.Page
		skipNoop = req.IsSkip
	}

	// 创建分页处理器
	handler := &SyncPageHandler{
		RemoteData: remoteItems,
		LocalIndex: localIndexInterface,
		CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
			remoteGame := item.(*vendors.GameInfo)
			localIdx := make(map[string]*ent.Game)
			for k, v := range localIndex {
				localIdx[k] = v.(*ent.Game)
			}
			return s.compareOne(remoteGame, localIdx, nil)
		},
		PageSize: pageSize,
		Page:     page,
		SkipNoop: skipNoop,
	}

	// 处理分页逻辑
	result := handler.Handle(int64(len(remoteResp.Data)), int64(len(local)))

	return result, remoteResp, nil
}

// validateForeignKeys 验证游戏的外键（分类、厂商、渠道是否存在）
// 返回对应的本地ID，如果外键不存在则返回错误
func (s *GameSyncService) validateRequiredTables(ctx context.Context) error {
	var count int64

	if err := s.db.WithContext(ctx).Model(&ent.GameCategory{}).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		logger.Error("[游戏同步] 本地分类表为空，无法执行游戏同步")
		return fmt.Errorf("本地分类表为空，无法执行游戏同步")
	}

	if err := s.db.WithContext(ctx).Model(&ent.GameProvider{}).
		Where("deleted_at IS NULL").
		Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		logger.Error("[游戏同步] 本地厂商表为空，无法执行游戏同步")
		return fmt.Errorf("本地厂商表为空，无法执行游戏同步")
	}

	if err := s.db.WithContext(ctx).Model(&ent.GameChannel{}).
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
	var category ent.GameCategory
	if err := tx.Where("source_id = ?", remoteGame.CatId).
		Where("deleted_at IS NULL").
		First(&category).Error; err != nil {
		logger.Errorf("[游戏同步-外键验证] 游戏 %s 分类不存在: source_id=%d, 错误: %v", remoteGame.Code, remoteGame.CatId, err)
		return 0, 0, 0, fmt.Errorf("分类不存在: source_id=%d", remoteGame.CatId)
	}
	categoryID = category.SourceId

	// 验证厂商是否存在
	if remoteGame.VenId == 0 {
		logger.Errorf("[游戏同步-外键验证] 游戏 %s 厂商ID为0（无效）", remoteGame.Code)
		return 0, 0, 0, fmt.Errorf("厂商ID为0")
	}
	var provider ent.GameProvider
	if err := tx.Where("source_id = ?", remoteGame.VenId).
		Where("deleted_at IS NULL").
		First(&provider).Error; err != nil {
		logger.Errorf("[游戏同步-外键验证] 游戏 %s 厂商不存在: source_id=%d, 错误: %v", remoteGame.Code, remoteGame.VenId, err)
		return 0, 0, 0, fmt.Errorf("厂商不存在: source_id=%d", remoteGame.VenId)
	}
	providerID = provider.SourceId

	// 验证渠道是否存在（渠道可以为0，表示厂商直连）
	if remoteGame.ChanId == 0 {
		logger.Debugf("[游戏同步-外键验证] 游戏 %s 渠道ID为0（厂商直连）", remoteGame.Code)
		channelID = 0
	} else {
		var channel ent.GameChannel
		if err := tx.Where("source_id = ?", remoteGame.ChanId).
			Where("deleted_at IS NULL").
			First(&channel).Error; err != nil {
			logger.Errorf("[游戏同步-外键验证] 游戏 %s 渠道不存在: source_id=%d, 错误: %v", remoteGame.Code, remoteGame.ChanId, err)
			return 0, 0, 0, fmt.Errorf("渠道不存在: source_id=%d", remoteGame.ChanId)
		}
		channelID = channel.SourceId
	}

	return categoryID, providerID, channelID, nil
}

// 处理游戏数据的创建和更新（从 Run 提取的业务逻辑）
func (s *GameSyncService) processGameDataWithProgress(ctx context.Context, tx *gorm.DB, remoteData []*vendors.GameInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
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
		var localGame ent.Game
		exists := tx.Where("source_id = ?", remoteGame.Id).
			Where("deleted_at IS NULL").
			First(&localGame).Error == nil

		logger.Infof("[游戏同步] 本地游戏存在: %v", exists)
		gameCode := remoteGame.Code + "_" + fmt.Sprintf("%d", i)
		nameI18n := map[string]interface{}{"default": remoteGame.Name}
		sourceNameI18n := map[string]interface{}{"default": remoteGame.Name}
		if !exists {
			// 新增游戏前验证外键
			categoryID, providerID, channelID, err := s.validateForeignKeys(ctx, tx, remoteGame, nil)
			if err != nil {
				logger.Errorf("[游戏同步] 游戏 %s 外键验证失败: %v", remoteGame.Code, err)
				applyResult.Failed++
				continue
			}
			logger.Infof("[游戏新增] 即将创建: Code=%s, ChannelID=%d, CategoryID=%d, ProviderID=%d",
				gameCode, channelID, categoryID, providerID)

			newGame := ent.Game{
				SourceId:          sql.NullInt64{Int64: remoteGame.Id, Valid: true},
				CategoryId:        categoryID,
				ProviderId:        providerID,
				ChannelId:         sql.NullInt64{Int64: channelID, Valid: true},
				GameCode:          gameCode,
				SourceGameCode:    remoteGame.Code,
				NameI18n:          JSONToString(nameI18n),
				SourceNameI18n:    JSONToString(sourceNameI18n),
				ImageUrl:          sql.NullString{String: remoteGame.Image, Valid: true},
				SourceImageUrl:    sql.NullString{String: remoteGame.Image, Valid: true},
				ProviderKey:       sql.NullString{String: remoteGame.VenKey, Valid: true},
				SourceProviderKey: sql.NullString{String: remoteGame.VenKey, Valid: true},
				SortNo:            remoteGame.Id,
				SourceSortNo:      sql.NullInt64{Int64: 0, Valid: true},
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
			updateData := map[string]interface{}{
				"source_id":           remoteGame.Id,
				"category_id":         categoryID,
				"provider_id":         providerID,
				"channel_id":          channelID,
				"game_code":           gameCode,
				"source_game_code":    remoteGame.Code,
				"source_name_i18n":    JSONToString(sourceNameI18n),
				"image_url":           sql.NullString{String: remoteGame.Image, Valid: true},
				"source_image_url":    sql.NullString{String: remoteGame.Image, Valid: true},
				"provider_key":        sql.NullString{String: remoteGame.VenKey, Valid: true},
				"source_provider_key": sql.NullString{String: remoteGame.VenKey, Valid: true},
				"source_sort_no":      sql.NullInt64{Int64: 0, Valid: true},
				"source_status":       1,
				"updated_at":          time.Now(),
			}
			if len(syncCols) > 0 {
				updateData = GetUpdatedData(map[string]interface{}{
					"source_id":           remoteGame.Id,
					"category_id":         categoryID,
					"provider_id":         providerID,
					"channel_id":          channelID,
					"game_code":           gameCode,
					"source_game_code":    remoteGame.Code,
					"name_i18n":           JSONToString(nameI18n),
					"source_name_i18n":    JSONToString(sourceNameI18n),
					"image_url":           sql.NullString{String: remoteGame.Image, Valid: true},
					"source_image_url":    sql.NullString{String: remoteGame.Image, Valid: true},
					"provider_key":        sql.NullString{String: remoteGame.VenKey, Valid: true},
					"source_provider_key": sql.NullString{String: remoteGame.VenKey, Valid: true},
					"sort_no":             sql.NullInt64{Int64: int64(remoteGame.Id), Valid: true},
					"source_sort_no":      sql.NullInt64{Int64: 0, Valid: true},
					"status":              1,
					"source_status":       1,
				}, syncCols)
				updateData["updated_at"] = time.Now()
			}

			if err := tx.Model(&localGame).
				Updates(updateData).Error; err != nil {
				logger.Errorf("[游戏更新] 失败: %v", err)
				return err
			}
			applyResult.Updated++
		}
	}

	return nil
}

// Run 执行游戏同步
func (s *GameSyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	logger.Infof("[游戏同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[游戏同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		logger.Errorf("[游戏同步] 预检查失败: %v", err)
		return err
	}
	// 创建并启动进度队列
	s.progressQueue.Start(ctx)
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal {
		defer s.progressQueue.Stop()
		logger.Infof("[游戏同步] 无需更新，直接返回")
		logger.Infof("[游戏同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
			previewResp.Stats.UpdateTotal, previewResp.Stats.CreateTotal, previewResp.Stats.DeleteTotal, previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
		// 更新checkpointID进度为100
		s.progressQueue.Send(&ProgressMessage{
			TableName:      "game",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
		})
		return nil
	}
	logger.Infof("[游戏同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 先校验本地依赖表是否已准备好
	logger.Infof("[游戏同步] 校验依赖表")
	if err := s.validateRequiredTables(ctx); err != nil {
		logger.Errorf("[游戏同步] 依赖表校验失败: %v", err)
		return err
	}
	logger.Infof("[游戏同步] ✓ 检查点已创建: checkpointID=%d", checkpointID)

	go func() {
		defer s.progressQueue.Stop()
		// 在事务中执行数据库操作
		logger.Infof("[游戏同步] 开始处理游戏数据")
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
			logger.Infof("[游戏同步] 处理批次: start=%d, end=%d", i, end)
			// 模拟等待
			logger.Debugf("[进度队列] 模拟处理延迟: 表=game, checkpointID=%d", checkpointID)
			time.Sleep(1000 * time.Millisecond)
			logger.Debugf("[进度队列] 开始处理消息: 表=game, checkpointID=%d", checkpointID)
			err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				return s.processGameDataWithProgress(ctx, tx, batch, apply, syncCols, checkpointID)
			})
			if err != nil {
				logger.Errorf("[游戏同步] 处理游戏数据失败: %v", err)
				return
			}
			progress := CalculateProgress(int64(i+1), int64(totalCount))
			msg := &ProgressMessage{
				TableName:      "game",
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
			logger.Infof("[游戏同步] 发送进度消息: 表=game, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
			s.progressQueue.Send(msg)
		}
		if err != nil {
			logger.Errorf("[游戏同步] 处理游戏数据失败: %v", err)
			return
		}

		// 执行逆向同步
		logger.Infof("[游戏同步] 开始执行逆向同步")
		err := s.ReverseSync(ctx, client, remoteResp, apply)
		if err != nil {
			logger.Errorf("❌ [游戏同步] 逆向同步失败: %v", err)
			return
		}
		logger.Infof("[游戏同步] 逆向同步完成: deleted=%d", apply.Deleted)

		// 最后一次进度更新（100%）
		logger.Infof("[游戏同步] 发送最终进度消息（100%%）")
		finalMsg := &ProgressMessage{
			TableName:      "game",
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
		logger.Infof("[游戏同步] ===== 同步完成 =====")
	}()
	// 构建执行结果
	return nil
}

// fetchLocal 获取本地所有游戏
func (s *GameSyncService) fetchLocal(ctx context.Context) ([]*ent.Game, map[string]*ent.Game, error) {
	var games []*ent.Game
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&games).Error; err != nil {
		return nil, nil, err
	}

	// 构建本地游戏索引（使用SourceGameCode作为键）
	index := make(map[string]*ent.Game)
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
func (s *GameSyncService) compareOne(remote *vendors.GameInfo, localIndex map[string]*ent.Game, depIndex *DependencyIndex) *platform_game.SyncDiff {
	local, exists := localIndex[remote.Code]

	// 初始化差异记录
	diff := &platform_game.SyncDiff{
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

	if CompareName(local.SourceNameI18n, remote.Name) {
		diff.Action = "update"
		diff.Reason = fmt.Sprintf("游戏名称变更: %s -> %s", local.SourceNameI18n, remote.Name)
		diff.ConflictType = "name_change"
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *GameSyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteGameResponse, apply *Apply) error {
	// 构建远程游戏编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteGame := range remoteResp.Data {
		remoteIndex[remoteGame.Code] = true
	}

	// 获取本地游戏数据
	var localGames []*ent.Game
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localGames).Error; err != nil {
		log.Printf("[游戏逆向同步] ✗ 获取本地数据失败: %v\n", err)
		return err
	}
	log.Printf("[游戏逆向同步] ✓ 获取本地数据成功, 共 %d 条\n", len(localGames))

	// 在事务中执行软删除
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localGame := range localGames {
			// 如果本地游戏在远程不存在，则软删除
			if !remoteIndex[localGame.SourceGameCode] {
				// // 检查唯一性：是否存在其他 deleted_at IS NULL 的记录有相同的 source_game_code
				// var duplicateCount int64
				// if err := tx.Model(&ent.Game{}).
				// 	Where("source_id = ?", localGame.SourceId).
				// 	Where("deleted_at IS NULL").
				// 	Count(&duplicateCount).Error; err != nil {
				// 	log.Printf("[游戏逆向同步] 唯一性检查失败: %v", err)
				// 	apply.Failed++
				// 	continue
				// }

				// if duplicateCount > 1 {
				// 	logger.Errorf("[游戏逆向同步] 存在 %d 条记录有相同的 source_id=%d, 跳过本条删除以保护数据", duplicateCount, localGame.SourceId)
				// 	apply.Failed++
				// 	continue
				// }

				if err := tx.Model(localGame).Update("deleted_at", time.Now()).Error; err != nil {
					log.Printf("[游戏逆向同步] 软删除失败: %v", err)
					apply.Failed++
					continue
				}
				apply.Deleted++
				log.Printf("[游戏逆向同步] ✓ 已软删除游戏: %s\n", localGame.SourceGameCode)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("[游戏逆向同步] ✓ 完成，删除: %d\n", apply.Deleted)
	return nil
}

func getRemoteGame(ctx context.Context, client vendors.VendorGameServiceClient, isGetRemoteClient bool) (*RemoteGameResponse, error) {
	if isGetRemoteClient {
		logger.Infof("[同步数据源] 使用远程客户端获取游戏信息")
		remoteResp, err := client.GetGame(ctx, &vendors.GetGameRequest{})
		if err != nil {
			logger.Errorf("❌ [同步数据源] 获取远程数据失败: %v", err)
			return nil, err
		}
		logger.Infof("[同步数据源] 获取到远程数据: count=%d", len(remoteResp.GameList))
		return &RemoteGameResponse{Data: remoteResp.GameList}, nil
	} else {
		logger.Infof("[同步数据源] 使用本地 JSON 数据获取游戏信息")
		var list []*vendors.GameInfo
		if err := loadLocalJSON("game.json", &list); err != nil {
			logger.Errorf("[同步数据源] 读取 game.json 失败: %v", err)
			return nil, err
		}
		logger.Infof("[同步数据源] ✓ 读取本地 JSON 数据成功, 共 %d 条", len(list))
		return &RemoteGameResponse{Data: list}, nil
	}
}

type RemoteGameResponse struct {
	Data []*vendors.GameInfo `json:"data"`
}
