package web

import (
	"net/http"
	"fmt"
)

func handleError(err error, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(500)
	fmt.Fprintf(w, err.Error())
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.NotFound(w, r)
}