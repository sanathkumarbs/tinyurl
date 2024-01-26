package tiny

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	openapi_types "github.com/oapi-codegen/runtime/types"
	srvType "sanathk.com/tinyurl/pkg/api/services/v1/tiny"
)

type APIHandler struct {
	svc *Service
}

// TODO: potentially maintain this type on the OpenAPI Spec as it is for responding back to the API?
type APIError struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"` // in UTC
}

type Error struct {
	APIError
}

func (e Error) Error() string {
	return fmt.Sprintf("[%v] %s", e.Timestamp, e.Message)
}

func NewAPIHandler(ctx context.Context, svc *Service) (*APIHandler, error) {
	slog.Info("initializing a TinyURL API Handler")

	return &APIHandler{
		svc: svc,
	}, nil
}

func (a *APIHandler) CreateTinyURL(ec echo.Context) error {
	slog.Info("handling a CreateTinyURL API request")
	body := srvType.TinyURLRequest{}

	err := ec.Bind(&body)
	if err != nil {
		slog.Error("could not process request body to create a TinyURL", slog.Any("error", err.Error()))
		ae := APIError{
			Message:   "could not process request body to create a TinyURL",
			Timestamp: time.Now().UTC().String(),
		}
		return ec.JSON(
			http.StatusBadRequest, Error{ae},
		)
	}

	var resp TinyURLResponse
	req := TinyURLRequest{
		Expiry:   body.Expiry.Time,
		Original: body.Original,
	}
	resp, err = a.svc.CreateTinyURL(ec, req)
	if err != nil {
		slog.Error("could not complete creating TinyURL", slog.Any("error", err.Error()))
		ae := APIError{
			// What actionable action can we provide to the user here?
			// As this error message by itself is not very useful
			Message:   "could not complete creating TinyURL",
			Timestamp: time.Now().UTC().String(),
		}
		return ec.JSON(
			http.StatusInternalServerError, Error{ae},
		)
	}

	return ec.JSON(
		http.StatusOK, srvType.TinyURLResponse{
			Expiry:   openapi_types.Date{Time: resp.Expiry},
			Original: resp.Original,
			Tinyurl:  resp.Tinyurl,
		},
	)
}
