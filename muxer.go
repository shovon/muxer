package main

import (
	"context"
	"net/http"
	"strings"
)

type key string

// The following are context keys for helping track sub muxers, and help extract
// the URL parameters

const (
	pathContextKey               key = "muxer_pathContextKey"
	pathOffsetContextKey         key = "muxer_pathOffsetContextKey"
	previousPathOffsetContextKey key = "muxer_previousPathOffsetContextKey"
)

// A type alias for a middleware.
type middleware func(http.Handler) http.Handler

// We're storing the middlewares in a linked list.
type middlewareNode struct {
	value middleware
	next  *middlewareNode
}

// For insterting an item into the middlewares linked list.
func (n *middlewareNode) insert(value middleware) *middlewareNode {
	return &middlewareNode{value, n}
}

// This is the struct that will serve as the intermediary HTTP handler that
// will multiplex the routers that have been appened to the muxer.
type wrapperServer struct {
	muxer *Muxer
}

// Implementation of ServeHTTP that actually handles the multiplexing.
func (ws wrapperServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

// TODO: determine if the field `routes` should not be a pointer.

// Muxer the main muxer library.
type Muxer struct {
	routes          *routes
	notFoundHandler http.Handler
	middlewares     *middlewareNode
}

// Just the handlerfunc used for the not found response.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("Not found"))
}

// NewMuxer creates a new muxer instance.
func NewMuxer() *Muxer {
	routes := newRouter()
	return &Muxer{
		&routes,
		http.HandlerFunc(notFound),
		nil,
	}
}

