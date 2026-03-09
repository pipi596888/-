package logic

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"tts-backend/user-api/internal/svc"
	"tts-backend/user-api/internal/types"
	"tts-backend/user-api/internal/model"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	existingUser, _ := l.svcCtx.UserModel.FindByUsername(req.Username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Username:       req.Username,
		Password:       string(hashedPassword),
		Email:          req.Email,
		Balance:        0,
		CharacterCount: 0,
	}

	id, err := l.svcCtx.UserModel.Insert(user)
	if err != nil {
		return nil, err
	}

	token, err := l.generateToken(id, req.Username)
	if err != nil {
		return nil, err
	}

	return &types.RegisterResp{
		Token: token,
		User: types.UserInfoResp{
			Id:              id,
			Username:        req.Username,
			Balance:         0,
			CharacterCount:  0,
		},
	}, nil
}

func (l *RegisterLogic) generateToken(userId int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.JwtSecret))
}
