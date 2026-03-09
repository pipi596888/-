package types

type Voice struct {
	Id         int64  `json:"id"`
	Name       string `json:"name"`
	Tone       string `json:"tone"`
	Gender     string `json:"gender"`
	PreviewUrl string `json:"previewUrl,omitempty"`
	IsDefault  bool   `json:"isDefault"`
}

type VoiceListReq struct{}

type VoiceListResp struct {
	List  []Voice `json:"list"`
	Total int64   `json:"total"`
}

type CreateVoiceReq struct {
	Name       string `json:"name"`
	Tone       string `json:"tone"`
	Gender     string `json:"gender"`
	PreviewUrl string `json:"previewUrl"`
}

type CreateVoiceResp struct {
	Voice Voice `json:"voice"`
}

type DeleteVoiceReq struct {
	Id int64 `json:"id"`
}

type SetDefaultReq struct {
	Id int64 `json:"id"`
}
