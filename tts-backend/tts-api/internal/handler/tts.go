package handler

import (
	"net/http"
	"strings"

	"tts-backend/tts-api/internal/logic"
	"tts-backend/tts-api/internal/svc"
	"tts-backend/tts-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GenerateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGenerateLogic(r.Context(), svcCtx)
		resp, err := l.Generate(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func QueryTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		taskId := strings.TrimPrefix(r.URL.Path, "/api/tts/task/")

		l := logic.NewQueryTaskLogic(r.Context(), svcCtx)
		resp, err := l.QueryTask(taskId)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
