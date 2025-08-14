package main

import (
	"context"
	"lo/config"
	"lo/internal/http"
	"lo/internal/usecase"
	"log"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	cfg := config.NewConfig()

	taskTracker := usecase.NewTracker(ctx)

	httpErrChan := make(chan error)

	shutdownFunc, err := http.NewRouter(ctx, taskTracker, cfg.HostPort, httpErrChan)
	if err != nil {
		log.Fatalf("error: %s", err)
	}

	select {
	case <-ctx.Done():
		log.Println("app - Run - interrupted with signal")
		graceCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := shutdownFunc(graceCtx); err != nil {
			log.Printf("Error during graceful shutdown: %v", err)
		}

		return
	case err = <-httpErrChan:
		log.Printf("app - Run - httpErrChan: %v\n", err)
		stop()

		return
	}

}
