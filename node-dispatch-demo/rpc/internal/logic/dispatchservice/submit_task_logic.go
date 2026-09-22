package dispatchservicelogic

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"oa.98ent.com/p9/node-dispatch/rpc/ent"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtask"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/dispatchtaskrun"
	"oa.98ent.com/p9/node-dispatch/rpc/ent/node"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/connection"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/protocol"
	"oa.98ent.com/p9/node-dispatch/rpc/internal/svc"
	"oa.98ent.com/p9/node-dispatch/rpc/pb/nodedispatchrpc/dispatchpb"

	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SubmitTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSubmitTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SubmitTaskLogic {
	return &SubmitTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// SubmitTask 提交调度任务
func (l *SubmitTaskLogic) SubmitTask(in *dispatchpb.SubmitTaskRequest) (*dispatchpb.SubmitTaskResponse, error) {
	// ---------- 参数整理 ----------

	// 整理任务参数
	requestNo := strings.TrimSpace(in.RequestNo)
	target := strings.TrimSpace(in.Target)
	taskType := strings.TrimSpace(in.TaskType)
	nodeCode := strings.TrimSpace(in.NodeCode)
	params := json.RawMessage(in.Params)

	// 校验基础参数
	if requestNo == "" {
		return nil, status.Error(codes.InvalidArgument, "request_no is required")
	}
	if target == "" {
		return nil, status.Error(codes.InvalidArgument, "target is required")
	}
	if taskType == "" {
		return nil, status.Error(codes.InvalidArgument, "task_type is required")
	}
	if nodeCode == "" {
		return nil, status.Error(codes.InvalidArgument, "node_code is required")
	}
	if !json.Valid(params) {
		return nil, status.Error(codes.InvalidArgument, "params must be valid json")
	}

	// ---------- 幂等检查 ----------

	// 根据请求编号检查任务是否已经存在
	existingTask, err := l.svcCtx.DB.DispatchTask.
		Query().
		Where(dispatchtask.RequestNoEQ(requestNo)).
		Only(l.ctx)
	if err == nil {
		// 获取任务首次执行记录
		existingRun, runErr := l.svcCtx.DB.DispatchTaskRun.
			Query().
			Where(
				dispatchtaskrun.TaskIDEQ(existingTask.ID),
				dispatchtaskrun.RunNoEQ(1),
			).
			Only(l.ctx)
		if runErr != nil {
			l.Logger.Errorw(
				"获取已有调度任务执行记录失败",
				logx.Field("request_no", requestNo),
				logx.Field("task_no", existingTask.TaskNo),
				logx.Field("error", runErr.Error()),
			)
			return nil, runErr
		}

		// 获取首次执行节点
		existingNode, nodeErr := l.svcCtx.DB.Node.Get(l.ctx, existingRun.NodeID)
		if nodeErr != nil {
			l.Logger.Errorw(
				"获取已有调度任务节点失败",
				logx.Field("request_no", requestNo),
				logx.Field("task_no", existingTask.TaskNo),
				logx.Field("error", nodeErr.Error()),
			)
			return nil, nodeErr
		}

		// 比较任务参数
		sameParams, compareErr := equalJSON(existingTask.Params, params)
		if compareErr != nil {
			l.Logger.Errorw(
				"比较调度任务参数失败",
				logx.Field("request_no", requestNo),
				logx.Field("task_no", existingTask.TaskNo),
				logx.Field("error", compareErr.Error()),
			)
			return nil, compareErr
		}

		// 相同请求编号必须保持任务内容一致
		if existingTask.Target != target ||
			existingTask.TaskType != taskType ||
			existingNode.Code != nodeCode ||
			!sameParams {
			return nil, status.Error(codes.AlreadyExists, "request_no already exists with different content")
		}

		// 相同请求直接返回已有任务编号, 不重复创建和下发
		return &dispatchpb.SubmitTaskResponse{
			TaskNo: existingTask.TaskNo, // 调度任务编号
		}, nil
	}
	if !ent.IsNotFound(err) {
		l.Logger.Errorw("查询调度任务失败", logx.Field("request_no", requestNo), logx.Field("error", err.Error()))
		return nil, err
	}

	// ---------- 节点校验 ----------

	// 获取执行节点
	nodeData, err := l.svcCtx.DB.Node.
		Query().
		Where(node.CodeEQ(nodeCode)).
		Only(l.ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, "node not found")
		}

		l.Logger.Errorw("查询执行节点失败", logx.Field("node_code", nodeCode), logx.Field("error", err.Error()))
		return nil, err
	}

	// 停用节点不能接收新的调度任务
	if nodeData.Status != 1 {
		return nil, status.Error(codes.FailedPrecondition, "node is disabled")
	}

	// 离线节点暂不创建调度任务
	if !l.svcCtx.Connections.IsOnline(nodeCode) {
		return nil, status.Error(codes.Unavailable, "node is offline")
	}

	// ---------- 任务创建 ----------

	// 生成调度中心任务编号
	taskNo := uuid.NewString()

	// 开启事务
	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		l.Logger.Errorw("开启调度任务事务失败", logx.Field("request_no", requestNo), logx.Field("error", err.Error()))
		return nil, err
	}

	// 创建调度任务
	taskData, err := tx.DispatchTask.
		Create().
		SetTaskNo(taskNo).       // 调度任务编号
		SetRequestNo(requestNo). // 调用方请求编号
		SetTarget(target).       // 目标服务
		SetTaskType(taskType).   // 任务类型
		SetParams(params).       // 任务参数
		Save(l.ctx)
	if err != nil {
		_ = tx.Rollback()
		l.Logger.Errorw(
			"创建调度任务失败",
			logx.Field("request_no", requestNo),
			logx.Field("task_no", taskNo),
			logx.Field("error", err.Error()),
		)
		return nil, err
	}

	// 创建首次任务执行记录
	runData, err := tx.DispatchTaskRun.
		Create().
		SetTaskID(taskData.ID). // 调度任务ID
		SetNodeID(nodeData.ID). // 执行节点ID
		SetRunNo(1).            // 首次执行序号
		Save(l.ctx)
	if err != nil {
		_ = tx.Rollback()
		l.Logger.Errorw("创建调度任务执行记录失败", logx.Field("task_no", taskNo), logx.Field("error", err.Error()))
		return nil, err
	}

	// 提交任务创建事务
	if err = tx.Commit(); err != nil {
		l.Logger.Errorw("提交调度任务事务失败", logx.Field("task_no", taskNo), logx.Field("error", err.Error()))
		return nil, err
	}

	// ---------- 任务下发 ----------

	// 编码任务下发数据
	dispatchData, err := json.Marshal(protocol.TaskDispatchData{
		TaskNo:   taskNo,        // 调度任务编号
		RunNo:    runData.RunNo, // 执行序号
		Target:   target,        // 目标服务
		TaskType: taskType,      // 任务类型
		Params:   params,        // 任务参数
	})
	if err != nil {
		l.Logger.Errorw("编码任务下发数据失败", logx.Field("task_no", taskNo), logx.Field("error", err.Error()))
		l.markDispatchFailed(taskData.ID, runData.ID, err.Error())
		return nil, err
	}

	// 通过WebSocket下发任务
	err = l.svcCtx.Connections.Send(l.ctx, nodeCode, protocol.Message{
		Type: protocol.MessageTypeTaskDispatch, // 消息类型
		Data: dispatchData,                     // 任务数据
	})
	if err != nil {
		l.Logger.Errorw(
			"下发调度任务失败",
			logx.Field("task_no", taskNo),
			logx.Field("run_no", runData.RunNo),
			logx.Field("node_code", nodeCode),
			logx.Field("error", err.Error()),
		)
		l.markDispatchFailed(taskData.ID, runData.ID, err.Error())

		if errors.Is(err, connection.ErrNodeOffline) {
			return nil, status.Error(codes.Unavailable, "node is offline")
		}

		return nil, status.Error(codes.Unavailable, "failed to dispatch task")
	}

	l.Logger.Infow(
		"调度任务已下发",
		logx.Field("task_no", taskNo),
		logx.Field("run_no", runData.RunNo),
		logx.Field("node_code", nodeCode),
	)

	// ---------- 返回结果 ----------

	// 返回提交结果
	return &dispatchpb.SubmitTaskResponse{
		TaskNo: taskNo, // 调度任务编号
	}, nil
}

