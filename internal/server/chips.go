package server

import (
	"fmt"
	"net/http"

	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/gorilla/mux"
	"github.com/warthog618/gpiod"
)

type ChipInfo struct {
	Name           string `json:"name"`
	Label          string `json:"label"`
	Lines          int    `json:"lines"`
	UapiAbiVersion int    `json:"uapiAbiVersion"`
}

func (s *Server) getAllChips(rw http.ResponseWriter, r *http.Request) {
	chips := []ChipInfo{}
	for _, chip := range core.Chips() {
		c, err := gpiod.NewChip(chip)
		if err != nil {
			s.serverError(rw, err)
			return
		}
		chips = append(chips, ChipInfo{
			Name:           c.Name,
			Label:          c.Label,
			Lines:          c.Lines(),
			UapiAbiVersion: c.UapiAbiVersion(),
		})
	}
	s.sendJSON(rw, chips)
}

func (s *Server) getOneChip(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	c, err := gpiod.NewChip(vars["chip"])
	if err != nil {
		s.clientError(rw, http.StatusNotFound, fmt.Sprintf("chip with name %s not found", vars["chip"]))
		return
	}
	s.sendJSON(rw, ChipInfo{
		Name:           c.Name,
		Label:          c.Label,
		Lines:          c.Lines(),
		UapiAbiVersion: c.UapiAbiVersion(),
	})
}
