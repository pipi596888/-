package logic

import (
	"context"

	"tts-backend/voice-api/internal/svc"
	"tts-backend/voice-api/internal/types"
	"tts-backend/voice-api/internal/model"
)

type CreateVoiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateVoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateVoiceLogic {
	return &CreateVoiceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateVoiceLogic) CreateVoice(req *types.CreateVoiceReq) (resp *types.CreateVoiceResp, err error) {
	voice := &model.Voice{
		Name:       req.Name,
		Tone:       req.Tone,
		Gender:     req.Gender,
		PreviewUrl: req.PreviewUrl,
		IsDefault:  false,
	}

	id, err := l.svcCtx.VoiceModel.Insert(voice)
	if err != nil {
		return nil, err
	}

	return &types.CreateVoiceResp{
		Voice: types.Voice{
			Id:         id,
			Name:       req.Name,
			Tone:       req.Tone,
			Gender:     req.Gender,
			PreviewUrl: req.PreviewUrl,
			IsDefault:  false,
		},
	}, nil
}
