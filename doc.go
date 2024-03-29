/*
# muxer

Package `muxer` is a minimal implementation of the `http.Handler` interface, and its purpose is to route and dispatche HTTP requests to appropriate handlers. `muxer` supports routing by pattern and by HTTP methods attached to those patterns. If a pattern specifies parameters, a `Params` function can be used to extract the value by the given parameter.

## Usage

```go
import "github.com/shovon/muxer"

func hello(w http.ResponseWriter, r *http.Request) {
	value := muxer.Params(r)["name"]
	w.Write([]byte(value))
}

func main() {
	mux := muxer.NewMuxer()

	mux.AddGetHandlerFunc("/hello/:name", hello)

	endpoint := "0.0.0.0:8080"
	fmt.Println("Server listening on " + endpoint)
	panic(http.ListenAndServe(endpoint, mux))
}
```
*/
package muxer
