package logic

import (
	"context"
	"errors"

	"tts-backend/user-api/internal/svc"
	"tts-backend/user-api/internal/types"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.UserInfoResp, err error) {
	userId := l.ctx.Value("userId")
	if userId == nil {
		return nil, errors.New("unauthorized")
	}

	user, err := l.svcCtx.UserModel.FindOne(userId.(int64))
	if err != nil {
		return nil, err
	}

	return &types.UserInfoResp{
		Id:              user.Id,
		Username:        user.Username,
		Balance:         user.Balance,
		CharacterCount:  user.CharacterCount,
	}, nil
}
