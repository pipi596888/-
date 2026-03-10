package logic

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"

	"tts-backend/tts-api/internal/model"
	"tts-backend/tts-api/internal/svc"
	"tts-backend/tts-api/internal/types"
)

var ErrInvalidVoice = errors.New("invalid voice")

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

func (l *GenerateLogic) Generate(req *types.GenerateReq, userId int64, isAdmin bool) (resp *types.GenerateResp, err error) {
	voiceIds := make([]int64, 0, len(req.Segments))
	for _, seg := range req.Segments {
		voiceIds = append(voiceIds, seg.VoiceId)
	}
	if err := l.svcCtx.VoiceAccess.ValidateVoiceIds(userId, isAdmin, voiceIds); err != nil {
		if errors.Is(err, model.ErrVoiceForbidden) {
			return nil, ErrForbidden
		}
		if errors.Is(err, model.ErrVoiceNotFound) || errors.Is(err, model.ErrInvalidVoiceID) {
			return nil, ErrInvalidVoice
		}
		return nil, ErrInvalidVoice
	}

	taskId := uuid.New().String()
	totalChars := 0
	for _, seg := range req.Segments {
		totalChars += len(seg.Text)
	}

	task := &model.TtsTask{
		TaskId:   taskId,
		UserId:   userId,
		Status:   "pending",
		Progress: 0,
		Format:   req.Format,
		Channel:  req.Channel,
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
