package main

import (
	"flag"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"tts-backend/user-api/internal/config"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/user/login",
				Handler: loginHandler(&c),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/user/register",
				Handler: registerHandler(&c),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/works/list",
				Handler: getWorksHandler(&c),
			},
		},
	)

	fmt.Printf("Starting user-api at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
