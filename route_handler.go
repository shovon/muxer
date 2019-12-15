package main

import (
	"net/http"
)

type routeHandler map[string]http.Handler

func (r *routeHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	handler, ok := (*r)[req.Method]
	if !ok {
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
		return
	}
	handler.ServeHTTP(w, req)
}
