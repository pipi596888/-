package logic

import (
	"context"

	"tts-backend/voice-api/internal/svc"
	"tts-backend/voice-api/internal/types"
)

type DeleteVoiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteVoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteVoiceLogic {
	return &DeleteVoiceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteVoiceLogic) DeleteVoice(req *types.DeleteVoiceReq) error {
	return l.svcCtx.VoiceModel.Delete(req.Id)
}

type SetDefaultVoiceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSetDefaultVoiceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SetDefaultVoiceLogic {
	return &SetDefaultVoiceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SetDefaultVoiceLogic) SetDefaultVoice(req *types.SetDefaultReq) error {
	return l.svcCtx.VoiceModel.SetDefault(req.Id)
}
