package sync

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gorm.io/gorm"
	"oa.98ent.com/p9/platform-game/common/logger"
	"oa.98ent.com/p9/platform-game/common/model"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 本地结构体（RPC 尚未实现，用于测试） =====

// GameCurrencyInfo 游戏货币信息
type GameCurrencyInfo struct {
	GameId     int64
	CurrencyID int64
}

// GetGameCurrencyResponse 获取游戏货币响应
type GetGameCurrencyResponse struct {
	Data []*GameCurrencyInfo
}

// ===== 游戏货币同步服务 =====

// CurrencySyncService 游戏货币同步服务
type CurrencySyncService struct {
	db *gorm.DB
}

// NewCurrencySyncService 创建游戏货币同步服务
func NewCurrencySyncService(db *gorm.DB) *CurrencySyncService {
	return &CurrencySyncService{db: db}
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
func (s *CurrencySyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncPreviewResp, error) {
	// 获取远程游戏货币数据
	// remoteResp, err := client.GetGameCurrency(ctx, &emptypb.Empty{})
	remoteResp, err := getTestGameCurrencies(s.db)
	if err != nil {
		logger.Errorf("[货币预检查] 获取远程数据失败: %v", err)
		return nil, err
	}
	remote := remoteResp.Data

	// 获取本地游戏货币数据
	local, err := s.fetchLocal(ctx)
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

	// 对每条远程游戏货币进行比较
	for _, remoteItem := range remote {
		diff := s.compareOne(remoteItem, local)
		result.Diffs = append(result.Diffs, diff)

		// 按操作类型统计
		switch diff.Action {
		case "create":
			result.Stats.CreateTotal++
		case "noop":
			result.Stats.NoopTotal++
		}
	}

	return result, nil
}

// processCurrencyData 处理游戏货币数据的创建和更新（从 Run 提取的业务逻辑）
func (s *CurrencySyncService) processCurrencyData(ctx context.Context, tx *gorm.DB, remoteData []*GameCurrencyInfo, applyResult *vendors.SyncApplyResult) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 currency key 的计数，确保唯一性
	for _, remoteCurrency := range remoteData {
		currencyKey := fmt.Sprintf("%d_%d", remoteCurrency.GameId, remoteCurrency.CurrencyID)
		if counterMap[currencyKey] > 0 {
			logger.Errorf("[货币同步] 检测到重复的远程货币编码: GameID=%d, CurrencyID=%d, 计数器: %d, 跳过处理", remoteCurrency.GameId, remoteCurrency.CurrencyID, counterMap[currencyKey])
			applyResult.Failed++
			continue
		}
		// 更新计数器
		counterMap[currencyKey]++
		i++
		// 检查本地是否存在该游戏货币
		var localCurrency model.GameCurrency
		exists := tx.Where("game_id = ?", remoteCurrency.GameId).Where("currency_id = ?", remoteCurrency.CurrencyID).
			Where("deleted_at IS NULL").First(&localCurrency).Error == nil

		if !exists {
			// 新增游戏货币
			newCurrency := model.GameCurrency{
				GameID:       remoteCurrency.GameId,
				CurrencyID:   remoteCurrency.CurrencyID,
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
			// 更新游戏货币 - 只更新源状态
			if err := tx.Model(&localCurrency).
				Updates(map[string]interface{}{
					"source_status": 1,
					"updated_at":    time.Now(),
				}).Error; err != nil {
				logger.Errorf("[货币更新] 失败: %v", err)
				return err
			}
			applyResult.Updated++
		}
	}
	return nil
}

// Run 执行游戏货币同步
func (s *CurrencySyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	logger.Infof("[货币同步] ===== 开始执行同步 =====")

	// 先执行预检查
	logger.Infof("[货币同步] 执行预检查")
	previewResp, err := s.Preview(ctx, client)
	if err != nil {
		logger.Errorf("[货币同步] 预检查失败: %v", err)
		return nil, err
	}
	logger.Infof("[货币同步] 预检查完成: remote_total=%d, local_total=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal)

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: previewResp,
		Apply:   &vendors.SyncApplyResult{},
	}

	// 获取远程游戏货币数据
	logger.Infof("[货币同步] 获取远程货币数据")
	remoteResp, err := getTestGameCurrencies(s.db)
	if err != nil {
		logger.Errorf("[货币同步] 获取远程数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[货币同步] 获取到远程数据: count=%d", len(remoteResp.Data))

	// 在事务中执行数据库操作
	logger.Infof("[货币同步] 开始处理货币数据")
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		return s.processCurrencyData(ctx, tx, remoteResp.Data, result.Apply)
	})
	if err != nil {
		logger.Errorf("[货币同步] 处理货币数据失败: %v", err)
		return nil, err
	}
	logger.Infof("[货币同步] 处理完成: created=%d, updated=%d, deleted=%d, failed=%d",
		result.Apply.Created, result.Apply.Updated, result.Apply.Deleted, result.Apply.Failed)

	// 执行逆向同步
	logger.Infof("[货币同步] 开始执行逆向同步")
	reverseResp, err := s.ReverseSync(ctx, client)
	if err != nil {
		logger.Errorf("[货币同步] 逆向同步失败: %v", err)
		return nil, err
	}
	result.Apply.Deleted = reverseResp.Apply.Deleted
	logger.Infof("[货币同步] 逆向同步完成: deleted=%d", result.Apply.Deleted)

	// 保存同步检查点
	logger.Infof("[货币同步] 准备保存同步检查点")
	checkpointMgr := NewCheckpointManager(s.db)
	checkpointValue := previewResp.Stats.RemoteTotal // 使用远程总数作为检查点值
	if err := checkpointMgr.UpdateCheckpoint(ctx, "CURRENCY",
		fmt.Sprintf("%d", checkpointValue), result, nil); err != nil {
		logger.Errorf("[货币同步] 保存检查点失败: %v", err)
		return nil, err
	}
	logger.Infof("[货币同步] ✓ 检查点已保存")

	logger.Infof("[货币同步] ===== 同步完成 =====")
	return result, nil
}

