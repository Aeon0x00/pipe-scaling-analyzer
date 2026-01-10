package httpserver

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Server struct {
	httpServer *http.Server
	handlers   *HTTPHandlers
}

func NewServer (port string, handlers *HTTPHandlers) *Server {
	router := mux.NewRouter()

	router.Path("/api/pipes").Methods("POST").HandlerFunc(handlers.HandleCreatePipe)
	router.Path("/api/pipes").Methods("GET").HandlerFunc(handlers.HandleListPipes)
	
	router.Path("/api/pipes/{id}").Methods("GET").HandlerFunc(handlers.HandleGetPipe)
	router.Path("/api/pipes/{id}").Methods("DELETE").HandlerFunc(handlers.HandleDeletePipe)


	router.Path("/api/calc").Methods("POST").HandlerFunc(handlers.HandleCalculateScaling)
	router.Path("/api/pipes/{id}/status").Methods("GET").HandlerFunc(handlers.HandleGetPipeStatus)
	
	router.Path("/api/pipes/{id}/history").Methods("GET").HandlerFunc(handlers.HandleGetPipeHistory)

	router.Path("/api/alerts").Methods("GET").HandlerFunc(handlers.HandleListAlerts)
	router.Path("/api/alerts/{id}/fix").Methods("PATCH").HandlerFunc(handlers.HandleFixAlert)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return &Server{
		httpServer: srv,
		handlers:   handlers,
	}
}
func (s *Server) Start() error {
	if err := s.httpServer.ListenAndServe(); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			
			return nil
		}
		return err
	}
	return  nil
}
func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
