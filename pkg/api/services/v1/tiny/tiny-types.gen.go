// Package tiny provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen/v2 version v2.0.0 DO NOT EDIT.
package tiny

import (
	openapi_types "github.com/oapi-codegen/runtime/types"
)

// TinyURLRequest defines model for TinyURLRequest.
type TinyURLRequest struct {
	Expiry   *openapi_types.Date `json:"expiry,omitempty"`
	Original string              `json:"original"`
}

// TinyURLResponse defines model for TinyURLResponse.
type TinyURLResponse struct {
	Expiry   openapi_types.Date `json:"expiry"`
	Original string             `json:"original"`
	Tinyurl  string             `json:"tinyurl"`
}

// TinyURLBody defines model for TinyURLBody.
type TinyURLBody = TinyURLRequest

// CreateTinyURLJSONRequestBody defines body for CreateTinyURL for application/json ContentType.
type CreateTinyURLJSONRequestBody = TinyURLRequest
