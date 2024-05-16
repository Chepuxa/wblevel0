package customerrors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"wbtech/level0/internal/logger"

	"go.uber.org/zap"
)

type envelope map[string]interface{}

func LogError(r *http.Request, err error) {
	logger.Logger.Error("An error occurred",
		zap.Error(err),
		zap.String("request_method", r.Method),
		zap.String("request_url", r.URL.String()),
	)
}

func ErrorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := writeJson(w, status, env, nil, logger.Logger)
	if err != nil {
		LogError(r, err)
		w.WriteHeader(500)
	}
}

func ServerErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	logger.Logger.Error("The server encountered a problem and could not process the request", zap.Error(err))
	ErrorResponse(w, r, http.StatusInternalServerError, "the server encountered a problem and could not process your request")
}

func NotFoundResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusNotFound, "the requested resource could not be found")
}

func MethodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	ErrorResponse(w, r, http.StatusMethodNotAllowed, fmt.Sprintf("the %s method is not supported for this resource", r.Method))
}

func writeJson(w http.ResponseWriter, status int, data envelope, headers http.Header, logger *zap.SugaredLogger) error {
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		logger.Error("Error marshaling JSON", zap.Error(err))
	}

	js = append(js, '\n')

	for key, value := range headers {
		w.Header()[key] = value
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(js)
	if err != nil {
		logger.Error("Error writing JSON response", zap.Error(err))
	}

	return nil
}
