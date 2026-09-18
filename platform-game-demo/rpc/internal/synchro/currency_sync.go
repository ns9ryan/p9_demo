package game_sync

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
	"oa.98ent.com/p9/platform-game/rpc/ent"
	"oa.98ent.com/p9/platform-game/rpc/internal/config"
	"oa.98ent.com/p9/platform-game/rpc/internal/constant"
	"oa.98ent.com/p9/platform-game/rpc/internal/dao"
	"oa.98ent.com/p9/platform-game/rpc/internal/locales"
	"oa.98ent.com/p9/platform-game/rpc/pb/platform_game"
	"oa.98ent.com/p9/platform-game/rpc/pb/vendors"
)

// ===== 游戏货币同步服务 =====

// CurrencySyncService 游戏货币同步服务
type CurrencySyncService struct {
	DAOManager    *dao.Manager
	config        config.Config
	progressQueue *ProgressQueue
	logx.Logger
}

// NewCurrencySyncService 创建游戏货币同步服务
func NewCurrencySyncService(ctx context.Context, config config.Config, daoManager *dao.Manager) *CurrencySyncService {
	return &CurrencySyncService{
		DAOManager:    daoManager,
		config:        config,
		progressQueue: NewProgressQueue(daoManager),
		Logger:        logx.WithContext(ctx),
	}
}

// Preview 预检查游戏货币同步
func (s *CurrencySyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteGameCurrencyResponse, error) {
	// 获取远程游戏货币数据
	remoteResp, err := s.getRemoteGameCurrency(ctx, client, s.DAOManager, isGetRemoteClient)
	if err != nil {
		s.Errorf("[货币预检查] 获取远程数据失败: %v", err)
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
				return s.compareAll(remoteCurrency, local)
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
			return s.compareAll(remoteCurrency, local)
		},
		PageSize: req.PageSize,
		Page:     req.Page,
		SkipNoop: req.IsSkip,
	}

	// 处理分页逻辑
	result := handler.Handle(int64(len(remote)), int64(len(local)))
	return result, remoteResp, nil
}

// 处理游戏货币数据的创建和更新（从 Run 提取的业务逻辑）
func (s *CurrencySyncService) processCurrencyDataWithProgress(ctx context.Context, remoteData []*vendors.GameCurrencyInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 currency key 的计数，确保唯一性
	createList := make([]*ent.GameCurrencyCreate, 0)
	for _, remoteCurrency := range remoteData {
		currencyKey := fmt.Sprintf("%d_%d", remoteCurrency.GameId, remoteCurrency.CurrencyId)
		if counterMap[currencyKey] > 0 {
			s.Errorf("[货币同步] 检测到重复的远程货币编码: GameID=%d, CurrencyID=%d, 计数器: %d, 跳过处理", remoteCurrency.GameId, remoteCurrency.CurrencyId, counterMap[currencyKey])
			applyResult.Failed++
			continue
		}
		// 更新计数器
		counterMap[currencyKey]++
		i++
		// 检查本地是否存在该游戏货币
		gameCurrencyRecord, err := s.DAOManager.GameCurrency.GetGameCurrencyByCurrcyIdAndGameId(ctx, remoteCurrency.GameId, remoteCurrency.CurrencyId)

		if err != nil || gameCurrencyRecord == nil {
			createList = append(createList, s.DAOManager.DB.GameCurrency.Create().
				SetGameID(remoteCurrency.GameId).
				SetCurrencyID(remoteCurrency.CurrencyId).
				SetSourceStatus(1).
				SetStatus(1).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()))
		} else {
			if s.compare(&platform_game.SyncDiff{}, remoteCurrency, gameCurrencyRecord) || len(syncCols) > 0 {
				s.Infof("[货币同步] 检测到需要更新的货币: GameID=%d, CurrencyID=%d", remoteCurrency.GameId, remoteCurrency.CurrencyId)
				update := s.DAOManager.DB.GameCurrency.
					UpdateOneID(gameCurrencyRecord.ID).
					SetUpdatedAt(time.Now())

				if len(syncCols) == 0 {
					syncCols = append(syncCols, "source_status")
				}
				if remoteCurrency.Status == 0 {
					remoteCurrency.Status = 1
				}

				for _, col := range syncCols {
					switch col {
					case "game_id":
						update.SetGameID(remoteCurrency.GameId)
					case "currency_id":
						update.SetCurrencyID(remoteCurrency.CurrencyId)
					case "source_status":
						update.SetSourceStatus(int64(remoteCurrency.Status))
					case "status":
						update.SetStatus(int64(remoteCurrency.Status))
					}
				}
				if _, err := update.Save(ctx); err != nil {
					s.Errorf("[货币同步] 更新失败: %v", err)
					applyResult.Failed++
					return err
				}
				applyResult.Updated++
			}
		}
	}
	s.Infof("[货币同步] 批量处理完成: Created=%d, Updated=%d, Failed=%d, created_list_len=%d", applyResult.Created, applyResult.Updated, applyResult.Failed, len(createList))
	if len(createList) > 0 {
		_, err := s.DAOManager.GameCurrency.BatchCreateGameCurrency(ctx, createList)
		if err != nil {
			s.Errorf("[货币新增] 批量创建失败: %v", err)
			applyResult.Failed += int32(len(createList))
		} else {
			applyResult.Created += int32(len(createList))
		}
	}

	return nil
}

