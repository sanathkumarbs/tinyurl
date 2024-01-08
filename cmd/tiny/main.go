package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"sanathk.com/tinyurl/internal/tiny"
	"sanathk.com/tinyurl/internal/tiny/db"
	"sanathk.com/tinyurl/pkg/postgres"
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

	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			svc.Hello()
		case <-time.After(5 * time.Minute):
			ticker.Stop()
			return nil
		}
	}
}