// fetchLocal 获取本地所有游戏货币
func (s *CurrencySyncService) fetchLocal(ctx context.Context) ([]*model.GameCurrency, error) {
	var currencies []*model.GameCurrency
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&currencies).Error; err != nil {
		return nil, err
	}
	return currencies, nil
}

// compareOne 比较单个游戏货币的本地和远程数据
func (s *CurrencySyncService) compareOne(remote *GameCurrencyInfo, local []*model.GameCurrency) *vendors.SyncDiff {
	// 初始化差异记录
	diff := &vendors.SyncDiff{
		ObjectType: "currency",
		ObjectId:   remote.GameId,
		// ObjectCode: remote.CurrencyID,
		RemoteId: remote.GameId,
		// RemoteCode: remote.CurrencyID,
	}

	// 检查本地是否存在该游戏货币关系
	// 简化策略：按 game_id 检查是否已有货币配置
	exists := false
	for _, curr := range local {
		if curr.GameID == remote.GameId {
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
func (s *CurrencySyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient) (*vendors.SyncRunResp, error) {
	log.Println("[币种逆向同步] 开始执行逆向同步...")

	// 获取远程币种数据
	// remoteResp, err := client.GetGameCurrency(ctx, &emptypb.Empty{})
	remoteResp, err := getTestGameCurrencies(s.db)
	if err != nil {
		log.Printf("[币种逆向同步] 获取测试数据失败: %v\n", err)
		return nil, err
	}
	log.Printf("[币种逆向同步] ✓ 获取远程数据成功, 共 %d 条\n", len(remoteResp.Data))

	// 构建远程币种唯一性索引（game_id + currency_id）
	remoteIndex := make(map[string]bool)
	for _, remoteCurr := range remoteResp.Data {
		key := fmt.Sprintf("%d_%d", remoteCurr.GameId, remoteCurr.CurrencyID)
		remoteIndex[key] = true
	}

	// 获取本地币种关联数据
	var localCurrencies []*model.GameCurrency
	if err := s.db.WithContext(ctx).Where("deleted_at IS NULL").Find(&localCurrencies).Error; err != nil {
		log.Printf("[币种逆向同步] ✗ 获取本地数据失败: %v\n", err)
		return nil, err
	}
	log.Printf("[币种逆向同步] ✓ 获取本地数据成功, 共 %d 条\n", len(localCurrencies))

	// 构建执行结果
	result := &vendors.SyncRunResp{
		Preview: &vendors.SyncPreviewResp{
			Stats: &vendors.SyncStats{
				LocalTotal: int64(len(localCurrencies)),
			},
		},
		Apply: &vendors.SyncApplyResult{},
	}

	// 在事务中执行软删除
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, localCurr := range localCurrencies {
			// 如果本地币种在远程不存在，则软删除
			key := fmt.Sprintf("%d_%d", localCurr.GameID, localCurr.CurrencyID)
			if !remoteIndex[key] {
				// 检查唯一性：是否存在其他 deleted_at IS NULL 的记录有相同的 game_id 和 currency_id
				var duplicateCount int64
				if err := tx.Model(&model.GameCurrency{}).
					Where("game_id = ?", localCurr.GameID).
					Where("currency_id = ?", localCurr.CurrencyID).
					Where("deleted_at IS NULL").
					Count(&duplicateCount).Error; err != nil {
					logger.Errorf("[币种逆向同步] 唯一性检查失败: %v", err)
					result.Apply.Failed++
					continue
				}

				// duplicateCount >= 1 包括了 localCurr 本身，如果 > 1 说明有其他重复记录
				if duplicateCount > 1 {
					logger.Errorf("[币种逆向同步] 存在 %d 条记录有相同的 game_id=%d, currency_id=%d, 跳过本条删除以保护数据", duplicateCount, localCurr.GameID, localCurr.CurrencyID)
					result.Apply.Failed++
					continue
				}

				if err := tx.Model(localCurr).Update("deleted_at", time.Now()).Error; err != nil {
					logger.Errorf("[币种逆向同步] 软删除失败: %v", err)
					result.Apply.Failed++
					continue
				}
				result.Apply.Deleted++
				log.Printf("[币种逆向同步] ✓ 已软删除币种关联: GameID=%d, CurrencyID=%d\n", localCurr.GameID, localCurr.CurrencyID)
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	log.Printf("[币种逆向同步] ✓ 完成，删除: %d\n", result.Apply.Deleted)
	return result, nil
}

// getTestGameCurrencies 获取测试用的游戏货币数据（RPC 尚未实现）
func getTestGameCurrencies(db *gorm.DB) (*GetGameCurrencyResponse, error) {
	if db == nil {
		return nil, fmt.Errorf("getTestGameCurrencies: db 为空")
	}

	currencyID, err := ensureSysCurrencyID(db, "USD")
	if err != nil {
		return nil, err
	}
	return &GetGameCurrencyResponse{
		Data: []*GameCurrencyInfo{
			{GameId: 1, CurrencyID: currencyID}, // USD
			{GameId: 3, CurrencyID: currencyID}, // USD
			{GameId: 2, CurrencyID: currencyID}, // USD
			{GameId: 4, CurrencyID: currencyID}, // USD
			{GameId: 5, CurrencyID: currencyID}, // USD
		},
	}, nil
}
