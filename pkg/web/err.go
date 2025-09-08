package web

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/marcelofabianov/fault"
)

type ErrorResponse struct {
	Error   string         `json:"error"`
	Details any            `json:"details,omitempty"`
	Context map[string]any `json:"context,omitempty"`
}

func ErrorResponder(w http.ResponseWriter, logger *slog.Logger, err error) {
	var f *fault.Error
	if errors.As(err, &f) {
		statusCode := mapFaultCodeToHTTPStatus(f.Code)
		response := ErrorResponse{
			Error:   f.Message,
			Context: f.Context,
		}
		RespondJSON(w, statusCode, response, logger)
		return
	}

	logger.Error("unhandled internal error", "error", err)
	response := ErrorResponse{Error: "an unexpected internal error occurred"}
	RespondJSON(w, http.StatusInternalServerError, response, logger)
}

func mapFaultCodeToHTTPStatus(code fault.Code) int {
	switch code {
	case fault.Conflict:
		return http.StatusConflict
	case fault.Invalid, fault.DomainViolation:
		return http.StatusBadRequest
	case fault.NotFound:
		return http.StatusNotFound
	case fault.Unauthorized:
		return http.StatusUnauthorized
	case fault.Forbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
