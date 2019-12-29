package express

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/valyala/fasthttp"
)

// Router :
func New() *Express {
	return &Express{
		methods: []string{"GET", "PUT", "POST", "DELETE", "HEAD", "PATCH", "OPTIONS", "TRACE", "CONNECT"},
	}
}

// Get :
func (r *Express) Get(args ...interface{}) {
	r.register("GET", args...)
}

// Put :
func (r *Express) Put(args ...interface{}) {
	r.register("PUT", args...)
}

// Post :
func (r *Express) Post(args ...interface{}) {
	r.register("POST", args...)
}

// Delete :
func (r *Express) Delete(args ...interface{}) {
	r.register("DELETE", args...)
}

// Head :
func (r *Express) Head(args ...interface{}) {
	r.register("HEAD", args...)
}

// Patch :
func (r *Express) Patch(args ...interface{}) {
	r.register("PATCH", args...)
}

// Options :
func (r *Express) Options(args ...interface{}) {
	r.register("OPTIONS", args...)
}

// Trace :
func (r *Express) Trace(args ...interface{}) {
	r.register("TRACE", args...)
}

// Connect :
func (r *Express) Connect(args ...interface{}) {
	r.register("CONNECT", args...)
}

// All :
func (r *Express) All(args ...interface{}) {
	for _, method := range r.methods {
		r.register(method, args...)
	}
}

// Use :
func (r *Express) Use(args ...interface{}) {
	for _, method := range r.methods {
		r.register(method, args...)
	}
}

// Get :
func (r *Express) Static(mount string, dir string) {
	clean_dir := filepath.Clean(dir)
	files, err := dirWalk(dir)
	if err != nil {
		panic(err)
	}
	for _, file := range files {
		file_mount := mount + strings.Replace(file, clean_dir, "", 1)
		file_mount = strings.Replace(file_mount, "//", "/", 1)
		file_path := file
		if filepath.Base(file_path) == "index.html" {
			r.routes = append(r.routes, &route{"GET", mount, nil, nil, func(c *Context) {
				c.SendFile(file_path)
			}})
		}
		r.routes = append(r.routes, &route{"GET", file_mount, nil, nil, func(c *Context) {
			c.SendFile(file_path)
		}})
	}
}

// register :
func (r *Express) register(method string, args ...interface{}) {
	var path string
	var pathOk bool
	var handler handlerFunc
	var handlerOk bool
	if len(args) == 1 {
		handler, handlerOk = args[0].(handlerFunc)
		if !handlerOk {
			panic("Invalid handler")
		}
	} else if len(args) == 2 {

		path, pathOk = args[0].(string)
		handler, handlerOk = args[1].(func(*Context))

		if !pathOk || !handlerOk {
			panic("Invalid path or handler")
		}
	}
	// Panic if first char does not begins with / or *
	if path[0] != '/' && path[0] != '*' {
		panic("Path must begin with slash '/' or wildcard '*'")
	}

	// Compile regix from path
	regex, err := regexp.Compile(pathToRegex(path))
	if err != nil {
		panic("Invalid url pattern: " + path)
	}
	// Strip parameters from path
	params := stripParameters(path)
	if haveParameters(path) {
		r.routes = append(r.routes, &route{method, path, regex, params, handler})
	} else {
		r.routes = append(r.routes, &route{method, path, nil, nil, handler})
	}
}

// handler :
func (r *Express) handler(fctx *fasthttp.RequestCtx) {
	// get custom context from sync pool
	ctx := acquireCtx(fctx)
	// get path and method from main context
	path := ctx.Path()
	method := ctx.Method()
	// loop trough routes
	for _, route := range r.routes {
		// Skip route if method is not allowed
		if route.method != method {
			continue
		}
		// First check if we match a static path
		if route.path == path {
			// Execute handler with context
			route.handler(ctx)
			// if next is not set, leave loop and release ctx
			if !ctx.next {
				break
			}
			// set next to false for next iteration
			ctx.next = false
			// continue to go to the next route
			continue
		}
		// Skip route if regex does not match
		if route.regex == nil || !route.regex.MatchString(path) {
			continue
		}
		// If we have parameters, lets find the matches
		if route.params != nil && len(route.params) > 0 {
			matches := route.regex.FindAllStringSubmatch(path, -1)
			// If we have matches, add params and values to context
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
