package externalapi

import (
	"encoding/json"
	"github.com/kyma-project/kyma/components/application-connector/internal/httpconsts"
	"github.com/kyma-project/kyma/components/application-connector/internal/httperrors"
	"net/http"
)

type ErrorHandler struct {
	Message string
	Code    int
}

func NewErrorHandler(code int, message string) *ErrorHandler {
	return &ErrorHandler{Message: message, Code: code}
}

func (eh *ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	responseBody := httperrors.ErrorResponse{Code: eh.Code, Error: eh.Message}

	w.Header().Set(httpconsts.HeaderContentType, httpconsts.ContentTypeApplicationJson)
	w.WriteHeader(eh.Code)
	json.NewEncoder(w).Encode(responseBody)
}
