package types

type Segment struct {
	VoiceId int64  `json:"voiceId"`
	Emotion string `json:"emotion"`
	Text    string `json:"text"`
}

type GenerateReq struct {
	Segments []Segment `json:"segments"`
	Format   string    `json:"format"`
	Channel  string    `json:"channel"`
}

type GenerateResp struct {
	TaskId string `json:"taskId"`
}

type TaskResp struct {
	TaskId   string `json:"taskId"`
	Status   string `json:"status"`
	Progress int    `json:"progress"`
	AudioUrl string `json:"audioUrl,omitempty"`
	Error    string `json:"error,omitempty"`
}
