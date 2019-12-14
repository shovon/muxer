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

type middleware func(http.Handler) http.Handler

type middlewareNode struct {
	value middleware
	next  *middlewareNode
}

func (n *middlewareNode) insert(value middleware) *middlewareNode {
	return &middlewareNode{value, n}
}

type wrapperServer struct {
	muxer *Muxer
}

func (ws wrapperServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// We want to strip the prefix. But under what logic?
	//
	// No matter how nested this instance is, we will typically get the full URL
	// path.
	//
	// And so, grab this variable to strip out the prefix.
	offset, ok := r.Context().Value(pathOffsetContextKey).(int)
	if !ok {
		offset = 0
	}

	// Get the newly sliced route
	pathComponents := strings.Split(r.URL.Path[1:], "/")[offset:]
	partialPath := "/" + strings.Join(pathComponents, "/")

	result := ws.muxer.routes.getPartial(partialPath)

	if !result.Retrieved {
		ws.muxer.notFoundHandler.ServeHTTP(w, r)
		return
	}

	switch handler := result.Value.(type) {
	case *RouteHandler:
		if handler == nil {
			ws.muxer.notFoundHandler.ServeHTTP(w, r)
		} else {
			h, ok := (*handler)[r.Method]
			if !ok {
				ws.muxer.ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		}
	case http.Handler:
		if handler == nil {
			ws.muxer.notFoundHandler.ServeHTTP(w, r)
		} else {
			handler.ServeHTTP(w, r)
		}
	}
}

// TODO: determine if the field `routes` should not be a pointer.

// Muxer the main muxer library.
type Muxer struct {
	routes          *routes
	notFoundHandler http.Handler
	middlewares     *middlewareNode
	chain           http.Handler
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(404)
	w.Write([]byte("Not found"))
}

// NewMuxer creates a new muxer instance.
func NewMuxer() *Muxer {
	routes := newRouter()
	return &Muxer{
		&routes,
		http.HandlerFunc(notFound),
		nil,
		nil,
	}
}

func (m *Muxer) wrapHandler(path string, h http.Handler) http.Handler {
	// This function will also be filtering out all requests that are not
	// wildcards.

	// This will take the parent's route and relay that over to the child route,
	// so that the child route can make adjustments.
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get the offset.
		pathOffset, ok := r.Context().Value(pathOffsetContextKey).(int)
		if !ok {
			pathOffset = 0
		}

		// The path without the wildcard.
		pathNoWildcard := extractRelevantPath(path)

		// The first slash is a distraction.
		components := strings.Split(pathNoWildcard[1:], "/")

		// Check to see if the path is a wildcard. If not, and the does not match,
		// just respond with a 404.
		if !pathHasWildcard(path) {
			// Cut out the irrelevant stuff from the HTTP request.
			relevantRequestPathComponents :=
				strings.Split(r.URL.Path[1:], "/")[pathOffset:]

			// Check to see if the relevant components match the route path
			// components.
			if len(components) != len(relevantRequestPathComponents) {
				w.WriteHeader(404)
				w.Write([]byte("Not found"))
				return
			}
		}

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

// Params grabs the parameters from the URL.
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

func pathHasWildcard(path string) bool {
	components := strings.Split(path[1:], "/")
	return components[len(components)-1] == "*"
}

func extractRelevantPath(path string) string {
	if pathHasWildcard(path) {
		// Remove the top two characters. First the *, then the /.
		//
		// For instance, /foo/bar/* will now becmoe /foo/bar
		return path[:len(path)-2]
	}
	return path
}

func (m *Muxer) addHandlerMethod(path string, method string, h http.Handler) {
	// TODO: don't only check for `RouteHandler`s. Sometimes, we want the
	// catch-all types.

	nonWildcardPath := extractRelevantPath(path)

	handler, ok := m.routes.get(nonWildcardPath).(*RouteHandler)
	h = m.wrapHandler(path, h)
	if handler == nil || !ok {
		m.routes.add(nonWildcardPath, &RouteHandler{method: h})
	} else {
		(*handler)[method] = h
	}
}

func (m *Muxer) addCatchAllHandler(path string, h http.Handler) {
	nonWildcardPath := extractRelevantPath(path)

	m.routes.add(nonWildcardPath, m.wrapHandler(path, h))
}

// Use adds a middleware.
func (m *Muxer) Use(f middleware) {
	// N.B. computational complexity of `Use` is linear.  Not quite sure on how
	// to optimize this. Nevertheless, for small numbers of inputs, the
	// performance impact should be minimal, and, this function almost always run
	// at initialization.

	m.middlewares = m.middlewares.insert(f)

	m.chain = wrapperServer{m}

	// This initial middleware is the last middleware, and it is what will be
	// invoking the remainder of the middlewares.
	n := m.middlewares

	for n != nil {
		m.chain = n.value(m.chain)
		n = m.middlewares.next
	}
}

// AddGetHandler adds an http.Handler associated with a GET request to the
// specified route.
func (m *Muxer) AddGetHandler(route string, h http.Handler) {
	m.addHandlerMethod(route, "GET", h)
}

// AddGetHandlerFunc adds a GET http.HandlerFunc associated with a GET request
// to the specified route.
func (m *Muxer) AddGetHandlerFunc(path string, h http.HandlerFunc) {
	m.addHandlerMethod(path, "GET", http.HandlerFunc(h))
}

// AddPostHandler adds an http.Handler associated with a POST request to the
// specified route.
func (m *Muxer) AddPostHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "POST", h)
}

// AddPutHandler adds an http.Handler associated with a PUT request to the
// specified route.
func (m *Muxer) AddPutHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "PUT", h)
}

// AddDeleteHandler adds an http.Handler associated with a DELETE request to the
// specified route.
func (m *Muxer) AddDeleteHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "DELETE", h)
}

// AddPatchHandler adds
func (m *Muxer) AddPatchHandler(path string, h http.Handler) {
	m.addHandlerMethod(path, "PATCH", h)
}

// AddHandler adds a handler associated with any HTTP method request to the
// specified route.
func (m *Muxer) AddHandler(path string, h http.Handler) {
	m.addCatchAllHandler(path, h)
}

// SetNotFoundHandler sets the not found handler.
func (m *Muxer) SetNotFoundHandler(h http.Handler) {
	m.notFoundHandler = h
}

// This is the entry-point for the entire muxer's HTTP request.
func (m *Muxer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if m.chain == nil {
		m.chain = wrapperServer{m}
	}
	m.chain.ServeHTTP(w, req)
}
