package apiserver

import (
	"log/slog"

	"github.com/labstack/echo/v4"
)

type EchoServer struct {
	*echo.Echo
}

func NewEchoServer() (*EchoServer, error) {
	slog.Info("creating an echo server")
	return &EchoServer{
		echo.New(),
	}, nil
}
