package utils

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type AudioMerger struct{}

func NewAudioMerger() *AudioMerger {
	return &AudioMerger{}
}

func (m *AudioMerger) MergeWavFiles(files [][]byte, outputFormat string) ([]byte, error) {
	if len(files) == 0 {
		return nil, fmt.Errorf("no files to merge")
	}

	if len(files) == 1 {
		return files[0], nil
	}

	tempDir := os.TempDir()
	inputList := filepath.Join(tempDir, "input_list.txt")

	inputFile, err := os.Create(inputList)
	if err != nil {
		return nil, err
	}

	tempFiles := make([]string, 0, len(files))
	for i, data := range files {
		tempFile := filepath.Join(tempDir, fmt.Sprintf("segment_%d.wav", i))
		err = os.WriteFile(tempFile, data, 0644)
		if err != nil {
			return nil, err
		}
		tempFiles = append(tempFiles, tempFile)
		inputFile.WriteString(fmt.Sprintf("file '%s'\n", tempFile))
	}
	inputFile.Close()

	outputFile := filepath.Join(tempDir, "merged.wav")

	cmd := exec.Command("ffmpeg", "-f", "concat", "-safe", "0", "-i", inputList, "-y", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg merge failed: %v, output: %s", err, string(output))
	}

	for _, f := range tempFiles {
		os.Remove(f)
	}
	os.Remove(inputList)

	result, err := os.ReadFile(outputFile)
	os.Remove(outputFile)
	if err != nil {
		return nil, err
	}

	if outputFormat == "mp3" {
		return m.ConvertWavToMp3(result)
	}

	return result, nil
}

func (m *AudioMerger) ConvertWavToMp3(wavData []byte) ([]byte, error) {
	tempDir := os.TempDir()
	inputFile := filepath.Join(tempDir, "convert_input.wav")
	outputFile := filepath.Join(tempDir, "convert_output.mp3")

	err := os.WriteFile(inputFile, wavData, 0644)
	if err != nil {
		return nil, err
	}
	defer os.Remove(inputFile)

	cmd := exec.Command("ffmpeg", "-i", inputFile, "-codec:a", "libmp3lame", "-qscale:a", "2", "-y", outputFile)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg convert failed: %v, output: %s", err, string(output))
	}
	defer os.Remove(outputFile)

	result, err := os.ReadFile(outputFile)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func MergeAudioChannels(left, right []byte) ([]byte, error) {
	return bytes.Join([][]byte{left, right}, nil), nil
}
