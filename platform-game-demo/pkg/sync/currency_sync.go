package sync

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 游戏货币同步服务 =====

// CurrencySyncService 游戏货币同步服务
type CurrencySyncService struct {
	db            *gorm.DB
	progressQueue *ProgressQueue
}

// NewCurrencySyncService 创建游戏货币同步服务
func NewCurrencySyncService(db *gorm.DB) *CurrencySyncService {
	return &CurrencySyncService{
		db:            db,
		progressQueue: NewProgressQueue(db),
	}
}

func ensureSysCurrencyID(tx *gorm.DB, code string) (int64, error) {
	code = strings.ToUpper(strings.TrimSpace(code))
	if code == "" {
		return 0, fmt.Errorf("currency_code 为空")
	}

	var currencyID int64
	if err := tx.Table("sys_currency").Select("id").Where("currency_code = ?", code).Scan(&currencyID).Error; err != nil {
		return 0, err
	}
	if currencyID > 0 {
		return currencyID, nil
	}

	now := time.Now()
	data := map[string]interface{}{
		"currency_code": code,
		"name_i18n":     map[string]string{"default": code},
		"currency_type": 1,
		"amount_factor": int64(1),
		"status":        int32(1),
		"sort_no":       0,
		"created_at":    now,
		"updated_at":    now,
	}
	if err := tx.Table("sys_currency").Create(data).Error; err != nil {
		return 0, err
	}

	if err := tx.Table("sys_currency").Select("id").Where("currency_code = ?", code).Scan(&currencyID).Error; err != nil {
		return 0, err
	}
	if currencyID == 0 {
		return 0, fmt.Errorf("currency_code=%s 插入后未查询到 id", code)
	}
	return currencyID, nil
}

