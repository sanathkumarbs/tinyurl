package tiny

import (
	"context"
	"io/fs"
	"log/slog"
	"time"
)

type DatabaseInterface interface {
	Migrate(files fs.FS, path string) error
}

type TinyURLRequest struct {
	Expiry   time.Time
	Original string
}

type TinyURLResponse struct {
	Expiry   time.Time
	Original string
	Tinyurl  string
}

type Service struct {
	db DatabaseInterface
}

func NewService(ctx context.Context, db DatabaseInterface) (*Service, error) {
	slog.Info("init TinyURL internal service")
	return &Service{
		db: db,
	}, nil
}

func (s *Service) CreateTinyURL(ctx context.Context, req TinyURLRequest) (TinyURLResponse, error) {
	slog.Info("processing a TinyURL request")
	return TinyURLResponse{
		Expiry:   req.Expiry,
		Original: req.Original,
		Tinyurl:  "sanathk.com/tinyurl123",
	}, nil
}
