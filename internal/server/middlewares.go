package server

import (
	"net/http"

	"github.com/AliRostami1/baagh/pkg/logy"
	"github.com/gorilla/mux"
)

func (s *Server) loggingMiddleware(log logy.Logger) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Do stuff here
			log.Infof("request received: %s", r.RequestURI)
			// Call the next handler, which can be another middleware in the chain, or the final handler.
			next.ServeHTTP(w, r)
		})
	}
}
