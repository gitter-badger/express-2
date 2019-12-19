package express

import (
	"fmt"
	"github.com/valyala/fasthttp"
)

type (
	// Routes stores all handlers
	Routes []*Route
	// Request represents the contextual data and environment while processing an incoming HTTP request.
	Request struct {
		*fasthttp.RequestCtx
		index  int
		routes Routes
	}
	// Route represents a method and handler
	Route struct {
		method  string
		path 		string
		handler func(*Request)
	}
	// Express manages routes and dispatches HTTP requests to the handlers of the matching routes.
	Express struct {
		routes Routes
	}
)

// New a Express instance
func New() *Express {
	return &Express{}
}

// Use register a middleware only handle ANY method
func (ex *Express) Use(args ...interface{}) {
	ex.Register("*", args...)
}
// Connect register a middleware only handle CONNECT method
func (ex *Express) Connect(args ...interface{}) {
	ex.Register("CONNECT", args...)
}
// Delete register a middleware only handle DELETE method
func (ex *Express) Delete(args ...interface{}) {
	ex.Register("DELETE", args...)
}
// Get register a middleware only handle GET method
func (ex *Express) Get(args ...interface{}) {
	ex.Register("GET", args...)
}
// Head register a middleware only handle HEAD method
func (ex *Express) Head(args ...interface{}) {
	ex.Register("HEAD", args...)
}
// Options register a middleware only handle OPTIONS method
func (ex *Express) Options(args ...interface{}) {
	ex.Register("OPTIONS", args...)
}
// Patch register a middleware only handle PATCH method
func (ex *Express) Patch(args ...interface{}) {
	ex.Register("PATCH", args...)
}
// Post register a middleware only handle POST method
func (ex *Express) Post(args ...interface{}) {
	ex.Register("POST", args...)
}
// Put register a middleware only handle PUT method
func (ex *Express) Put(args ...interface{}) {
	ex.Register("PUT", args...)
}
// Trace register a middleware only handle TRACE method
func (ex *Express) Trace(args ...interface{}) {
	ex.Register("TRACE", args...)
}

// Register route to stack
func (ex *Express) Register(method string, args ...interface{}) {

	var path string
	var handler func(*Request)
	// func(path string, req handler)
	if len(args) == 2 {
		path, _ = args[0].(string)
		handler, _ = args[1].(func(*Request))
	}
	// func(req handler)
	if len(args) == 1 {
		path = "*"
		handler, _ = args[0].(func(*Request))
	}

	if path == "" || path[0] != '/' && path[0] != '*' {
		panic("The first params of Use func must be a string which start with '/' or '*'")
	}

	// Append Route to array
	ex.routes = append(ex.routes, &Route{method, path, handler})
	//ex.routes[path] = append(ex.routes[path], &Route{method, handler})
}

// Context sets the request and response of the contex
func (ex *Express) Context(fctx *fasthttp.RequestCtx) *Request {
	ctx := &Request{}

	ctx.RequestCtx = fctx
	ctx.index = 0
	ctx.routes = ex.routes

	return ctx
}

// Incomming handles the HTTP request.
func (ex *Express) Incomming(fctx *fasthttp.RequestCtx) {
	// Create custom context that contains index and routes
	ctx := ex.Context(fctx)
	ctx.FindHandler()

}
// FindHandler trying to match a route in array
func (ctx *Request) FindHandler() {
	method := string(ctx.Method())
	path := string(ctx.Path())
	for l := len(ctx.routes); ctx.index < l; ctx.index++ {
		x := *ctx.routes[ctx.index]
		if x.method == "*" && x.path == "*" {
			x.handler(ctx)
			return
		}
		if x.method == "*" && x.path == path {
			x.handler(ctx)
			return
		}
		if x.method == method && x.path == "*" {
			x.handler(ctx)
			return
		}
		if x.method == method && x.path == path {
			x.handler(ctx)
			return
		}
	}
	fmt.Println("No handler found")

}

// Next calls the rest of the handlers associated with the current route.
func (ctx *Request) Next() {
	ctx.index++
	ctx.FindHandler()
}

// Listen for incomming requests on specified port
func (ex *Express) Listen(port int) error {
	return fasthttp.ListenAndServe(fmt.Sprintf(":%v", port), ex.Incomming)
}