// Run 执行游戏货币同步
func (s *CurrencySyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	s.Infof("[货币同步] ===== 开始执行同步 =====")

	// 先执行预检查
	s.Infof("[货币同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		s.Errorf("[货币同步] 预检查失败: %v", err)
		return err
	}
	// 创建并启动进度队列
	s.progressQueue.Start(ctx)
	defer s.progressQueue.Stop()
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal && len(syncCols) == 0 {
		s.Infof("[货币同步] 无需更新，直接返回")
		s.Infof("[货币同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
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
	// 把remoteResp.Data拆分为多条一批进行处理（可根据实际情况调整批次大小）
	batchSize := s.config.SyncBatchSize
	s.Infof("[货币同步] 预检查完成: remote_total=%d, local_total=%d, 检查点已创建: checkpointID=%d, batch_size=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal, checkpointID, batchSize)
	totalCount := len(remoteResp.Data)
	apply := &Apply{}
	for i := 0; i < len(remoteResp.Data); i += batchSize {
		end := i + batchSize
		if end > len(remoteResp.Data) {
			end = len(remoteResp.Data)
		}
		batch := remoteResp.Data[i:end]
		s.Infof("[货币同步] 处理批次: start=%d, end=%d", i, end)
		// 模拟等待
		s.Debugf("[进度队列] 模拟处理延迟: 表=currency, checkpointID=%d", checkpointID)
		time.Sleep(1000 * time.Millisecond)
		s.Debugf("[进度队列] 开始处理消息: 表=currency, checkpointID=%d", checkpointID)
		err = s.processCurrencyDataWithProgress(ctx, batch, apply, syncCols, checkpointID)
		if err != nil {
			s.Errorf("[货币同步] 处理货币数据失败: %v", err)
			return err
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
		s.Infof("[货币同步] 发送进度消息: 表=currency, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
		s.progressQueue.Send(msg)
	}
	if err != nil {
		s.Errorf("[货币同步] 处理货币数据失败: %v", err)
		return err
	}

	// 执行逆向同步
	s.Infof("[货币同步] 开始执行逆向同步")
	err = s.ReverseSync(ctx, client, remoteResp, apply)
	if err != nil {
		s.Errorf("❌ [货币同步] 逆向同步失败: %v", err)
		return err
	}
	s.Infof("[货币同步] 逆向同步完成: deleted=%d", apply.Deleted)

	// 最后一次进度更新（100%）
	s.Infof("[货币同步] 发送最终进度消息（100%%）")
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

	s.Infof("[货币同步] ===== 同步完成 =====")
	// 构建执行结果
	return nil
}

// 获取本地所有游戏货币
func (s *CurrencySyncService) fetchLocal(ctx context.Context) ([]*ent.GameCurrency, error) {
	currencies, err := s.DAOManager.GameCurrency.GetAllGameCurrency(ctx, 0, 0)
	if err != nil {
		return nil, err
	}
	return currencies, nil
}

// 比较单个游戏货币的本地和远程数据
func (s *CurrencySyncService) compareAll(remote *vendors.GameCurrencyInfo, local []*ent.GameCurrency) *platform_game.SyncDiff {
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
		if curr.GameID == remote.GameId && curr.CurrencyID == remote.CurrencyId {
			s.compare(diff, remote, curr)
			exists = true
			break
		}
	}

	if !exists {
		diff.Action = "create"
		diff.Reason = "本地不存在该游戏货币配置"
		return diff
	}
	return diff
}

// 比较单个游戏货币的本地和远程数据
func (s *CurrencySyncService) compare(diff *platform_game.SyncDiff, remote *vendors.GameCurrencyInfo, local *ent.GameCurrency) bool {
	if local.SourceStatus != int64(remote.Status) {
		diff.Action = constant.SyncActionUpdate
		diff.Reason = fmt.Sprintf("状态变更: %d -> %d", local.SourceStatus, remote.Status)
		diff.ConflictType = "status_changed"
		return true
	}
	diff.Action = constant.SyncActionNoop
	diff.Reason = "本地和远程数据一致"
	return false
}

// ReverseSync 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *CurrencySyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteGameCurrencyResponse, apply *Apply) error {
	s.Infof("[币种逆向同步] 开始执行逆向同步...")
	// 构建远程币种唯一性索引（game_id + currency_id）
	remoteIndex := make(map[string]bool)
	for _, remoteCurr := range remoteResp.Data {
		key := fmt.Sprintf("%d_%d", remoteCurr.GameId, remoteCurr.CurrencyId)
		remoteIndex[key] = true
	}

	// 获取本地币种关联数据
	localCurrencies, err := s.DAOManager.GameCurrency.GetAllGameCurrency(ctx, 0, 0)
	if err != nil {
		s.Errorf("[币种逆向同步] ✗ 获取本地数据失败: %v", err)
		return err
	}
	s.Infof("[币种逆向同步] ✓ 获取本地数据成功, 共 %d 条", len(localCurrencies))

	for _, localCurr := range localCurrencies {
		// 如果本地币种在远程不存在，则软删除
		key := fmt.Sprintf("%d_%d", localCurr.GameID, localCurr.CurrencyID)
		if !remoteIndex[key] {
			_, err := s.DAOManager.GameCurrency.UpdateGameCurrency(ctx, localCurr.ID, map[string]interface{}{"deleted_at": time.Now()})
			if err != nil {
				s.Errorf("[币种逆向同步] 软删除失败: %v", err)
				apply.Failed++
				continue
			}
			apply.Deleted++
			s.Infof("[币种逆向同步] ✓ 已软删除币种关联: GameID=%d, CurrencyID=%d", localCurr.GameID, localCurr.CurrencyID)
		}
	}
	if err != nil {
		return err
	}

	s.Infof("[币种逆向同步] ✓ 完成，删除: %d", apply.Deleted)
	return nil
}

// 获取远程游戏货币数据
func (s *CurrencySyncService) getRemoteGameCurrency(ctx context.Context, client vendors.VendorGameServiceClient, daoManager *dao.Manager, isGetRemoteClient bool) (*RemoteGameCurrencyResponse, error) {
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

	s.Infof("[同步数据源] 使用本地 JSON 数据获取币种关联信息")
	var list []*vendors.GameCurrencyInfo
	if err := locales.LoadLocalVendorRemoteJSON("currency.json", &list); err != nil {
		s.Errorf("[同步数据源] 读取 currency.json 失败: %v", err)
		return nil, err
	}
	s.Infof("[同步数据源] ✓ 读取本地 JSON 数据成功, 共 %d 条", len(list))
	return &RemoteGameCurrencyResponse{Data: list}, nil
}

// GetGameCurrencyResponse 获取游戏货币响应
type RemoteGameCurrencyResponse struct {
	Data []*vendors.GameCurrencyInfo
}
