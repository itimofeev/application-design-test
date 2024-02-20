package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/kelseyhightower/envconfig"
	"golang.org/x/sync/errgroup"

	"applicationDesignTest/internal/app/booking"
	"applicationDesignTest/internal/repository/inmemory"
	"applicationDesignTest/internal/server"
)

type config struct {
	Address string `envconfig:"APP_ADDRESS" default:":8080"`
}

func main() {
	config := config{}
	if err := envconfig.Process("", &config); err != nil {
		panic(err)
	}

	ctx := signalContext(context.Background())

	if err := run(ctx, config); err != nil {
		slog.Error("error on running main", "error", err)
		panic(fmt.Sprintf("error on running main: %s", err))
	}
	slog.Info("app finished")
}

func run(ctx context.Context, cfg config) error {
	errGr, errGrCtx := errgroup.WithContext(ctx)

	repository, err := inmemory.New(inmemory.Config{})
	if err != nil {
		return fmt.Errorf("failed to create repository: %w", err)
	}

	bookingApp, err := booking.New(booking.Config{
		Repository: repository,
	})
	if err != nil {
		return fmt.Errorf("failed to create booking app: %w", err)
	}

	srv, err := server.NewServer(server.Config{
		Address:    cfg.Address,
		BookingApp: bookingApp,
	})
	if err != nil {
		return err
	}

	errGr.Go(func() error {
		slog.InfoContext(ctx, "start http server")

		return srv.Serve(errGrCtx)
	})

	// other workers or processes can be run via errgroup here

	return errGr.Wait()
}

// signalContext returns a context that is canceled if either SIGTERM or SIGINT signal is received.
func signalContext(ctx context.Context) context.Context {
	cnCtx, cancel := context.WithCancel(ctx)

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		sig := <-c
		slog.InfoContext(cnCtx, "received signal", "signal", sig)
		cancel()
	}()

	return cnCtx
}
