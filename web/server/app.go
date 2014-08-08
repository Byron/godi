// Package server implements a godi web server, hosting a thick client
package server

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

// Returns a handler suitable to provide a godi web frontend
func NewHandler() *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(
		&assetfs.AssetFS{
			Asset:    Asset,
			AssetDir: AssetDir,
			Prefix:   "",
		},
	))

	webHandler := NewWebSocketHandler()
	baseURL := "/api/v1/"
	socketURL := baseURL + "websocket"
	mux.Handle(socketURL, webHandler.handler())
	mux.Handle(baseURL+"state", NewRestHandler(webHandler.restStateHandler, socketURL))

	return mux
}
