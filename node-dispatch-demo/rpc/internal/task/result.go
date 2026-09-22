package task

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
)

// MarkResult 标记任务执行结果
func (s *Service) MarkResult(
	ctx context.Context,
	taskNo string,
	runNo int64,
	nodeCode string,
	success bool,
	result []byte,
	errorMessage string,
) error {
	// 查询调度任务
	taskData, err := s.db.DispatchTask.
		Query().
		Where(dispatchtask.TaskNoEQ(taskNo)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("查询调度任务失败: %w", err)
	}

	// 查询任务执行记录
	runData, err := s.db.DispatchTaskRun.
		Query().
		Where(
			dispatchtaskrun.TaskIDEQ(taskData.ID),
			dispatchtaskrun.RunNoEQ(runNo),
		).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("查询任务执行记录失败: %w", err)
	}

	// 查询执行节点
	nodeData, err := s.db.Node.
		Query().
		Where(node.IDEQ(runData.NodeID)).
		Only(ctx)
	if err != nil {
		return fmt.Errorf("查询执行节点失败: %w", err)
	}

	// 校验执行节点
	if nodeData.Code != nodeCode {
		return fmt.Errorf("任务执行节点不匹配: expect=%s actual=%s", nodeData.Code, nodeCode)
	}

	// 已完成的执行记录不重复处理
	if runData.Status == 3 || runData.Status == 4 {
		return nil
	}

	// 开启事务
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}

	// 创建执行结果更新
	runUpdate := tx.DispatchTaskRun.
		Update().
		Where(
			dispatchtaskrun.IDEQ(runData.ID),
			dispatchtaskrun.StatusIn(1, 2),
		).
		SetFinishedAt(time.Now()) // 执行结束时间

	if success {
		// 标记执行成功
		runUpdate.
			SetStatus(3).     // 执行状态: 3成功
			SetResult(result) // 执行结果
	} else {
		// 标记执行失败
		runUpdate.
			SetStatus(4).                 // 执行状态: 4失败
			SetErrorMessage(errorMessage) // 执行失败原因
	}

	// 保存执行结果
	affected, err := runUpdate.Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("更新任务执行结果失败: %w", err)
	}

	// 执行记录已经被其他结果处理时直接忽略
	if affected == 0 {
		_ = tx.Rollback()
		return nil
	}

	// 当前执行记录后面存在新的执行记录时, 不覆盖任务整体状态
	hasNewerRun, err := tx.DispatchTaskRun.
		Query().
		Where(
			dispatchtaskrun.TaskIDEQ(taskData.ID),
			dispatchtaskrun.RunNoGT(runNo),
		).
		Exist(ctx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("检查后续任务执行记录失败: %w", err)
	}

	// 当前为最新执行记录时更新任务整体状态
	if !hasNewerRun {
		taskStatus := int64(4)
		if success {
			taskStatus = 3
		}

		if err = tx.DispatchTask.
			UpdateOneID(taskData.ID).
			SetStatus(taskStatus). // 任务状态: 3成功, 4失败
			Exec(ctx); err != nil {
			_ = tx.Rollback()
			return fmt.Errorf("更新任务状态失败: %w", err)
		}
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}
