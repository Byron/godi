// Package server implements a godi web server, hosting a thick client
package server

import (
	"net/http"

	"github.com/elazarl/go-bindata-assetfs"
)

// Return a new server instance which is initialized with the given data
type Server struct {
	srv http.Server
}

// Initialize the server to serve on the given address
func New(address string) *Server {
	mux := http.NewServeMux()

	s := Server{
		srv: http.Server{
			Addr:    address,
			Handler: mux,
		},
	}

	mux.Handle("/", http.FileServer(
		&assetfs.AssetFS{
			Asset:    Asset,
			AssetDir: AssetDir,
			Prefix:   "",
		},
	))

	return &s
}

// Start the server to listen from the address it was initialized with.
// Will not return unless abort is requested
// Error may occur if network resources couldn't be used
func (s *Server) Run() error {
	return s.srv.ListenAndServe()
}
