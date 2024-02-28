package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Hymiside/hezzl-api/pkg/models"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

func StartServer(ctx context.Context, handler *chi.Mux, cfg models.ConfigServer) error {
	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Handler: handler,
	}

	go func() {
		<-ctx.Done()
		if err := httpServer.Shutdown(ctx); err != nil {
			log.Fatalf("failed to shutdown server: %v", err)
		}
	}()

	log.Infof("server started on http://%s:%s/", cfg.Host, cfg.Port)
	return httpServer.ListenAndServe()
}