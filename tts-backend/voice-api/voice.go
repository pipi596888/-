package main

import (
	"flag"
	"fmt"
	"net/http"

	"tts-backend/voice-api/internal/config"
	"tts-backend/voice-api/internal/handler"
	"tts-backend/voice-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/voice-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	svcCtx := svc.NewServiceContext(&c)

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodGet,
				Path:    "/api/voice/list",
				Handler: handler.GetVoiceListHandler(svcCtx),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/voice/create",
				Handler: handler.CreateVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/voice/:id",
				Handler: handler.DeleteVoiceHandler(svcCtx),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/voice/default/:id",
				Handler: handler.SetDefaultVoiceHandler(svcCtx),
			},
		},
	)

	fmt.Printf("Starting voice-api at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