// The purpose of this function is to handle path offsetting. Offsetting is done
// through the help of contexts that are embedded directly within the HTTP
// handlers. The offset is incremented at every handler.
//
// The one caveat is that if a non muxer handler is supplied at any level, then
// we would end up losing track. Maybe we might need to provide a workaround.
func (m *Muxer) wrapHandler(path string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		pathOffset, ok := r.Context().Value(pathOffsetContextKey).(int)
		if !ok {
			pathOffset = 0
		}

		pathNoWildcard := extractRelevantPath(path)

		// The first slash is a distraction.
		components := strings.Split(pathNoWildcard[1:], "/")

		if !pathHasWildcard(path) {
			// Cut out the irrelevant stuff from the HTTP request.
			relevantRequestPathComponents :=
				strings.Split(r.URL.Path[1:], "/")[pathOffset:]

			if len(components) != len(relevantRequestPathComponents) {
				w.WriteHeader(404)
				w.Write([]byte("Not found"))
				return
			}
		}

		newOffset := pathOffset + len(components)

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

// Determines if the given path ends with a wildcard character.
func pathHasWildcard(path string) bool {
	components := strings.Split(path[1:], "/")
	return components[len(components)-1] == "*"
}

// For now, this function will just drop the wildcard (*) character.
func extractRelevantPath(path string) string {
	if pathHasWildcard(path) {
		// Remove the top two characters. First the *, then the /.
		//
		// For instance, /foo/bar/* will now become /foo/bar
		return path[:len(path)-2]
	}
	return path
}

func (m *Muxer) addHandlerMethod(path string, method string, h http.Handler) {
	// TODO: don't only check for `RouteHandler`s. Sometimes, we want the
	// catch-all types.

	nonWildcardPath := extractRelevantPath(path)

	handler, ok := m.routes.get(nonWildcardPath).(*routeHandler)
	h = m.wrapHandler(path, h)
	if handler == nil || !ok {
		m.routes.add(nonWildcardPath, &routeHandler{method: h})
	} else {
		(*handler)[method] = h
	}
}

func (m *Muxer) addCatchAllHandler(path string, h http.Handler) {
	nonWildcardPath := extractRelevantPath(path)

	m.routes.add(nonWildcardPath, m.wrapHandler(path, h))
}

// AddGetHandler adds an http.Handler associated with a GET request to the
// specified route.
func (m *Muxer) AddGetHandler(route string, h http.Handler) {
	m.addHandlerMethod(route, "GET", h)
}

// AddGetHandlerFunc adds a GET http.HandlerFunc associated with a GET request
// to the specified route.
func (m *Muxer) AddGetHandlerFunc(route string, h http.HandlerFunc) {
	m.AddGetHandler(route, http.HandlerFunc(h))
}

// AddPostHandler adds an http.Handler associated with a POST request to the
// specified route.
func (m *Muxer) AddPostHandler(route string, h http.Handler) {
	m.addHandlerMethod(route, "POST", h)
}

// AddPostHandlerFunc adds a POST http.HandlerFunc associated with a POST
// request to the specified route.
func (m *Muxer) AddPostHandlerFunc(route string, h http.HandlerFunc) {
	m.AddPostHandler(route, http.HandlerFunc(h))
}

// AddPutHandler adds an http.Handler associated with a PUT request to the
// specified route.
func (m *Muxer) AddPutHandler(route string, h http.Handler) {
	m.addHandlerMethod(route, "PUT", h)
}

// AddPutHandlerFunc adds an http.HandlerFUnc associated with a PUT request to
// the specified route.
func (m *Muxer) AddPutHandlerFunc(route string, h http.HandlerFunc) {
	m.AddPutHandler(route, h)
}

// AddDeleteHandler adds an http.Handler associated with a DELETE request to the
// specified route.
func (m *Muxer) AddDeleteHandler(route string, h http.Handler) {
	m.addHandlerMethod(route, "DELETE", h)
}

// AddDeleteHandlerFunc adds an http.HandlerFunc associated with a DELETE
// request to the specified route.
func (m *Muxer) AddDeleteHandlerFunc(route string, h http.HandlerFunc) {
	m.AddDeleteHandler(route, http.HandlerFunc(h))
}

// AddPatchHandler adds an http.Handler associated with a PATCH request to the
// specified route.
func (m *Muxer) AddPatchHandler(route string, h http.Handler) {
	m.addHandlerMethod(route, "PATCH", h)
}

// AddPatchHandlerFunc adds an http.Handler associated with a PATCH request to
// the specified route.
func (m *Muxer) AddPatchHandlerFunc(route string, h http.HandlerFunc) {
	m.AddPatchHandler(route, http.HandlerFunc(h))
}

// AddCustomMethodHandler adds an http.Handler associated with a custom method
// to the specified route.
func (m *Muxer) AddCustomMethodHandler(method, route string, h http.Handler) {
	m.addHandlerMethod(route, method, h)
}

// AddCustomMethodHandlerFunc adds a http.HandlerFunc associated with a custom
// method to the specified route.
func (m *Muxer) AddCustomMethodHandlerFunc(
	method,
	route string,
	h http.HandlerFunc,
) {
	m.AddCustomMethodHandler(route, method, http.HandlerFunc(h))
}

// AddHandler adds a http.Handler associated with any HTTP method request to the
// specified route.
func (m *Muxer) AddHandler(path string, h http.Handler) {
	m.addCatchAllHandler(path, h)
}

// AddHandlerFunc adds a http.HandlerFunc associated with any HTTP method
// request to the specified route.
func (m *Muxer) AddHandlerFunc(path string, h http.HandlerFunc) {
	m.AddHandler(path, http.HandlerFunc(h))
}

// SetNotFoundHandler sets the not found handler.
func (m *Muxer) SetNotFoundHandler(h http.Handler) {
	m.notFoundHandler = h
}

// ServeHTTP is the entry-point for the entire muxer's HTTP request.
func (m *Muxer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// We want to strip the prefix. But under what logic?
	//
	// No matter how nested this instance is, we will typically get the full URL
	// path.
	//
	// And so, grab offset context variable, to strip out the prefix.
	offset, ok := req.Context().Value(pathOffsetContextKey).(int)
	if !ok {
		offset = 0
	}

	// Extract the relevant part of the path.
	pathComponents := strings.Split(req.URL.Path[1:], "/")[offset:]
	partialPath := "/" + strings.Join(pathComponents, "/")

	result := m.routes.getShortCircuited(partialPath)

	if !result.retrieved {
		m.notFoundHandler.ServeHTTP(w, req)
		return
	}

	switch handler := result.value.(type) {
	case *routeHandler:
		if handler == nil {
			m.notFoundHandler.ServeHTTP(w, req)
		} else {
			h, ok := (*handler)[req.Method]
			if !ok {
				m.notFoundHandler.ServeHTTP(w, req)
			} else {
				h.ServeHTTP(w, req)
			}
		}
	case http.Handler:
		if handler == nil {
			m.notFoundHandler.ServeHTTP(w, req)
		} else {
			handler.ServeHTTP(w, req)
		}
	default:
		m.notFoundHandler.ServeHTTP(w, req)
	}
}
