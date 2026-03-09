package logic

import (
	"context"

	"tts-backend/voice-api/internal/svc"
	"tts-backend/voice-api/internal/types"
)

type GetVoiceListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetVoiceListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetVoiceListLogic {
	return &GetVoiceListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetVoiceListLogic) GetVoiceList() (resp *types.VoiceListResp, err error) {
	voices, err := l.svcCtx.VoiceModel.FindAll()
	if err != nil {
		return nil, err
	}

	var list []types.Voice
	for _, v := range voices {
		list = append(list, types.Voice{
			Id:         v.Id,
			Name:       v.Name,
			Tone:       v.Tone,
			Gender:     v.Gender,
			PreviewUrl: v.PreviewUrl,
			IsDefault:  v.IsDefault,
		})
	}

	return &types.VoiceListResp{
		List:  list,
		Total: int64(len(list)),
	}, nil
}
