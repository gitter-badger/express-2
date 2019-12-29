package express

import (
	"regexp"

	"github.com/valyala/fasthttp"
)

// Context :
type Context struct {
	next   bool
	params *[]string
	values []string
	//params   *map[string]string
	Fasthttp *fasthttp.RequestCtx
}

// handlerFunc :
type handlerFunc func(*Context)

type route struct {
	method  string
	path    string
	regex   *regexp.Regexp
	params  []string
	handler handlerFunc
}

// Express :
type Express struct {
	routes  []*route
	methods []string
}

// Cookie
type CookieOptions struct {
	domain   string
	expires  int
	httpOnly bool
	maxAge   int
	path     string
	secure   bool
	signed   bool
	sameSite string // "strict" or "laks"
}
