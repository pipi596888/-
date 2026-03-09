package engine

type TTSProvider interface {
	Generate(text string, voiceId int64, emotion string) ([]byte, error)
	GetVoiceName(voiceId int64) string
}
