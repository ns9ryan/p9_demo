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

// ===== 分类同步服务 =====

// CategorySyncService 分类同步服务
type CategorySyncService struct {
	DAOManager    *dao.Manager
	config        config.Config
	progressQueue *ProgressQueue
	logx.Logger
}

// NewCategorySyncService 创建分类同步服务
func NewCategorySyncService(ctx context.Context, config config.Config, daoManager *dao.Manager) *CategorySyncService {
	return &CategorySyncService{
		DAOManager:    daoManager,
		config:        config,
		progressQueue: NewProgressQueue(daoManager),
		Logger:        logx.WithContext(ctx),
	}
}

// 预检查分类同步
func (s *CategorySyncService) Preview(ctx context.Context, client vendors.VendorGameServiceClient, req *platform_game.SyncPreviewRequest, isGetRemoteClient bool) (*platform_game.SyncPreviewResp, *RemoteCategoryResponse, error) {
	s.Info("[分类预检查] 开始执行")

	// 获取远程分类数据
	s.Info("[分类预检查] 正在调用 GetGameCategory RPC")
	remoteResp, err := s.getRemoteGameCategory(ctx, client, isGetRemoteClient)
	if err != nil {
		s.Errorf("❌ [分类预检查] 获取测试数据失败: %v", err)
		return nil, nil, err
	}

	s.Infof("[分类预检查] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.Data))

	// 获取本地分类数据
	s.Info("[分类预检查] 正在获取本地数据")
	local, localIndex, err := s.fetchLocal(ctx)
	if err != nil {
		s.Errorf("❌ [分类预检查] 获取本地数据失败: %v", err)
		return nil, nil, err
	}
	s.Infof("[分类预检查] ✓ 获取本地数据成功, 共 %d 条", len(local))

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
				return s.compareAll(remoteCategory, localIdx)
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
			return s.compareAll(remoteCategory, localIdx)
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
func (s *CategorySyncService) processCategoryDataWithProgress(ctx context.Context, remoteData []*vendors.GameCategoryInfo, applyResult *Apply, syncCols []string, checkpointID int64) error {
	i := 0
	counterMap := make(map[string]int) // 用于记录每个 source_category_code 的计数，确保唯一性
	nameMap := make(map[string]string)
	createList := make([]*ent.GameCategoryCreate, 0, len(remoteData))
	for _, remoteCat := range remoteData {
		if counterMap[remoteCat.Code] > 0 {
			s.Errorf("[分类同步] 检测到重复的远程分类编码: %s, 计数器: %d, 跳过处理", remoteCat.Code, counterMap[remoteCat.Code])
			applyResult.Failed++
			continue
		}
		nameMap[remoteCat.Code] = remoteCat.Name
		// 更新计数器
		counterMap[remoteCat.Code]++
		i++
		gameCategoryRecord, err := s.DAOManager.GameCategory.GetGameCategoryBySourceCode(ctx, remoteCat.Code)
		categoryCode := remoteCat.Code + "_" + fmt.Sprintf("%d", i)
		if err != nil || gameCategoryRecord == nil {
			createList = append(createList, s.DAOManager.DB.GameCategory.Create().
				SetSourceID(remoteCat.Id).
				SetCategoryCode(categoryCode).
				SetSourceCategoryCode(remoteCat.Code).
				SetSortNo(remoteCat.Id).
				SetSourceSortNo(0).
				SetStatus(int64(remoteCat.Status)).
				SetSourceStatus(int64(remoteCat.Status)).
				SetCreatedAt(time.Now()).
				SetUpdatedAt(time.Now()))
		} else {
			if s.compare(&platform_game.SyncDiff{}, remoteCat, gameCategoryRecord) || len(syncCols) > 0 {
				s.Infof("[分类同步] 检测到需要更新的分类: %s", remoteCat.Code)
				update := s.DAOManager.DB.GameCategory.
					UpdateOneID(gameCategoryRecord.ID).
					SetUpdatedAt(time.Now())

				if len(syncCols) == 0 {
					syncCols = append(syncCols, "category_code")
					syncCols = append(syncCols, "source_category_code")
					syncCols = append(syncCols, "source_status")
					syncCols = append(syncCols, "source_sort_no")
				}

				for _, col := range syncCols {
					switch col {
					case "category_code":
						update.SetCategoryCode(categoryCode)
					case "source_category_code":
						update.SetSourceCategoryCode(remoteCat.Code)
					case "status":
						update.SetStatus(int64(remoteCat.Status))
					case "source_status":
						update.SetSourceStatus(int64(remoteCat.Status))
					case "sort_no":
						update.SetSortNo(remoteCat.Id)
					case "source_sort_no":
						update.SetSourceSortNo(0)
					}
				}

				_, err := update.Save(ctx)
				if err != nil {
					s.Errorf("[分类更新] ent 失败: %v", err)
					applyResult.Failed++
					continue
				}
				applyResult.Updated++
				s.Infof("[分类更新] ✓ 已更新分类: %s", remoteCat.Code)
			}
		}
	}
	if len(createList) > 0 {
		records, err := s.DAOManager.GameCategory.BatchCreateGameCategory(ctx, createList)
		if err != nil {
			s.Errorf("[分类新增] 批量创建失败: %v", err)
			applyResult.Failed += int32(len(createList))
		} else {
			applyResult.Created += int32(len(records))
		}
	}
	// 将 nameMap 写入本地 JSON 文件
	if err := locales.MergeLocalJSON("game_category_name_map.json", nameMap, constant.CategoryBiz); err != nil {
		s.Errorf("[分类同步] 保存名称映射失败: %v", err)
		// 不返回错误，继续执行后续逻辑
	}
	return nil
}

// Run 执行分类同步
func (s *CategorySyncService) Run(ctx context.Context, client vendors.VendorGameServiceClient, syncCols []string, checkpointID int64, isGetRemoteClient bool) error {
	s.Infof("[分类同步] ===== 开始执行同步 =====")
	// 先执行预检查
	s.Infof("[分类同步] 执行预检查")
	previewResp, remoteResp, err := s.Preview(ctx, client, nil, isGetRemoteClient)
	if err != nil {
		s.Errorf("[分类同步] 预检查失败: %v", err)
		return err
	}
	s.progressQueue.Start(ctx)
	defer s.progressQueue.Stop()
	if previewResp.Stats.UpdateTotal == 0 && previewResp.Stats.CreateTotal == 0 && previewResp.Stats.DeleteTotal == 0 && previewResp.Stats.RemoteTotal == previewResp.Stats.LocalTotal && len(syncCols) == 0 {
		s.Infof("[分类同步] 无需更新，直接返回")
		s.Infof("[分类同步] UpdateTotal=%d, CreateTotal=%d, DeleteTotal=%d, RemoteTotal=%d, LocalTotal=%d",
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
	// 把remoteResp.Data拆分为多条一批进行处理（可根据实际情况调整批次大小）
	batchSize := s.config.SyncBatchSize
	s.Infof("[分类同步] 预检查完成: remote_total=%d, local_total=%d, 检查点已创建: checkpointID=%d, batch_size=%d",
		previewResp.Stats.RemoteTotal, previewResp.Stats.LocalTotal, checkpointID, batchSize)
	totalCount := len(remoteResp.Data)
	apply := &Apply{}
	for i := 0; i < len(remoteResp.Data); i += batchSize {
		end := i + batchSize
		if end > len(remoteResp.Data) {
			end = len(remoteResp.Data)
		}
		batch := remoteResp.Data[i:end]
		s.Infof("[分类同步] 处理批次: start=%d, end=%d", i, end)
		// 模拟等待
		s.Debugf("[进度队列] 模拟处理延迟: 表=category, checkpointID=%d", checkpointID)
		time.Sleep(1000 * time.Millisecond)
		s.Debugf("[进度队列] 开始处理消息: 表=category, checkpointID=%d", checkpointID)
		err = s.processCategoryDataWithProgress(ctx, batch, apply, syncCols, checkpointID)
		if err != nil {
			s.Errorf("[分类同步] 处理分类数据失败: %v", err)
			return err
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
		s.Infof("[分类同步] 发送进度消息: 表=category, 处理数=%d, 总数=%d, 进度=%d%%", i+1, totalCount, progress)
		s.progressQueue.Send(msg)
	}

	// 执行逆向同步
	s.Infof("[分类同步] 开始执行逆向同步")
	_, err = s.ReverseSync(ctx, client, remoteResp, apply)
	if err != nil {
		s.Errorf("❌ [分类同步] 逆向同步失败: %v", err)
		return err
	}
	s.Infof("[分类同步] 逆向同步完成: deleted=%d", apply.Deleted)

	// 最后一次进度更新（100%）
	s.Infof("[分类同步] 发送最终进度消息（100%%）")
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

	s.Infof("[分类同步] ===== 同步完成 =====")

	return nil
}

// 获取本地所有分类
func (s *CategorySyncService) fetchLocal(ctx context.Context) ([]*ent.GameCategory, map[string]*ent.GameCategory, error) {
	categories, err := s.DAOManager.GameCategory.GetAllGameCategories(ctx)
	if err != nil {
		s.Errorf("fetchLocal fail: %d", err)
		return nil, nil, err
	}

	// 构建本地分类索引（按上游分类编码）
	index := make(map[string]*ent.GameCategory)
	for _, cat := range categories {
		index[cat.SourceCategoryCode] = cat
	}

	return categories, index, nil
}

// 比较单个分类的本地和远程数据
func (s *CategorySyncService) compareAll(remote *vendors.GameCategoryInfo, localIndex map[string]*ent.GameCategory) *platform_game.SyncDiff {
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

	s.compare(diff, remote, local)
	return diff
}

func (s *CategorySyncService) compare(diff *platform_game.SyncDiff, remote *vendors.GameCategoryInfo, local *ent.GameCategory) bool {
	if local.SourceStatus != int64(remote.Status) {
		diff.Action = constant.SyncActionUpdate
		diff.Reason = fmt.Sprintf("分类状态变更: %d -> %d", local.SourceStatus, remote.Status)
		diff.ConflictType = "status_changed"
		return true
	}
	diff.Action = constant.SyncActionNoop
	diff.Reason = "本地和远程数据一致"
	return false
}

// 逆向同步：检查本地数据在远程是否存在，不存在则软删除
func (s *CategorySyncService) ReverseSync(ctx context.Context, client vendors.VendorGameServiceClient, remoteResp *RemoteCategoryResponse, apply *Apply) (*Apply, error) {
	s.Infof("[分类逆向同步] 开始执行逆向同步")

	// 获取远程分类数据
	s.Infof("[分类逆向同步] ✓ 获取远程数据成功, 共 %d 条", len(remoteResp.Data))

	// 构建远程分类编码索引
	remoteIndex := make(map[string]bool)
	for _, remoteCat := range remoteResp.Data {
		remoteIndex[remoteCat.Code] = true
		s.Debugf("[分类逆向同步] 远程分类编码: %s", remoteCat.Code)
	}
	s.Infof("[分类逆向同步] 远程分类编码总数: %d", len(remoteIndex))

	// 获取本地分类数据
	localCategories, err := s.DAOManager.GameCategory.GetAllGameCategories(ctx)
	if err != nil {
		s.Errorf("❌ [分类逆向同步] 获取本地数据失败: %v", err)
		return nil, err
	}
	s.Infof("[分类逆向同步] ✓ 获取本地数据成功, 共 %d 条", len(localCategories))

	// 打印本地所有分类编码
	for _, localCat := range localCategories {
		s.Debugf("[分类逆向同步] 本地分类编码: %s (ID: %d, deleted_at: %v)",
			localCat.SourceCategoryCode, localCat.ID, localCat.DeletedAt)
	}

	for _, localCat := range localCategories {
		// 如果本地分类在远程不存在，则软删除
		if !remoteIndex[localCat.SourceCategoryCode] {
			_, err := s.DAOManager.GameCategory.UpdateGameCategory(ctx, localCat.ID, map[string]interface{}{"deleted_at": time.Now()})
			if err != nil {
				s.Errorf("❌ [分类逆向同步] 软删除失败: %v", err)
				apply.Failed++
				continue
			}
			apply.Deleted++
			s.Infof("[分类逆向同步] ✓ 已软删除分类: %s (ID: %d)", localCat.SourceCategoryCode, localCat.ID)
		}
	}

	// if err != nil {
	// 	return nil, err
	// }

	s.Infof("[分类逆向同步] ✓ 完成，删除: %d", apply.Deleted)
	return apply, nil
}

func (s *CategorySyncService) getRemoteGameCategory(ctx context.Context, client vendors.VendorGameServiceClient, isGetRemoteClient bool) (*RemoteCategoryResponse, error) {
	if isGetRemoteClient {
		s.Infof("[同步数据源] 正在获取远程分类数据...")
		remoteResp, err := client.GetGameCategory(ctx, &vendors.Empty{})
		if err != nil {
			s.Errorf("❌ [同步数据源] 获取远程数据失败: %v", err)
			return nil, err
		}
		s.Infof("[同步数据源] ✓ 获取远程分类数据成功, 共 %d 条", len(remoteResp.CategoryList))
		return &RemoteCategoryResponse{Data: remoteResp.CategoryList}, nil
	} else {
		s.Infof("[同步数据源] 正在获取本地分类数据...")
		var list []*vendors.GameCategoryInfo
		if err := locales.LoadLocalVendorRemoteJSON("category.json", &list); err != nil {
			s.Errorf("[同步数据源] 读取 category.json 失败: %v", err)
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
		s.Infof("[同步数据源] ✓ 获取本地分类数据成功, 共 %d 条", len(categoryList))
		return &RemoteCategoryResponse{Data: categoryList}, nil
	}
}

type RemoteCategoryResponse struct {
	Data []*vendors.GameCategoryInfo `json:"data"`
}
