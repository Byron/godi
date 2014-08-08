package server

import (
	"encoding/json"
	"time"

	"github.com/Byron/godi/api"
	"github.com/Byron/godi/io"

	"code.google.com/p/go.net/websocket"
)

const (
	connectionTimeout = 1 * time.Minute
	writeTimeout      = 10 * time.Second
)

type MessageState uint8

const (
	StateChanged MessageState = iota
	StateResult
	StateFinished
)

// A JSON serializable struct to inform about what's going on
type jsonMessage struct {
	Message    string         `json:message`
	Error      string         `json:error`
	Importance api.Importance `json:importance`

	State MessageState `json:state`
}

func (m *jsonMessage) fromResult(r api.Result) {
	if r == nil {
		m.State = StateChanged
	} else {
		m.Message, m.Importance = r.Info()
		if r.Error() != nil {
			m.Error = r.Error().Error()
		}
		m.State = StateResult
	}
}

func (j *jsonMessage) String() string {
	res, _ := json.Marshal(j)
	return string(res)
}

type webSocketHandler struct {
	// The writer keeping our clients
	pmw *io.ParallelMultiWriter

	// Contains pre-formatted messages that we are to send to every client
	jsonBroadcasts chan string

	// This gates new clients
	newClients chan webClient
}

func NewWebSocketHandler() webSocketHandler {
	wsh := webSocketHandler{
		io.NewParallelMultiWriter(nil),
		// Yes, we block if too much is going on, which might slow down the entire operation
		make(chan string),
		make(chan webClient),
	}

	go wsh.eventHandler()

	return wsh
}

// To be run as go-loop to process our requests
// NOTE: currently, it never returns. It's fine, as connections are closed automatically
func (w *webSocketHandler) eventHandler() {
	for {
		select {
		case nc := <-w.newClients:
			{
				w.pmw.AutoInsert(&nc)
			}
		case msg := <-w.jsonBroadcasts:
			{
				// all we have to do is to write. Writers prune themselves if writing fails
				// We can't really check for errors as there is noone to tell.
				// Also we don't have a log right now ... .
				w.pmw.Write([]byte(msg))

				// Check for errors and nil the writers accordingly
				for wid := 0; wid < w.pmw.Capacity(); wid++ {
					wc, err := w.pmw.WriterAtIndex(wid)
					if err != nil {
						// NOTE: We should log this somewhere ...
						println("SERVER ERR", err.Error())
						wc.(*webClient).Close()
						w.pmw.SetWriterAtIndex(wid, nil)
					}
				}
			}
		} // channel select
	} // loop forever
}

// Returns a handler suited to listen on a particular web URL
func (w *webSocketHandler) handler() websocket.Handler {
	return websocket.Handler(w.handleConnection)
}

func (w *webSocketHandler) handleConnection(ws *websocket.Conn) {
	w.newClients <- webClient{ws}
}

// This one runs synchronously too
func (w *webSocketHandler) restStateHandler(isEnd bool, res api.Result) {
	m := jsonMessage{}
	if isEnd {
		m.State = StateFinished
	} else {
		m.fromResult(res)
	}

	w.jsonBroadcasts <- m.String()
}

// Implemnets a client which deals
type webClient struct {
	c *websocket.Conn // the connection to which to send something
}

func (w *webClient) Write(b []byte) (int, error) {
	w.c.SetWriteDeadline(time.Now().Add(writeTimeout))
	println("SERV PRE WRITE", len(b))
	n, err := w.c.Write(b)
	println("SERV WRITE DONE", n)

	// there is no need to check the error - the websocketHandler takes care
	// of dealing with is, closing the connection if needed
	w.c.SetWriteDeadline(time.Time{})
	return n, err
}

func (w *webClient) Close() error {
	return w.c.Close()
}
