package muxer

import (
	"net/http"
	"strings"
)

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