// Preview 预检查游戏货币同步
func (s *CurrencySyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteGameCurrencyResponse, error) {
	// 获取远程游戏货币数据
	remoteResp, err := getRemoteGameCurrency(ctx, client, s.db, isGetRemoteClient)
	if err != nil {
		logger.Errorf("[货币预检查] 获取远程数据失败: %v", err)
		return nil, nil, err
	}
	remote := remoteResp.Data

	// 获取本地游戏货币数据
	local, err := s.fetchLocal(ctx)
	if err != nil {
		return nil, nil, err
	}

	// 转换为接口切片
	remoteItems := make([]interface{}, len(remote))
	for i, v := range remote {
		remoteItems[i] = v
	}

	// 当 req 为 nil 时，返回全量结果不做分页
	if req == nil {
		handler := &SyncPageHandler{
			RemoteData: remoteItems,
			LocalIndex: make(map[string]interface{}),
			CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
				remoteCurrency := item.(*vendors.GameCurrencyInfo)
				return s.compareOne(remoteCurrency, local)
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
		LocalIndex: make(map[string]interface{}),
		CompareFunc: func(item interface{}, localIndex map[string]interface{}) *platform_game.SyncDiff {
			remoteCurrency := item.(*vendors.GameCurrencyInfo)
			return s.compareOne(remoteCurrency, local)
		},
		PageSize: req.PageSize,
		Page:     req.Page,
		SkipNoop: req.IsSkip,
	}

	// 处理分页逻辑
	result := handler.Handle(int64(len(remote)), int64(len(local)))
	return result, remoteResp, nil
}

// processCurrencyDataWithProgress 处理游戏货币数据的创建和更新（从 Run 提取的业务逻辑）
func (s *CurrencySyncService) processCurrencyDataWithProgress(ctx context.Context, tx *gorm.DB, remoteData []*vendors.GameCurrencyInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 currency key 的计数，确保唯一性
	for _, remoteCurrency := range remoteData {
		currencyKey := fmt.Sprintf("%d_%d", remoteCurrency.GameId, remoteCurrency.CurrencyId)
		if counterMap[currencyKey] > 0 {
			logger.Errorf("[货币同步] 检测到重复的远程货币编码: GameID=%d, CurrencyID=%d, 计数器: %d, 跳过处理", remoteCurrency.GameId, remoteCurrency.CurrencyId, counterMap[currencyKey])
			applyResult.Failed++
			continue
		}
		// 更新计数器
		counterMap[currencyKey]++
		i++
		// 检查本地是否存在该游戏货币
		var localCurrency ent.GameCurrency
		exists := tx.Where("game_id = ?", remoteCurrency.GameId).Where("currency_id = ?", remoteCurrency.CurrencyId).
			Where("deleted_at IS NULL").First(&localCurrency).Error == nil

		if !exists {
			// 新增游戏货币
			newCurrency := ent.GameCurrency{
				GameId:       remoteCurrency.GameId,
				CurrencyId:   remoteCurrency.CurrencyId,
				SourceStatus: 1,
				Status:       1,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			}
			if err := tx.Create(&newCurrency).Error; err != nil {
				logger.Errorf("[货币新增] 失败: %v", err)
				return err
			}
			applyResult.Created++
		} else {
			// 更新游戏货币
			updateData := map[string]interface{}{
				"source_status": 1,
				"updated_at":    time.Now(),
			}
			if len(syncCols) > 0 {
				updateData = GetUpdatedData(map[string]interface{}{
					"status":        1,
					"source_status": 1,
				}, syncCols)
				updateData["updated_at"] = time.Now()
			}
			if err := tx.Model(&localCurrency).
				Updates(updateData).Error; err != nil {
				logger.Errorf("[货币更新] 失败: %v", err)
				return err
			}
			applyResult.Updated++
		}
	}

	return nil
}

// Run 执行游戏货币同步
func (s *CurrencySyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	logger.Infof("[货币同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[货币同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		logger.Errorf("[货币同步] 预检查失败: %v", err)
		return err
	}
	// 创建并启动进度队列
	s.progressQueue.Start(ctx)
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal {
		defer s.progressQueue.Stop()
		logger.Infof("[货币同步] 无需更新，直接返回")
		logger.Infof("[货币同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
			previewResp.Stats.UpdateTotal, previewResp.Stats.CreateTotal, previewResp.Stats.DeleteTotal, previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
		// 更新checkpointID进度为100
		s.progressQueue.Send(&ProgressMessage{
			TableName:      "currency",
			ProcessedCount: previewResp.Stats.RemoteTotal,
			RemoteTotal:    previewResp.Stats.RemoteTotal,
			LocalTotal:     previewResp.Stats.LocalTotal,
			Progress:       100,
			CheckpointID:   checkpointID,
		})
		return nil
	}
	logger.Infof("[货币同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)
	// 创建初始检查点记录（此时progress=0）
	logger.Infof("[货币同步] 创建初始检查点")
	logger.Infof("[货币同步] ✓ 检查点已创建: checkpointID=%d", checkpointID)

	go func() {
		defer s.progressQueue.Stop()
		// 在事务中执行数据库操作
		logger.Infof("[货币同步] 开始处理货币数据")
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
			logger.Infof("[货币同步] 处理批次: start=%d, end=%d", i, end)
			// 模拟等待
			logger.Debugf("[进度队列] 模拟处理延迟: 表=currency, checkpointID=%d", checkpointID)
			time.Sleep(1000 * time.Millisecond)
			logger.Debugf("[进度队列] 开始处理消息: 表=currency, checkpointID=%d", checkpointID)
			err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
				return s.processCurrencyDataWithProgress(ctx, tx, batch, apply, syncCols, checkpointID)
			})
			if err != nil {
				logger.Errorf("[货币同步] 处理货币数据失败: %v", err)
				return
			}
			progress := CalculateProgress(int64(i+1), int64(totalCount))
			msg := &ProgressMessage{
				TableName:      "currency",
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
			logger.Infof("[货币同步] 发送进度消息: 表=currency, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
			s.progressQueue.Send(msg)
		}
		if err != nil {
			logger.Errorf("[货币同步] 处理货币数据失败: %v", err)
			return
		}

		// 执行逆向同步
		logger.Infof("[货币同步] 开始执行逆向同步")
		err = s.ReverseSync(ctx, client, remoteResp, apply)
		if err != nil {
			logger.Errorf("❌ [货币同步] 逆向同步失败: %v", err)
			return
		}
		logger.Infof("[货币同步] 逆向同步完成: deleted=%d", apply.Deleted)

		// 最后一次进度更新（100%）
		logger.Infof("[货币同步] 发送最终进度消息（100%%）")
		finalMsg := &ProgressMessage{
			TableName:      "currency",
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

		logger.Infof("[货币同步] ===== 同步完成 =====")
	}()
	// 构建执行结果
	return nil
}

// 获取本地所有游戏货币
func (s *CurrencySyncService) fetchLocal(ctx context.Context) ([]*ent.GameCurrency, error) {
	var currencies []*ent.GameCurrency
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

// compareOne 比较单个游戏货币的本地和远程数据
func (s *CurrencySyncService) compareOne(remote *vendors.GameCurrencyInfo, local []*ent.GameCurrency) *platform_game.SyncDiff {
	// 初始化差异记录
	diff := &platform_game.SyncDiff{
		ObjectType: "currency",
		ObjectId:   remote.CurrencyId,
		// ObjectCode: remote.CurrencyID,
		RemoteId: remote.CurrencyId,
		// RemoteCode: remote.CurrencyID,
	}

	// 检查本地是否存在该游戏货币关系
	// 简化策略：按 game_id 检查是否已有货币配置
	exists := false
	for _, curr := range local {
		if curr.GameId == remote.GameId {
			exists = true
			break
		}
	}

	if !exists {
		diff.Action = "create"
		diff.Reason = "本地不存在该游戏货币配置"
		return diff
	}

	// 本地和远程一致，无需操作
	diff.Action = "noop"
	diff.Reason = "本地和远程数据一致"
	return diff
}

// ReverseSync 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *CurrencySyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteGameCurrencyResponse, apply *Apply) error {
	log.Println("[币种逆向同步] 开始执行逆向同步...")
	// 构建远程币种唯一性索引（game_id + currency_id）
	remoteIndex := make(map[string]bool)
	for _, remoteCurr := range remoteResp.Data {
		key := fmt.Sprintf("%d_%d", remoteCurr.GameId, remoteCurr.CurrencyId)
		remoteIndex[key] = true
	}

	// 获取本地币种关联数据
	var localCurrencies []*ent.GameCurrency
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localCurrencies).Error; err != nil {
		log.Printf("[币种逆向同步] ✗ 获取本地数据失败: %v\n", err)
		return err
	}
	log.Printf("[币种逆向同步] ✓ 获取本地数据成功, 共 %d 条\n", len(localCurrencies))

	// 在事务中执行软删除
	err := s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localCurr := range localCurrencies {
			// 如果本地币种在远程不存在，则软删除
			key := fmt.Sprintf("%d_%d", localCurr.GameId, localCurr.CurrencyId)
			if !remoteIndex[key] {
				// // 检查唯一性：是否存在其他 deleted_at IS NULL 的记录有相同的 game_id 和 currency_id
				// var duplicateCount int64
				// if err := tx.Model(&ent.GameCurrency{}).
				// 	Where("game_id = ?", localCurr.GameId).
				// 	Where("currency_id = ?", localCurr.CurrencyId).
				// 	Where("deleted_at IS NULL").
				// 	Count(&duplicateCount).Error; err != nil {
				// 	logger.Errorf("[币种逆向同步] 唯一性检查失败: %v", err)
				// 	apply.Failed++
				// 	continue
				// }

				// // duplicateCount >= 1 包括了 localCurr 本身，如果 > 1 说明有其他重复记录
				// if duplicateCount > 1 {
				// 	logger.Errorf("[币种逆向同步] 存在 %d 条记录有相同的 game_id=%d, currency_id=%d, 跳过本条删除以保护数据", duplicateCount, localCurr.GameId, localCurr.CurrencyId)
				// 	apply.Failed++
				// 	continue
				// }

				if err := tx.Model(localCurr).Update("deleted_at", time.Now()).Error; err != nil {
					logger.Errorf("[币种逆向同步] 软删除失败: %v", err)
					apply.Failed++
					continue
				}
				apply.Deleted++
				log.Printf("[币种逆向同步] ✓ 已软删除币种关联: GameID=%d, CurrencyID=%d\n", localCurr.GameId, localCurr.CurrencyId)
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	log.Printf("[币种逆向同步] ✓ 完成，删除: %d\n", apply.Deleted)
	return nil
}

// 获取远程游戏货币数据
func getRemoteGameCurrency(ctx context.Context, client vendors.VendorGameServiceClient, db *gorm.DB, isGetRemoteClient bool) (*RemoteGameCurrencyResponse, error) {
	// if isGetRemoteClient {
	// 	remoteResp, err := client.GetGameCurrency(ctx, &vendors.Empty{})
	// 	if err != nil {
	// 		logger.Errorf("❌ [币种逆向同步] 获取远程数据失败: %v", err)
	// 		return nil, err
	// 	}
	// 	return &RemoteGameCurrencyResponse{Data: remoteResp.GameCurrencyList}, nil
	// }
	// // 如果不使用远程客户端，可以在这里返回本地测试数据或其他逻辑
	// return &RemoteGameCurrencyResponse{Data: []*vendors.GameCurrencyInfo{}}, nil

	// if db == nil {
	// 	return nil, fmt.Errorf("getTestGameCurrencies: db 为空")
	// }

	// currencyID, err := ensureSysCurrencyID(db, "USD")
	// if err != nil {
	// 	return nil, err
	// }
	// return &RemoteGameCurrencyResponse{
	// 	Data: []*vendors.GameCurrencyInfo{
	// 		&vendors.GameCurrencyInfo{GameId: 1, CurrencyId: currencyID}, // USD
	// 		&vendors.GameCurrencyInfo{GameId: 3, CurrencyId: currencyID}, // USD
	// 		&vendors.GameCurrencyInfo{GameId: 2, CurrencyId: currencyID}, // USD
	// 		&vendors.GameCurrencyInfo{GameId: 4, CurrencyId: currencyID}, // USD
	// 		&vendors.GameCurrencyInfo{GameId: 5, CurrencyId: currencyID}, // USD
	// 	},
	// }, nil

	logger.Infof("[同步数据源] 使用本地 JSON 数据获取币种关联信息")
	var list []*vendors.GameCurrencyInfo
	if err := loadLocalJSON("currency.json", &list); err != nil {
		logger.Errorf("[同步数据源] 读取 currency.json 失败: %v", err)
		return nil, err
	}
	logger.Infof("[同步数据源] ✓ 读取本地 JSON 数据成功, 共 %d 条", len(list))
	return &RemoteGameCurrencyResponse{Data: list}, nil
}

// GetGameCurrencyResponse 获取游戏货币响应
type RemoteGameCurrencyResponse struct {
	Data []*vendors.GameCurrencyInfo
}
