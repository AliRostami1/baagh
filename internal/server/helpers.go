package server

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func (s *Server) serverError(rw http.ResponseWriter, e error) {
	s.logger.Errorf("internal server error: %v", e)

	rw.WriteHeader(http.StatusInternalServerError)

	errRes, _ := json.Marshal(ErrorResponse{
		Success: false,
		Message: http.StatusText(http.StatusInternalServerError),
	})

	rw.Write(errRes)
}

func (s *Server) clientError(rw http.ResponseWriter, status int) {
	rw.WriteHeader(status)

	errRes, _ := json.Marshal(ErrorResponse{
		Success: false,
		Message: http.StatusText(status),
	})

	rw.Write(errRes)
}

func (s *Server) notFoundError(rw http.ResponseWriter) {
	s.clientError(rw, http.StatusNotFound)
}

type SuccessResponse struct {
	Success bool
	Data    interface{}
}

func (s *Server) sendJSON(rw http.ResponseWriter, data interface{}) {
	sucRes, err := json.Marshal(SuccessResponse{Success: true, Data: data})
	if err != nil {
		s.serverError(rw, err)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(sucRes)
}
