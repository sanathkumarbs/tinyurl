package postgres

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

type Postgres struct {
	connString *url.URL
	conn       *pgx.Conn
}

// URL should include scheme://username:password@address:port/databasename
// example: "postgres://admin:admin@localhost:5432/testdb"
func NewPostgres(ctx context.Context, connString *url.URL) (Postgres, error) {
	pg := Postgres{
		connString: connString,
	}
	conn, err := pg.connect(ctx)
	if err != nil {
		return pg, fmt.Errorf("connecting to postgres: %w", err)
	}
	pg.conn = conn

	return pg, nil
}

func (pg Postgres) connect(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, pg.connString.String())
}

func (pg Postgres) Conn() *pgx.Conn {
	return pg.conn
}

func (pg Postgres) Migrate(files fs.FS, path string) error {
	c, err := pgx.ParseConfig(pg.connString.String())
	if err != nil {
		return fmt.Errorf("parsing postgres connString: %w", err)
	}

	sourceInstance, err := iofs.New(files, path)
	if err != nil {
		return fmt.Errorf("setting up sourceInstance: %w", err)
	}

	dbName := pg.connString.Path
	pc := new(postgres.Config)
	pc.DatabaseName = dbName

	databaseInstance, err := postgres.WithInstance(
		stdlib.OpenDB(*c),
		pc,
	)
	if err != nil {
		return fmt.Errorf("setting up databaseInstance: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", sourceInstance, dbName, databaseInstance)
	if err != nil {
		return fmt.Errorf("setting up migrateInstance: %w", err)
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return fmt.Errorf("migrating up: %w", err)
	}
	return nil
}
