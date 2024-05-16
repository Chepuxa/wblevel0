package api

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"wbtech/level0/internal/api/handlers"
	"wbtech/level0/internal/api/routes"
	"wbtech/level0/internal/config"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"
)

type API struct {
	Router *chi.Mux
	Config *config.Config
	Logger *zap.SugaredLogger
	Wg     *sync.WaitGroup
}

func NewAPI(log *zap.SugaredLogger, cfg *config.Config, h *handlers.Handlers, wg *sync.WaitGroup) *API {

	router := chi.NewRouter()

	router.Use(middleware.Logger)

	routes.SetupRoutes(router, h, cfg)

	return &API{
		Router: router,
		Logger: log,
		Config: cfg,
		Wg:     wg,
	}
}

func (a *API) Run() error {
	srv := &http.Server{
		Addr:         a.Config.Address,
		Handler:      a.Router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	shutdownError := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sign := <-quit

		a.Logger.Infow("Caught signal", "signal", sign.String())

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		if err != nil {
			shutdownError <- err
		}

		a.Logger.Infow("Finishing background tasks", "addr", srv.Addr)

		a.Wg.Wait()
		shutdownError <- nil
	}()

	a.Logger.Infow("Starting server", "addr", srv.Addr)

	err := srv.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	err = <-shutdownError
	if err != nil {
		return err
	}

	a.Logger.Infow("Stopped server", "addr", srv.Addr)

	return nil
}
