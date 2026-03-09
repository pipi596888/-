package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"tts-backend/tts-worker/internal/config"
	"tts-backend/tts-worker/internal/engine"
	"tts-backend/tts-worker/internal/model"
	"tts-backend/tts-worker/internal/utils"
)

type TTSWorker struct {
	config       *config.Config
	taskModel    model.TtsTaskModel
	segmentModel model.TtsSegmentModel
	engine       engine.TTSProvider
	merger       *utils.AudioMerger
}

type TaskMessage struct {
	TaskId string `json:"taskId"`
}

func NewTTSWorker(c *config.Config) *TTSWorker {
	db, err := sql.Open("mysql", c.Mysql.DataSource)
	if err != nil {
		panic(err)
	}

	var ttsEngine engine.TTSProvider
	if c.Aliyun.AccessKeyId != "" {
		ttsEngine = engine.NewMockEngine()
	} else {
		ttsEngine = engine.NewMockEngine()
	}

	return &TTSWorker{
		config:       c,
		taskModel:    model.NewTtsTaskModel(db),
		segmentModel: model.NewTtsSegmentModel(db),
		engine:       ttsEngine,
		merger:       utils.NewAudioMerger(),
	}
}

func (w *TTSWorker) ProcessTask(taskId string) error {
	log.Printf("Processing task: %s", taskId)

	err := w.taskModel.UpdateStatus(taskId, "processing", 0)
	if err != nil {
		return err
	}

	task, err := w.taskModel.FindByTaskId(taskId)
	if err != nil {
		w.taskModel.UpdateError(taskId, err.Error())
		return err
	}

	segments, err := w.segmentModel.FindByTaskId(taskId)
	if err != nil {
		w.taskModel.UpdateError(taskId, err.Error())
		return err
	}

	if len(segments) == 0 {
		w.taskModel.UpdateError(taskId, "no segments found")
		return fmt.Errorf("no segments found")
	}

	totalSegments := len(segments)
	audioDataList := make([][]byte, 0, totalSegments)

	for i, seg := range segments {
		log.Printf("Generating segment %d/%d: %s", i+1, totalSegments, seg.Text)

		audioData, err := w.engine.Generate(seg.Text, seg.VoiceId, seg.Emotion)
		if err != nil {
			w.taskModel.UpdateError(taskId, err.Error())
			return err
		}

		audioDataList = append(audioDataList, audioData)

		progress := (i + 1) * 80 / totalSegments
		w.taskModel.UpdateStatus(taskId, "processing", progress)
	}

	mergedAudio, err := w.merger.MergeWavFiles(audioDataList, task.Format)
	if err != nil {
		w.taskModel.UpdateError(taskId, err.Error())
		return err
	}

	audioUrl := fmt.Sprintf("https://%s.%s/tts/%s.%s",
		w.config.Oss.BucketName,
		w.config.Oss.Endpoint,
		taskId,
		task.Format,
	)

	log.Printf("Task %s completed, audio URL: %s", taskId, audioUrl)

	w.taskModel.UpdateAudioUrl(taskId, audioUrl)

	return nil
}

func (w *TTSWorker) Start(ctx context.Context) error {
	log.Println("TTS Worker started")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker stopping")
			return nil
		case <-ticker.C:
			w.processPendingTasks()
		}
	}
}

func (w *TTSWorker) processPendingTasks() {
	log.Println("Checking for pending tasks...")
}

func HandleTaskMessage(data []byte) error {
	var msg TaskMessage
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return err
	}

	fmt.Printf("Received task: %s\n", msg.TaskId)
	return nil
}