// markDispatchFailed 标记任务下发失败
func (l *SubmitTaskLogic) markDispatchFailed(taskID, runID int64, errorMessage string) {
	// ---------- 状态更新 ----------

	// 开启失败状态更新事务
	tx, err := l.svcCtx.DB.Tx(l.ctx)
	if err != nil {
		l.Logger.Errorw("开启任务失败状态事务失败", logx.Field("task_id", taskID), logx.Field("error", err.Error()))
		return
	}

	// 标记调度任务失败
	err = tx.DispatchTask.
		UpdateOneID(taskID).
		SetStatus(4).
		Exec(l.ctx)
	if err != nil {
		_ = tx.Rollback()
		l.Logger.Errorw("更新调度任务失败状态失败", logx.Field("task_id", taskID), logx.Field("error", err.Error()))
		return
	}

	// 标记任务执行记录失败
	err = tx.DispatchTaskRun.
		UpdateOneID(runID).
		SetStatus(4).
		SetErrorMessage(errorMessage).
		SetFinishedAt(time.Now()).
		Exec(l.ctx)
	if err != nil {
		_ = tx.Rollback()
		l.Logger.Errorw("更新任务执行记录失败状态失败", logx.Field("run_id", runID), logx.Field("error", err.Error()))
		return
	}

	// 提交失败状态
	if err = tx.Commit(); err != nil {
		l.Logger.Errorw("提交任务失败状态事务失败", logx.Field("task_id", taskID), logx.Field("error", err.Error()))
	}
}

// equalJSON 比较JSON内容是否一致
func equalJSON(left, right json.RawMessage) (bool, error) {
	var leftValue any
	if err := json.Unmarshal(left, &leftValue); err != nil {
		return false, err
	}

	var rightValue any
	if err := json.Unmarshal(right, &rightValue); err != nil {
		return false, err
	}

	return reflect.DeepEqual(leftValue, rightValue), nil
}
