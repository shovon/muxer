package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAddGetHandler(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddGetHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))

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

func TestAddGetHandlerFunc(t *testing.T) {
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

func TestAddPostHandler(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddPostHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))

	req, err := http.NewRequest("POST", "/foo", nil)
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

func TestAddPostHandlerFunc(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddPostHandlerFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	})

	req, err := http.NewRequest("POST", "/foo", nil)
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

func TestAddPutHandler(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddPutHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))

	req, err := http.NewRequest("PUT", "/foo", nil)
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

func TestAddPutHandlerFunc(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddPutHandlerFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	})

	req, err := http.NewRequest("PUT", "/foo", nil)
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

func TestAddDeleteHandler(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddDeleteHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))

	req, err := http.NewRequest("DELETE", "/foo", nil)
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

func TestAddDeleteHandlerFunc(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddDeleteHandlerFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	})

	req, err := http.NewRequest("DELETE", "/foo", nil)
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

func TestAddPatchHandler(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddPatchHandler("/foo", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	}))

	req, err := http.NewRequest("PATCH", "/foo", nil)
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

func TestAddPatchHandlerFunc(t *testing.T) {
	expected := "Haha"

	muxer := NewMuxer()
	muxer.AddPatchHandlerFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(expected))
	})

	req, err := http.NewRequest("PATCH", "/foo", nil)
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
