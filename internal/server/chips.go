package server

import (
	"net/http"

	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/warthog618/gpiod"
)

func (s *Server) getAllChips(rw http.ResponseWriter, r *http.Request) {
	chips := []gpiod.Chip{}
	for _, chip := range core.Chips() {
		c, err := gpiod.NewChip(chip)
		if err != nil {
			s.serverError(rw, err)
			return
		}
		chips = append(chips, *c)
	}
	s.sendJSON(rw, chips)
}
