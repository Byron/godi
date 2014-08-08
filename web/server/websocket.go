package server

import (
	"code.google.com/p/go.net/websocket"

	"github.com/Byron/godi/api"
)

type webSocketHandler struct {
}

// Returns a handler suited to listen on a particular web URL
func (w *webSocketHandler) handler() websocket.Handler {
	return websocket.Handler(w.handleConnection)
}

func (w *webSocketHandler) handleConnection(ws *websocket.Conn) {

}

func (w *webSocketHandler) restStateHandler(isEnd bool, res api.Result) {

	// Perform a braodcast to all clients, in parallel, timeout per client
	if res != nil {
		msg, prio := res.Info()
		println(isEnd, msg, prio)
	} else {
		println(isEnd)
	}

}
