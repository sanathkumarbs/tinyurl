package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"sanathk.com/tinyurl/internal/tiny"
	"sanathk.com/tinyurl/internal/tiny/db"
	srvStub "sanathk.com/tinyurl/pkg/api/services/v1/tiny"
	"sanathk.com/tinyurl/pkg/apiserver"
	"sanathk.com/tinyurl/pkg/postgres"
)

const (
	// TODO: this should be definied and configured via env var/flags
	servicePort    = "8080"
	defaultAPIPath = "/api/v1"
)

func main() {
	if err := cmd(); err != nil {
		os.Exit(1)
	}
}

func cmd() error {
	rootCmd := &cobra.Command{
		Use:   "tiny",
		Short: "TinyURL Service",
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
	return rootCmd.Execute()
}

func run() error {
	slog.Info("running TinyURL service")
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	// TODO: remove hardcoded url, get from env vars and into a config
	rawURL := "postgres://postgres:admin@postgres:5432/tiny"
	pgURL, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	pg, err := postgres.NewPostgres(ctx, pgURL)
	if err != nil {
		return fmt.Errorf("creating postgres: %w", err)
	}

	err = pg.Migrate(db.Migrations, db.MigrationsPath)
	if err != nil {
		return err
	}

	svc, err := tiny.NewService(ctx, pg)
	if err != nil {
		return err
	}

	api, err := tiny.NewAPIHandler(ctx, svc)
	if err != nil {
		return err
	}

	srv, err := apiserver.NewEchoServer()
	if err != nil {
		return err
	}
	srv.Group(defaultAPIPath)
	srvStub.RegisterHandlersWithBaseURL(srv, api, defaultAPIPath)

	if err := srv.Start(fmt.Sprintf(":%s", servicePort)); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	return nil
}
