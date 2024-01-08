package tiny

import (
	"context"
	"fmt"
	"io/fs"

	"github.com/jackc/pgx/v5"
	"sanathk.com/tinyurl/pkg/postgres"
)

type Postgreser interface {
	Conn() *pgx.Conn
	Migrate(files fs.FS, path string) error
}

type Service struct {
	pg Postgreser
}

// URL should include scheme://username:password@address:port/databasename
// example: "postgres://admin:admin@localhost:5432/testdb"
func NewService(ctx context.Context, pg postgres.Postgres) (*Service, error) {
	return &Service{
		pg: pg,
	}, nil
}

func (s *Service) Hello() {
	fmt.Println("hello")
}
