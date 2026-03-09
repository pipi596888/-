package svc

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"tts-backend/voice-api/internal/config"
	"tts-backend/voice-api/internal/model"
)

type ServiceContext struct {
	Config    *config.Config
	VoiceModel model.VoiceModel
}

func NewServiceContext(c *config.Config) *ServiceContext {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:    c,
		VoiceModel: model.NewVoiceModel(db),
	}
}
