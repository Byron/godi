package rest

import (
	"net/http"
)

type dirHandler struct {
}

func NewDirHandler() http.Handler {
	return &dirHandler{}
}

func (r *dirHandler) ServeHTTP(w http.ResponseWriter, rq *http.Request) {
	if rq.Method != "GET" {
		http.Error(w, "Only GET is allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set(ContentKey, JsonContent)

	println(rq.RequestURI)
}
