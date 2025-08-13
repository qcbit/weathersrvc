// Package handler handles the HTTP request and response for the service.
package handler

import (
	"net/http"
	"os"

	"go.uber.org/zap"

	v1 "github.com/qcbit/weathersrvc/app/service/weather-api/handler/v1"
	"github.com/qcbit/weathersrvc/foundation/web"
)

// MuxConfig contains the configuration for the HTTP server.
type MuxConfig struct {
	Shutdown chan os.Signal
	Log      *zap.SugaredLogger
}

// PublicMux constructs a http.Handler with the application routes.
func PublicMux(cfg MuxConfig) http.Handler {
	// Construct the web app.
	app := web.NewApp(cfg.Shutdown)

	// Load routes.
	v1.PublicRoutes(app, v1.Config{
		Log: cfg.Log,
	})
	return app
}
