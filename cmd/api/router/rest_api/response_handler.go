package rest_api

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

var (
	ErrMethodNotAllowed = &ErrorResponse{StatusCode: 405, Message: "Method not allowed"}
	ErrNotFound         = &ErrorResponse{StatusCode: 404, Message: "Resource not found"}
	ErrBadRequest       = &ErrorResponse{StatusCode: 400, Message: "Bad request"}
)

func RenderBadRequest(err error) *ErrorResponse {
	return ErrBadRequest
}
func ServerErrorRender(err error) *ErrorResponse {
	return &ErrorResponse{
		StatusCode: 500,
		Message:    err.Error(),
	}
}

func responseWithJson(writer http.ResponseWriter, status int, object interface{}) {
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(status)
	json.NewEncoder(writer).Encode(object)
}
