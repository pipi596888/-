package svc

import (
	"database/sql"

	"tts-backend/user-api/internal/config"
	"tts-backend/user-api/internal/model"
)

type ServiceContext struct {
	Config    *config.Config
	UserModel model.UserModel
}

func NewServiceContext(c *config.Config) *ServiceContext {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	
	return &ServiceContext{
		Config:    c,
		UserModel: model.NewUserModel(db),
	}
}
