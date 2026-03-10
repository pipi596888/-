package handler

import (
	"net/http"

	"tts-backend/voice-api/internal/auth"
	"tts-backend/voice-api/internal/logic"
	"tts-backend/voice-api/internal/svc"
	"tts-backend/voice-api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetVoiceListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewGetVoiceListLogic(r.Context(), svcCtx)
		userId, isAdmin, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			userId = 0
			isAdmin = false
		}
		resp, err := l.GetVoiceList(userId, isAdmin)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func CreateVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateVoiceReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCreateVoiceLogic(r.Context(), svcCtx)
		resp, err := l.CreateVoice(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}

func DeleteVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, ok := parseIDWithPrefix(r.URL.Path, "/api/voice/")
		if !ok {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "id is required"})
			return
		}

		l := logic.NewDeleteVoiceLogic(r.Context(), svcCtx)
		err := l.DeleteVoice(&types.DeleteVoiceReq{Id: id})
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}

func SetDefaultVoiceHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, _, err := auth.ParseUserIDFromRequest(r, svcCtx.Config.JwtSecret)
		if err != nil {
			httpx.WriteJson(w, http.StatusUnauthorized, map[string]interface{}{"code": 401, "message": "unauthorized"})
			return
		}

		id, ok := parseIDWithPrefix(r.URL.Path, "/api/voice/default/")
		if !ok {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "id is required"})
			return
		}

		if _, err := svcCtx.VoiceModel.FindOne(id); err != nil {
			httpx.WriteJson(w, http.StatusBadRequest, map[string]interface{}{"code": 400, "message": "voice not found"})
			return
		}

		l := logic.NewSetDefaultVoiceLogic(r.Context(), svcCtx)
		err = l.SetDefaultVoice(userId, &types.SetDefaultReq{Id: id})
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.Ok(w)
		}
	}
}
