package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/gorilla/mux"
)

type ItemGet struct {
}

func (s *Server) getOneItem(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chipName, itemOffset := vars["chip"], vars["offset"]
	ioInt, err := strconv.Atoi(itemOffset)
	if err != nil {
		// this error should never happen as the route only
		// cathces offsets that are integers, but we have it
		// here anyways just in case
		s.clientError(rw, http.StatusBadRequest, "offset should be integer")
		return
	}

	item, err := core.GetItem(chipName, ioInt)
	if err != nil {
		s.clientError(rw, http.StatusBadRequest, fmt.Sprintf("item with %s offset isn't registered on chip %s", itemOffset, chipName))
		return
	}

	itemInfo, err := item.Info()
	if err != nil {
		s.serverError(rw, err)
		return
	}

	s.sendJSON(rw, itemInfo)
}

type createItem struct {
	Mode  string `json:"mode"`
	Pull  string `json:"pull"`
	State string `json:"state"`
}

func (s *Server) createOneItem(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chip, itemOffset := vars["chip"], vars["offset"]
	ioInt, err := strconv.Atoi(itemOffset)
	if err != nil {
		// this error should never happen as the route only
		// cathces offsets that are integers, but we have it
		// here anyways just in case
		s.clientError(rw, http.StatusBadRequest, "offset should be integer")
		return
	}

	var createItemData createItem
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&createItemData)

	if err != nil {
		s.serverError(rw, err)
		return
	}
	defer r.Body.Close()

	item, err := core.RequestItemHelper(chip, ioInt, createItemData.Mode, createItemData.Pull, createItemData.State)
	if err != nil {
		s.clientError(rw, http.StatusBadRequest, fmt.Sprintf("failed requesting the item %v", err))
		return
	}

	itemInfo, err := item.Info()
	if err != nil {
		s.serverError(rw, err)
		item.Close()
		return
	}

	s.sendJSON(rw, itemInfo)
}

func (s *Server) deleteOneItem(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	chip, itemOffset := vars["chip"], vars["offset"]
	ioInt, err := strconv.Atoi(itemOffset)
	if err != nil {
		// this error should never happen as the route only
		// cathces offsets that are integers, but we have it
		// here anyways just in case
		s.clientError(rw, http.StatusBadRequest, "offset should be integer")
		return
	}

	item, err := core.GetItem(chip, ioInt)
	if err != nil {
		s.clientError(rw, http.StatusBadRequest, "item not registered")
		return
	}

	err = item.Close()
	if err != nil {
		s.serverError(rw, err)
	}

	itemInfo, err := item.Info()
	if err != nil {
		s.sendJSON(rw, "item got fully cleaned up")
		return

	}

	s.sendJSON(rw, itemInfo)
}

func (s *Server) watchOneItem(rw http.ResponseWriter, r *http.Request) {
	// core.GetItem()
}

func (s *Server) getAllItems(rw http.ResponseWriter, r *http.Request) {
	chip := mux.Vars(r)["chip"]

	info, err := core.Info()
	if err != nil {
		s.serverError(rw, err)
		return
	}

	items, ok := info[chip]
	if !ok {
		s.clientError(rw, http.StatusBadRequest, "chip not found")
		return
	}

	s.sendJSON(rw, items)

}

func (s *Server) watchAllItems(rw http.ResponseWriter, r *http.Request) {
	// core.GetItem()
}
