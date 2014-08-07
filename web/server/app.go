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

	mux.Handle("/api/v1/state", new(restHandler))

	return mux
}
