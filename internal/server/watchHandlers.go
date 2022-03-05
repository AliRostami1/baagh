package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AliRostami1/baagh/pkg/controller/core"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
)

// Time allowed to read the next pong message from the client.
const pongWait = 60 * time.Second

var upgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
}

func reader(ws *websocket.Conn) {
	defer ws.Close()
	ws.SetReadLimit(512)
	ws.SetReadDeadline(time.Now().Add(pongWait))
	ws.SetPongHandler(func(string) error {
		ws.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, _, err := ws.ReadMessage()
		if err != nil {
			break
		}
	}
}

func (s *Server) watchOneItem(rw http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(rw, r, nil)
	if err != nil {
		s.logger.Errorf("upgrade:", err)
		return
	}
	defer ws.Close()

	vars := mux.Vars(r)
	chip, itemOffset := vars["chip"], vars["offset"]
	ioInt, err := strconv.Atoi(itemOffset)
	if err != nil {
		// this error should never happen as the route only
		// cathces offsets that are integers, but we have it
		// here anyways just in case
		// s.clientError(rw, http.StatusBadRequest, "offset should be integer")
		s.wsServerError(ws, err)
		return
	}

	watcher, err := core.NewWatcher(chip, ioInt)
	if err != nil {
		s.wsServerError(ws, err)
		return
	}
	defer watcher.Close()

	item, err := core.GetItem(chip, ioInt)
	if err != nil {
		s.wsServerError(ws, err)
		return
	}

	itemInfo, err := item.Info()
	if err != nil {
		s.wsServerError(ws, err)
		return
	}

	s.wsSendJSON(ws, itemInfo)

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case ev, ok := <-watcher.Watch():
				if !ok {
					return
				}
				s.wsSendJSON(ws, ev.Info)
			}
		}
	}()

	reader(ws)
}
