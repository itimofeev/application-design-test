package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-playground/validator/v10"

	"applicationDesignTest/internal/app/booking"
)

type Config struct {
	Address    string `validate:"required"`
	BookingApp *booking.App
}

type Server struct {
	srv        *http.Server
	bookingApp *booking.App
	handler    handler
}

func NewServer(cfg Config) (*Server, error) {
	err := validator.New().Struct(cfg)
	if err != nil {
		return nil, err
	}

	s := &Server{
		srv: &http.Server{
			Addr: cfg.Address,
		},
		handler: handler{
			bookingApp: cfg.BookingApp,
		},
	}

	s.init(cfg)

	return s, nil
}

func (s *Server) Serve(ctx context.Context) error {
	go func() {
		<-ctx.Done()
		shutdownTimeout := time.Second * 5 // todo pass via config
		ctxShutdown, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		if err := s.srv.Shutdown(ctxShutdown); err != nil {
			slog.WarnContext(ctx, "failed to shutdown server", "error", err)
		}
	}()

	if err := s.srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) init(cfg Config) {
	r := chi.NewRouter()
	r.Use(
		middleware.Recoverer,
	)

	r.Group(func(router chi.Router) {
		// metrics, version, pprof, etc. handlers
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.DefaultLogger)

		r.Route("/api/v1", func(api chi.Router) {
			api.Post("/orders", s.handler.createOrder)
		})
	})

	s.srv.Handler = r
}
