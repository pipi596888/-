package logic

import (
	"context"
	"fmt"
	"github.com/google/uuid"

	"tts-backend/tts-api/internal/svc"
	"tts-backend/tts-api/internal/types"
	"tts-backend/tts-api/internal/model"
)

type GenerateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGenerateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateLogic {
	return &GenerateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GenerateLogic) Generate(req *types.GenerateReq) (resp *types.GenerateResp, err error) {
	taskId := uuid.New().String()

	userId := int64(1)
	totalChars := 0
	for _, seg := range req.Segments {
		totalChars += len(seg.Text)
	}

	task := &model.TtsTask{
		TaskId:  taskId,
		UserId:  userId,
		Status:  "pending",
		Progress: 0,
		Format:  req.Format,
		Channel: req.Channel,
	}

	_, err = l.svcCtx.TaskModel.Insert(task)
	if err != nil {
		return nil, err
	}

	segments := make([]*model.TtsSegment, 0, len(req.Segments))
	for i, seg := range req.Segments {
		segments = append(segments, &model.TtsSegment{
			TaskId:  taskId,
			VoiceId: seg.VoiceId,
			Emotion: seg.Emotion,
			Text:    seg.Text,
			Sort:    i,
		})
	}

	err = l.svcCtx.SegmentModel.BatchInsert(segments)
	if err != nil {
		return nil, err
	}

	fmt.Printf("Task %s created with %d segments, total chars: %d\n", taskId, len(req.Segments), totalChars)

	return &types.GenerateResp{
		TaskId: taskId,
	}, nil
}
