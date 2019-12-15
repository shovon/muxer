# muxer

Package `muxer` is a minimal implementation of the `http.Handler` interface, and its purpose is to route and dispatche HTTP requests to appropriate handlers. `muxer` supports routing by pattern and by HTTP methods attached to those patterns. If a pattern specifies parameters, a `Params` function can be used to extract the value by the given parameter.

## Usage

```go
import "github.com/shovon/muxer"

func main() {
  muxer := NewMuxer()
}
```