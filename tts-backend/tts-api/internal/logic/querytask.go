package logic

import (
	"context"

	"tts-backend/tts-api/internal/svc"
	"tts-backend/tts-api/internal/types"
)

type QueryTaskLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewQueryTaskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QueryTaskLogic {
	return &QueryTaskLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QueryTaskLogic) QueryTask(taskId string) (resp *types.TaskResp, err error) {
	task, err := l.svcCtx.TaskModel.FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}

	return &types.TaskResp{
		TaskId:   task.TaskId,
		Status:   task.Status,
		Progress: task.Progress,
		AudioUrl: task.AudioUrl,
	}, nil
}
