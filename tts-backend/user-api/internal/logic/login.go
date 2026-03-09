package logic

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"tts-backend/user-api/internal/svc"
	"tts-backend/user-api/internal/types"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	user, err := l.svcCtx.UserModel.FindByUsername(req.Username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("password mismatch")
	}

	token, err := l.generateToken(user.Id, user.Username)
	if err != nil {
		return nil, err
	}

	return &types.LoginResp{
		Token: token,
		User: types.UserInfoResp{
			Id:              user.Id,
			Username:        user.Username,
			Balance:         user.Balance,
			CharacterCount:  user.CharacterCount,
		},
	}, nil
}

func (l *LoginLogic) generateToken(userId int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"userId":   userId,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(l.svcCtx.Config.JwtSecret))
}
