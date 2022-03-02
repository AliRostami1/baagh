package server

import (
	"net/http"
)

func (s *Server) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// Do stuff here
		s.logger.Infof("%s request received on: %s%s", r.Method, r.Host, r.URL.Path)
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(rw, r)
	})
}
