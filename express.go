package express

import (
	"fmt"
	"bytes"
	"github.com/valyala/fasthttp"
)

type (
	// Request represents the contextual data and environment while processing an incoming HTTP request.
	Request struct {
		*fasthttp.RequestCtx
		// the index of the currently executing handler in routes
		index  int
		routes *Routes
	}
	// Handler is the function for handling HTTP requests.
	Handler func(*Request)
	// Route represents a method and handler
	Route struct {
		method  []byte
		handler Handler
	}
	// Routes stores all handlers
	Routes map[string][]*Route
	// Express manages routes and dispatches HTTP requests to the handlers of the matching routes.
	Express struct {
		routes Routes
	}
)

// New a Express instance
func New() *Express {
	return &Express{
		routes: make(Routes),
	}
}

// Connect register a middleware only handle CONNECT method
func (ex *Express) Connect(path string, handler Handler) {
	ex.Register([]byte("CONNECT"), path, handler)
}
// Delete register a middleware only handle DELETE method
func (ex *Express) Delete(path string, handler Handler) {
	ex.Register([]byte("DELETE"), path, handler)
}
// Get register a middleware only handle GET method
func (ex *Express) Get(path string, handler Handler) {
	ex.Register([]byte("GET"), path, handler)
}
// Head register a middleware only handle HEAD method
func (ex *Express) Head(path string, handler Handler) {
	ex.Register([]byte("HEAD"), path, handler)
}
// Options register a middleware only handle OPTIONS method
func (ex *Express) Options(path string, handler Handler) {
	ex.Register([]byte("OPTIONS"), path, handler)
}
// Patch register a middleware only handle PATCH method
func (ex *Express) Patch(path string, handler Handler) {
	ex.Register([]byte("PATCH"), path, handler)
}
// Post register a middleware only handle POST method
func (ex *Express) Post(path string, handler Handler) {
	ex.Register([]byte("POST"), path, handler)
}
// Put register a middleware only handle PUT method
func (ex *Express) Put(path string, handler Handler) {
	ex.Register([]byte("PUT"), path, handler)
}
// Trace register a middleware only handle TRACE method
func (ex *Express) Trace(path string, handler Handler) {
	ex.Register([]byte("TRACE"), path, handler)
}

// Register route to stack
func (ex *Express) Register(method []byte, path string, handler Handler) {
	if path == "" || path[0] != '/' {
		panic("The first params of Use func must be a string which start with '/'")
	}
	// Append Route to path array
	ex.routes[path] = append(ex.routes[path], &Route{method, handler})
}

// Context sets the request and response of the contex
func (ex *Express) Context(fctx *fasthttp.RequestCtx) *Request {
	ctx := &Request{}

	ctx.RequestCtx = fctx
	ctx.index = 0
	ctx.routes = &ex.routes

	return ctx
}

// Incomming handles the HTTP request.
func (ex *Express) Incomming(fctx *fasthttp.RequestCtx) {
	// Create string url from bytes
	url := string(fctx.Path())
	// Check if path exist
	if ex.routes[url] == nil {
		return // path not found
	}
	// Create custom context that contains index and pointer to routes
	ctx := ex.Context(fctx)

	// Loop trough all routes until we have a method match, this is our starting point
	for l := len((*ctx.routes)[url]); ctx.index < l; ctx.index++ {
		if bytes.Compare((*ctx.routes)[url][ctx.index].method, ctx.Method()) == 0 {
			(*ctx.routes)[url][ctx.index].handler(ctx)
			return
		}
	}

	// no match found with specific method
	// do something

}

// Next calls the rest of the handlers associated with the current route.
func (ctx *Request) Next() {
	// Create string url from bytes
	url := string(ctx.Path())
	ctx.index++
	for l := len((*ctx.routes)[url]); ctx.index < l; ctx.index++ {
		if bytes.Compare((*ctx.routes)[url][ctx.index].method, ctx.Method()) == 0 {
			(*ctx.routes)[url][ctx.index].handler(ctx)
			return
		}
	}

}

// Listen for incomming requests on specified port
func (ex *Express) Listen(port int) error {
	return fasthttp.ListenAndServe(fmt.Sprintf(":%v", port), ex.Incomming)
}
