// This is the main package for the weather API service.
package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"go.uber.org/zap"

	"github.com/qcbit/weathersrvc/app/service/weather-api/handler"
	"github.com/qcbit/weathersrvc/foundation/logger"
)

// main is the entry point for the weather API service.
func main() {
	log, err := logger.New("WEATHER-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

func run(log *zap.SugaredLogger) error {
	// --------------------------------------------------------------
	// GOMAXPROCS
	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// ---------------------------------------------------------------
	// Web server startup
	log.Infow("startup", "status", "starting web server")
	publicHost := ":8080"

	// ---------------------------------------------------------------
	// Signal handling for graceful shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverErrors := make(chan error, 1)

	publicMux := handler.PublicMux(handler.MuxConfig{
		Shutdown: shutdown,
		Log:      log,
	})

	public := http.Server{
		Addr:    publicHost,
		Handler: publicMux,
	}

	// Start the public API server in a goroutine; Serve and listen
	go func() {
		log.Infow("startup", "status", "public API router started", "host", public.Addr)
		serverErrors <- public.ListenAndServe()
	}()

	// Wait for shutdown signal
	<-shutdown
	log.Infow("shutdown", "status", "shutdown started", "signal")
	defer log.Infow("shutdown", "status", "shutdown complete", "signal")

	return nil
}
