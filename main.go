package main

import (
	"fmt"
	"net/http"
)

func main() {
	// router := NewRouter()
	// router.Add("/foo/bar", 24)
	// router.Add("/foo/baz", 42)
	// router.Add("/foo/bar/:number", 3)
	// router.Add("/foo/bar/:number/something", "haha")
	// router.Add("/foo/bar", 10)

	// fmt.Println(router.Get("/foo/bar"))
	// fmt.Println(router.Get("/foo/baz"))
	// fmt.Println(router.Get("/foo/noexist"))
	// fmt.Println(router.Get("/foo/bar/baz"))
	// fmt.Println(router.Get("/foo/bar/baz/something"))

	muxer := NewMuxer()
	// muxer.AddGetHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, World!"))
	// }))
	// muxer.AddPostHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	w.Write([]byte("Hello, Post!"))
	// }))
	// muxer.AddGetHandler("/foo/:bar", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	params := Params(r)
	// 	w.Write([]byte(params["bar"]))
	// }))
	// muxer.AddGetHandler("/foo/:bar/baz/:widgets", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	params := Params(r)
	// 	w.Write([]byte(params["bar"]))
	// 	w.Write([]byte("\n"))
	// 	w.Write([]byte(params["widgets"]))
	// 	w.Write([]byte("\n"))
	// }))
	subMuxer := NewMuxer()
	subMuxer.AddGetHandler("/foo/:something", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params := Params(r)
		w.Write([]byte(params["something"]))
		w.Write([]byte("\n"))
	}))
	muxer.AddGetHandler("/bar/*", subMuxer)
	muxer.AddGetHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("LOL\n"))
	}))

	fmt.Println("Hopefully server is liistening on 0.0.0.0:8080")
	panic(http.ListenAndServe("0.0.0.0:8080", muxer))
}
