package server

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
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

func (s *Server) clientError(rw http.ResponseWriter, status int, message string) {
	rw.WriteHeader(status)

	errRes, _ := json.Marshal(ErrorResponse{
		Success: false,
		Message: message,
	})

	rw.Write(errRes)
}

func (s *Server) notFoundError(rw http.ResponseWriter) {
	s.clientError(rw, http.StatusNotFound, http.StatusText(http.StatusNotFound))
}

type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
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

func (s *Server) wsServerError(ws *websocket.Conn, e error) {
	s.logger.Errorf("internal server error: %v", e)

	err := ws.WriteJSON(ErrorResponse{
		Success: false,
		Message: http.StatusText(http.StatusInternalServerError),
	})

	if marshalError, ok := err.(*json.MarshalerError); ok {
		s.logger.Errorf("marshaller error: %v", marshalError)
		return
	} else if err != nil {
		s.logger.Errorf("write json error: %v", err)
	}
}

func (s *Server) wsSendJSON(ws *websocket.Conn, data interface{}) {
	err := ws.WriteJSON(SuccessResponse{Success: true, Data: data})
	if err != nil {
		s.wsServerError(ws, err)
		return
	}
}
