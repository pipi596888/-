package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"tts-backend/user-api/internal/config"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if err := ensureWorksSchema(db); err != nil {
		panic(err)
	}
	if err := ensureAdminSchema(db); err != nil {
		panic(err)
	}

	server.AddRoutes(
		[]rest.Route{
			{
				Method:  http.MethodPost,
				Path:    "/api/user/login",
				Handler: loginHandler(&c, db),
			},
			{
				Method: http.MethodGet,
				Path:   "/api/user/public-key",
				Handler: func(w http.ResponseWriter, r *http.Request) {
					resp, err := getPublicKey()
					if err != nil {
						httpx.ErrorCtx(r.Context(), w, err)
						return
					}
					httpx.OkJsonCtx(r.Context(), w, resp)
				},
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/user/register",
				Handler: registerHandler(&c, db),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/user/info",
				Handler: requireAuth(&c, getUserInfoHandler(db)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/works/list",
				Handler: requireAuth(&c, getWorksHandler(db)),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/works/:taskId/title",
				Handler: requireAuth(&c, updateWorkTitleHandler(db)),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/works/:taskId",
				Handler: requireAuth(&c, deleteWorkHandler(db)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/feedback",
				Handler: requireAuth(&c, createFeedbackHandler(db)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/feedback/my",
				Handler: requireAuth(&c, listMyFeedbackHandler(db)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/admin/feedback/list",
				Handler: requireAdmin(&c, db, listAllFeedbackHandler(db)),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/admin/feedback/:id/reply",
				Handler: requireAdmin(&c, db, replyFeedbackHandler(db)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/admin/system/stats",
				Handler: requireAdmin(&c, db, systemStatsHandler(db)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/admin/users",
				Handler: requireAdmin(&c, db, listAdminUsersHandler(db)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/admin/users",
				Handler: requireAdmin(&c, db, createAdminUserHandler(db)),
			},
			{
				Method:  http.MethodPut,
				Path:    "/api/admin/users/:id",
				Handler: requireAdmin(&c, db, updateAdminUserHandler(db)),
			},
			{
				Method:  http.MethodDelete,
				Path:    "/api/admin/users/:id",
				Handler: requireAdmin(&c, db, deleteAdminUserHandler(db)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/admin/roles",
				Handler: requireAdmin(&c, db, listAdminRolesHandler(db)),
			},
			{
				Method:  http.MethodGet,
				Path:    "/api/admin/logs",
				Handler: requireAdmin(&c, db, listAdminLogsHandler(db)),
			},
			{
				Method:  http.MethodPost,
				Path:    "/api/admin/logs",
				Handler: requireAdmin(&c, db, appendAdminLogHandler(db)),
			},
		},
	)

	fmt.Printf("Starting user-api at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
