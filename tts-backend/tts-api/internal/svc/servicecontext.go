package svc

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"tts-backend/tts-api/internal/config"
	"tts-backend/tts-api/internal/model"
)

type ServiceContext struct {
	Config       *config.Config
	TaskModel    model.TtsTaskModel
	SegmentModel model.TtsSegmentModel
	VoiceAccess  model.VoiceAccessModel
}

func NewServiceContext(c *config.Config) *ServiceContext {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	return &ServiceContext{
		Config:       c,
		TaskModel:    model.NewTtsTaskModel(db),
		SegmentModel: model.NewTtsSegmentModel(db),
		VoiceAccess:  model.NewVoiceAccessModel(db),
	}
}
