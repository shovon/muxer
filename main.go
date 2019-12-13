package main

import (
	"fmt"
	"net/http"
)

func main() {
	muxer := NewMuxer()
	muxer.AddGetHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	}))
	muxer.AddPostHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, Post!"))
	}))
	muxer.AddGetHandler("/foo/:bar", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := Params(r)
		w.Write([]byte(params["bar"]))
	}))
	muxer.AddGetHandler("/foo/:bar/baz/:widgets", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := Params(r)
		w.Write([]byte(params["bar"]))
		w.Write([]byte("\n"))
		w.Write([]byte(params["widgets"]))
		w.Write([]byte("\n"))
	}))

	subMuxer := NewMuxer()
	subMuxer.AddGetHandler("/foo/:something", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := Params(r)
		w.Write([]byte(params["something"]))
		w.Write([]byte("\n"))
	}))
	subMuxer.AddPostHandler("/some-post", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Some post\n"))
	}))

	muxer.AddHandler("/bar/*", subMuxer)
	muxer.AddGetHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("LOL\n"))
	}))

	fmt.Println("Hopefully server is liistening on 0.0.0.0:8080")
	panic(http.ListenAndServe("0.0.0.0:8080", muxer))
}
