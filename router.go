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
	if path[0] != '/' && path[0] != '*' {
		panic("Path must begin with slash '/'")
	}

	regex, err := regexp.Compile(pathToRegex(path))
	if err != nil {
		panic(err)
	}

	params := findWildcards(path)

	r.routes = append(r.routes, &route{method, path, regex, params, handler})

}

// handler :
func (r *Express) handler(fctx *fasthttp.RequestCtx) {

	path := string(fctx.Path())
	method := string(fctx.Method())

	ctx := acquireCtx(fctx)

	for _, route := range r.routes {
		if route.method != method {
			continue
		}
		// static path, wohoo no regex
		if route.path == path {
			route.handler(ctx)

			if !ctx.next {
				break
			}
			ctx.next = false
			continue
		}
		// Oh boy, here we go *.*
		if !route.regex.MatchString(path) {
			continue
		}
		// Prolly best to check if we have params in the first place
		// if route.params > 0 {
		matches := route.regex.FindAllStringSubmatch(path, -1)

		if len(matches) > 0 && len(matches[0]) > 1 {
			ctx.params = route.params
			ctx.values = matches[0][1:len(matches[0])]
		}
		route.handler(ctx)

		if !ctx.next {
			break
		}
		ctx.next = false
	}

	releaseCtx(ctx)
}

// Listen :
func (r *Express) Listen(port int) error {
	return fasthttp.ListenAndServe(fmt.Sprintf(":%v", port), r.handler)
}
