package engine

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type VITSEngine struct {
	modelPath string
}

func NewVITSEngine(modelPath string) *VITSEngine {
	return &VITSEngine{
		modelPath: modelPath,
	}
}

func (e *VITSEngine) Generate(text string, voiceId int64, emotion string) ([]byte, error) {
	tempDir := os.TempDir()
	inputFile := filepath.Join(tempDir, fmt.Sprintf("tts_input_%d.txt", voiceId))
	outputFile := filepath.Join(tempDir, fmt.Sprintf("tts_output_%d.wav", voiceId))

	err := os.WriteFile(inputFile, []byte(text), 0644)
	if err != nil {
		return nil, err
	}
	defer os.Remove(inputFile)

	cmd := exec.Command("python", "inference.py", "--text", text, "--output", outputFile)
	cmd.Dir = e.modelPath

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("vits inference failed: %v, output: %s", err, string(output))
	}
	defer os.Remove(outputFile)

	audioData, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, err
	}

	return audioData, nil
}

func (e *VITSEngine) GetVoiceName(voiceId int64) string {
	return fmt.Sprintf("voice_%d", voiceId)
}

type MockEngine struct{}

func NewMockEngine() *MockEngine {
	return &MockEngine{}
}

func (e *MockEngine) Generate(text string, voiceId int64, emotion string) ([]byte, error) {
	return []byte("mock audio data"), nil
}

func (e *MockEngine) GetVoiceName(voiceId int64) string {
	return fmt.Sprintf("voice_%d", voiceId)
}
