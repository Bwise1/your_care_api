package rest

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/bwise1/your_care_api/util"
	"github.com/bwise1/your_care_api/util/tracing"
	"github.com/bwise1/your_care_api/util/values"
)

type ServerResponse struct {
	Err        error           `json:"err,omitempty"`
	Message    string          `json:"message"`
	Status     string          `json:"status"`
	StatusCode int             `json:"status_code"`
	Context    context.Context `json:"context,omitempty"`
	Data       interface{}     `json:"data,omitempty"`
}

func respondWithJSONPayload(ctx *tracing.Context, data interface{}, status, message string) *ServerResponse {

	payload, err := json.Marshal(data)
	if err != nil {
		return respondWithError(err, "failed to marshal json payload", values.Error, ctx)
	}

	return &ServerResponse{
		Status:     values.Success,
		StatusCode: util.StatusCode(status),
		Message:    message,
		Data:       payload,
	}
}

// respondWithError logs an error with zap and parses the error to the ServerResponse
func respondWithError(err error, message, status string, tracingContext *tracing.Context) *ServerResponse {
	// logger.Log.Error(err.Error(), message, status)
	return &ServerResponse{
		// Err:        err,
		Message:    message,
		Status:     status,
		StatusCode: util.StatusCode(status),
	}
}

func writeJSONResponse(w http.ResponseWriter, content []byte, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if _, err := w.Write(content); err != nil {
		// logger.Log.Error("unable to write json response")
		log.Println("unable to write json response")
	}
}

// writeErrorResponse writes an error response to the client
func writeErrorResponse(w http.ResponseWriter, err error, status, errMessage string) {
	r := respondWithError(err, errMessage, status, nil)
	response, _ := json.Marshal(r)
	writeJSONResponse(w, response, r.StatusCode)
}
