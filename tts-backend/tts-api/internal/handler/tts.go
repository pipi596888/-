package handler

import (
	"errors"
	"net/http"
	"strings"

	"tts-backend/tts-api/internal/auth"
	"tts-backend/tts-api/internal/logic"
	"tts-backend/tts-api/internal/svc"
	"tts-backend/tts-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GenerateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
			return
		}

		var req types.GenerateReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewGenerateLogic(r.Context(), svcCtx)
		resp, err := l.Generate(&req, userId, isAdmin)
		if err != nil {
			if errors.Is(err, logic.ErrForbidden) {
				httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{
					"code":    http.StatusForbidden,
					"message": "forbidden",
				})
				return
			}
			if errors.Is(err, logic.ErrInvalidVoice) {
				httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{
					"code":    http.StatusBadRequest,
					"message": "invalid voice",
				})
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func QueryTaskHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
			return
		}

		taskId := strings.TrimPrefix(r.URL.Path, "/api/tts/task/")

		l := logic.NewQueryTaskLogic(r.Context(), svcCtx)
		resp, err := l.QueryTask(taskId, userId, isAdmin)
		if err != nil {
			if errors.Is(err, logic.ErrForbidden) {
				httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{
					"code":    http.StatusForbidden,
					"message": "forbidden",
				})
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func QueryTaskDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{
				"code":    http.StatusUnauthorized,
				"message": "unauthorized",
			})
			return
		}

		taskId := strings.TrimPrefix(r.URL.Path, "/api/tts/task/")
		taskId = strings.TrimSuffix(taskId, "/detail")
		taskId = strings.Trim(taskId, "/")

		l := logic.NewQueryTaskDetailLogic(r.Context(), svcCtx)
		resp, err := l.QueryTaskDetail(taskId, userId, isAdmin)
		if err != nil {
			if errors.Is(err, logic.ErrForbidden) {
				httpx.WriteJson(w, http.StatusForbidden, map[string]interface{}{
					"code":    http.StatusForbidden,
					"message": "forbidden",
				})
				return
			}
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
