package main

import (
	"context"
	"net/http"
	"strings"
)

type key string

const (
	pathContextKey               key = "muxer_pathContextKey"
	pathOffsetContextKey         key = "muxer_pathOffsetContextKey"
	previousPathOffsetContextKey key = "muxer_previousPathOffsetContextKey"
)

type Muxer struct {
	routes *Routes
}

func NewMuxer() *Muxer {
	routes := NewRouter()
	return &Muxer{&routes}
}

// This is the wrapHandler method
func wrapHandler(path string, h http.Handler) http.Handler {
	// path is /foo/bar/baz
	// offset would then be 3
	//
	// path is /something/another
	// offset would then be 5
	//
	// path is /even/more
	// offset would be 7

	// This will take the parent's route and relay that over to the child route,
	// so that the child route can make adjustments.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathOffset, ok := r.Context().Value(pathOffsetContextKey).(int)
		if !ok {
			pathOffset = 0
		}

		// The first slash is a distractino.
		components := strings.Split(path[1:], "/")
		newOffset := pathOffset + len(components)

		// Store the route's path (not the request path)
		// ctx = context.WithValue(r.Context(), pathContextKey, path)
		ctx := context.WithValue(r.Context(), pathOffsetContextKey, newOffset)
		ctx = context.WithValue(
			ctx,
			previousPathOffsetContextKey,
			pathOffset,
		)
		ctx = context.WithValue(ctx, pathContextKey, path)

		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func Params(r *http.Request) map[string]string {
	previousPathOffset, ok :=
		r.Context().Value(previousPathOffsetContextKey).(int)
	if !ok {
		previousPathOffset = 0
	}
	routePath, ok := r.Context().Value(pathContextKey).(string)
	if !ok {
		routePath = "/"
	}
	requestPathComponents :=
		strings.Split(r.URL.Path[1:], "/")[previousPathOffset:]
	routePathComponents := strings.Split(routePath[1:], "/")

	result := make(map[string]string)
	for i := 0; i < len(routePathComponents); i++ {
		if routePathComponents[i][0] == ':' {
			result[routePathComponents[i][1:]] = requestPathComponents[i]
		}
	}

	return result
}

func (m *Muxer) addHandlerMethod(path string, method string, h http.Handler) {
	handler, ok := m.routes.Get(path).(*RouteHandler)
	h = wrapHandler(path, h)
	if handler == nil || !ok {
		m.routes.Add(path, &RouteHandler{method: h})
	} else {
		(*handler)[method] = h
	}
}

func (m *Muxer) AddGetHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "GET", h)
}

func (m *Muxer) AddPostHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "POST", h)
}

func (m *Muxer) AddPutHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "PUT", h)
}

func (m *Muxer) AddDeleteHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "DELETE", h)
}

// This is the entry-point for the entire muxer's HTTP request.
func (m *Muxer) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	// We want to strip the prefix. But under what logic?
	//
	// No matter how nested this instance is, we will typically get the full URL
	// path.

	offset, ok := req.Context().Value(pathOffsetContextKey).(int)
	if !ok {
		offset = 0
	}

	// Get the newly sliced route
	pathComponents := strings.Split(req.URL.Path[1:], "/")[offset:]
	partialPath := "/" + strings.Join(pathComponents, "/")

	// Get the handler associated with the
	result := m.routes.GetPartial(partialPath)

	if !result.Retrieved {
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
		return
	}

	handler, ok := result.Value.(*RouteHandler)
	if handler == nil || !ok {
		w.WriteHeader(404)
		w.Write([]byte("Not found"))
	} else {
		handler.ServeHTTP(w, req)
	}

}
