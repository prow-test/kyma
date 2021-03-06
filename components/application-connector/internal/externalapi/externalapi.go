package externalapi

import (
	"github.com/gorilla/mux"
	"net/http"
)

type MetadataHandler interface {
	CreateService(w http.ResponseWriter, r *http.Request)
	GetService(w http.ResponseWriter, r *http.Request)
	GetServices(w http.ResponseWriter, r *http.Request)
	UpdateService(w http.ResponseWriter, r *http.Request)
	DeleteService(w http.ResponseWriter, r *http.Request)
}

func NewHandler(handler MetadataHandler) http.Handler {
	router := mux.NewRouter()

	router.Path("/v1/health").Handler(NewHealthCheckHandler()).Methods(http.MethodGet)

	metadataRouter := router.PathPrefix("/{remoteEnvironment}/v1/metadata").Subrouter()
	metadataRouter.HandleFunc("/services", handler.CreateService).Methods(http.MethodPost)
	metadataRouter.HandleFunc("/services", handler.GetServices).Methods(http.MethodGet)
	metadataRouter.HandleFunc("/services/{serviceId}", handler.GetService).Methods(http.MethodGet)
	metadataRouter.HandleFunc("/services/{serviceId}", handler.UpdateService).Methods(http.MethodPut)
	metadataRouter.HandleFunc("/services/{serviceId}", handler.DeleteService).Methods(http.MethodDelete)

	router.NotFoundHandler = NewErrorHandler(404, "Requested resource could not be found.")
	router.MethodNotAllowedHandler = NewErrorHandler(405, "Method not allowed.")

	return router
}
