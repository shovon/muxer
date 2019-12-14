package main

import "testing"

import "net/http"

import "net/http/httptest"

func TestAdd(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddGetHandlerFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	})

	req, err := http.NewRequest("GET", "/foo", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	muxer.ServeHTTP(rr, req)

	expected = "Haha"
	if body := rr.Body.String(); body != expected {
		t.Errorf("handler returned unexpected: want %v, but got %v", expected, body)
	}
}

func TestMethods(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddGetHandlerFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(expected))
	})

	muxer.AddPostHandlerFunc(
		"/foo",
		func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte(expected))
		},
	)
}
