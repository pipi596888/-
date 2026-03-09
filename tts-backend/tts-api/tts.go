package main

import (
	"flag"
	"fmt"
	"net/http"

	"tts-backend/tts-api/internal/config"
	"tts-backend/tts-api/internal/handler"
	"tts-backend/tts-api/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/tts-api.yaml", "the config file")

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
				Method:  http.MethodPost,
				Path:    "/api/tts/generate",
				Handler: handler.GenerateHandler(svcCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/tts/task/:taskId",
				Handler: handler.QueryTaskHandler(svcCtx),
			},
		},
	)

	fmt.Printf("Starting tts-api at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
