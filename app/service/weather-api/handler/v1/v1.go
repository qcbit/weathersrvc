// Package v1 contains the HTTP handlers for version 1 of the service.
package v1

import (
	"net/http"

	"go.uber.org/zap"

	"github.com/qcbit/weathersrvc/app/service/weather-api/handler/v1/public"
	"github.com/qcbit/weathersrvc/foundation/web"
)

type Config struct {
	Log *zap.SugaredLogger
}

// PublicRoutes adds all of the routes for the public API.
func PublicRoutes(app *web.App, cfg Config) {

	pbl := &public.Handler{
		Log: cfg.Log,
	}
	app.Handle(http.MethodGet, "/location/:coordinates", pbl.Coordinates)
}
