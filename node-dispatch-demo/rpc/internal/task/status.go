package task

import (
	"context"
	"fmt"
	"time"

	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
)

// MarkRunning 标记任务执行中
func (s *Service) MarkRunning(ctx context.Context, taskNo string, runNo int64, nodeCode string) error {
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

	// 已经开始或结束的执行记录不重复处理
	if runData.Status != 1 {
		return nil
	}

	// 开启事务
	tx, err := s.db.Tx(ctx)
	if err != nil {
		return fmt.Errorf("开启事务失败: %w", err)
	}

	// 仅将待执行记录标记为执行中, 避免延迟ACK覆盖已完成状态
	affected, err := tx.DispatchTaskRun.
		Update().
		Where(
			dispatchtaskrun.IDEQ(runData.ID),
			dispatchtaskrun.StatusEQ(1),
		).
		SetStatus(2).             // 执行状态: 2执行中
		SetStartedAt(time.Now()). // 开始执行时间
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("更新任务执行状态失败: %w", err)
	}

	// 当前执行记录已经被其他消息更新时直接忽略本次ACK
	if affected == 0 {
		_ = tx.Rollback()
		return nil
	}

	// 仅将待执行任务标记为执行中, 避免覆盖已完成任务状态
	_, err = tx.DispatchTask.
		Update().
		Where(
			dispatchtask.IDEQ(taskData.ID),
			dispatchtask.StatusEQ(1),
		).
		SetStatus(2). // 任务状态: 2执行中
		Save(ctx)
	if err != nil {
		_ = tx.Rollback()
		return fmt.Errorf("更新任务状态失败: %w", err)
	}

	// 提交事务
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("提交事务失败: %w", err)
	}

	return nil
}
