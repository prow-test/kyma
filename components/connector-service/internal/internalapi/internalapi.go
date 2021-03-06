package internalapi

import (
	"github.com/gorilla/mux"
	"github.com/kyma-project/kyma/components/connector-service/internal/errorhandler"
	"net/http"
)

type TokenHandler interface {
	CreateToken(w http.ResponseWriter, r *http.Request)
}

func NewHandler(handler TokenHandler) http.Handler {
	router := mux.NewRouter()

	tokenRouter := router.PathPrefix("/v1/remoteenvironments").Subrouter()

	tokenRouter.HandleFunc("/{reName}/tokens", handler.CreateToken).Methods(http.MethodPost)

	router.NotFoundHandler = errorhandler.NewErrorHandler(404, "Requested resource could not be found.")
	router.MethodNotAllowedHandler = errorhandler.NewErrorHandler(405, "Method not allowed.")

	return router
}
