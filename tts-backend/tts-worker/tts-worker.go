package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"tts-backend/tts-worker/internal/config"
	"tts-backend/tts-worker/internal/worker"

	"github.com/zeromicro/go-zero/core/conf"
)

var configFile = flag.String("f", "etc/tts-worker.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	log.Printf("Starting TTS Worker...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	w := worker.NewTTSWorker(&c)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	if err := w.Start(ctx); err != nil {
		log.Printf("Worker error: %v", err)
	}

	log.Println("TTS Worker stopped")
}
