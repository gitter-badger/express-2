package express

import (
	"fmt"
	"regexp"

	"github.com/valyala/fasthttp"
)

// Router :
func Router() *Express {
	return &Express{
		methods: []string{"GET", "PUT", "POST", "DELETE", "HEAD", "PATCH", "OPTIONS", "TRACE", "CONNECT"},
	}
}

// Get :
func (r *Express) Get(path string, handler handlerFunc) {
	r.register("GET", path, handler)
}

// Put :
func (r *Express) Put(path string, handler handlerFunc) {
	r.register("PUT", path, handler)
}

// Post :
func (r *Express) Post(path string, handler handlerFunc) {
	r.register("POST", path, handler)
}

// Delete :
func (r *Express) Delete(path string, handler handlerFunc) {
	r.register("DELETE", path, handler)
}

// Head :
func (r *Express) Head(path string, handler handlerFunc) {
	r.register("HEAD", path, handler)
}

// Patch :
func (r *Express) Patch(path string, handler handlerFunc) {
	r.register("PATCH", path, handler)
}

// Options :
func (r *Express) Options(path string, handler handlerFunc) {
	r.register("OPTIONS", path, handler)
}

// Trace :
func (r *Express) Trace(path string, handler handlerFunc) {
	r.register("TRACE", path, handler)
}

// Connect :
func (r *Express) Connect(path string, handler handlerFunc) {
	r.register("CONNECT", path, handler)
}

// Use :
func (r *Express) Use(path string, handler handlerFunc) {
	for _, method := range r.methods {
		r.register(method, path, handler)
	}
}

// register :
func (r *Express) register(method string, path string, handler handlerFunc) {
	// Panic if first char does not begins with / or *
	if path[0] != '/' && path[0] != '*' {
		panic("Path must begin with slash '/' or wildcard '*'")
	}
	// Compile regix from path
	regex, err := regexp.Compile(pathToRegex(path))
	if err != nil {
		panic(err)
	}
	// Set parameters
	params := findParams(path)
	// Add to route
	r.routes = append(r.routes, &route{method, path, regex, params, handler})
}

// handler :
func (r *Express) handler(fctx *fasthttp.RequestCtx) {
	// get path and method from main context
	path := string(fctx.Path())
	method := string(fctx.Method())
	// get custom context from sync pool
	ctx := acquireCtx(fctx)
	// loop trough routes
	for _, route := range r.routes {
		// Skip route if method is not allowed
		if route.method != method {
			continue
		}
		// Check if path equals static
		if route.path == path {
			// Execute handler with context
			route.handler(ctx)
			// if next is not set, leave loop and release ctx
			if !ctx.next {
				break
			}
			// set next to false for next iteration
			ctx.next = false
			// continue to skip the regex calls
			continue
		}
		// Skip route if regex does not match
		if !route.regex.MatchString(path) {
			continue
		}
		// If we have parameters, lets find the matches
		if len(route.params) > 0 {
			matches := route.regex.FindAllStringSubmatch(path, -1)
			// If matches, add params and values to context
			if len(matches) > 0 && len(matches[0]) > 1 {
				ctx.params = &route.params
				ctx.values = matches[0][1:len(matches[0])]
			}
		}
		// Execute handler with context
		route.handler(ctx)
		// if next is not set, leave loop and release ctx
		if !ctx.next {
			break
		}
		// set next to false for next iteration
		ctx.next = false
	}
	// release context back into sync pool
	releaseCtx(ctx)
}

// Listen :
func (r *Express) Listen(port int) error {
	return fasthttp.ListenAndServe(fmt.Sprintf(":%v", port), r.handler)
}
