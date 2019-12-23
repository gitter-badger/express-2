package express

import (
	"io"
	"mime"
	"path/filepath"
	"sync"

	"github.com/valyala/fasthttp"
)

var ctxPool = sync.Pool{
	New: func() interface{} {
		return new(Context)
	},
}

func acquireCtx(fctx *fasthttp.RequestCtx) *Context {
	ctx := ctxPool.Get().(*Context)
	ctx.Fasthttp = fctx
	return ctx
}

func releaseCtx(ctx *Context) {
	ctx.next = false
	ctx.params = nil
	ctx.Fasthttp = nil
	ctxPool.Put(ctx)
}

// Next :
func (ctx *Context) Next() {
	ctx.next = true
	ctx.params = nil
}

// ParseRange https://expressjs.com/en/4x/api.html#req.ip
func (ctx *Context) ParseRange() string {
	//fasthttp.ParseByteRange(string("Range: bytes=0-1023"), 5000)
	return ctx.Fasthttp.RemoteIP().String()
}

// Ip https://expressjs.com/en/4x/api.html#req.ip
func (ctx *Context) Ip() string {
	return ctx.Fasthttp.RemoteIP().String()
}

// Url https://expressjs.com/en/4x/api.html#req.originalUrl
func (ctx *Context) Url() string {
	return string(ctx.Fasthttp.RequestURI())
}

// Query https://expressjs.com/en/4x/api.html#req.query
func (ctx *Context) Query(key string) string {
	return string(ctx.Fasthttp.QueryArgs().Peek(key))
}

// Params https://expressjs.com/en/4x/api.html#req.params
func (ctx *Context) Params(key string) string {
	if ctx.params == nil {
		return ""
	}
	for i := range ctx.params {
		if ctx.params[i] == key {
			return ctx.values[i]
		}
	}
	return ""
}

// Method https://expressjs.com/en/4x/api.html#req.method
func (ctx *Context) Method() string {
	return string(ctx.Fasthttp.Method())
}

// Path https://expressjs.com/en/4x/api.html#req.path
func (ctx *Context) Path() string {
	return string(ctx.Fasthttp.Path())
}

// Secure https://expressjs.com/en/4x/api.html#req.secure
func (ctx *Context) Secure() bool {
	return ctx.Fasthttp.IsTLS()
}

// Xhr https://expressjs.com/en/4x/api.html#req.xhr
func (ctx *Context) Xhr() bool {
	if ctx.Get("X-Requested-With") == "XMLHttpRequest" {
		return true
	}
	return false
}

// Protocol https://expressjs.com/en/4x/api.html#req.protocol
func (ctx *Context) Protocol() string {
	if ctx.Fasthttp.IsTLS() {
		return "https"
	}
	return "http"
}

// Is https://expressjs.com/en/4x/api.html#req.is
func (ctx *Context) Is(ext string) bool {
	if ext[0] != '.' {
		ext = "." + ext
	}
	extensions, _ := mime.ExtensionsByType(ctx.Get("Content-Type"))
	if len(extensions) > 0 {
		for _, item := range extensions {
			if item == ext {
				return true
			}
		}
	}
	return false
}

// Type https://expressjs.com/en/4x/api.html#res.type
func (ctx *Context) Type(ext string) {
	if ext[0] != '.' {
		ext = "." + ext
	}
	m := mime.TypeByExtension(ext)
	ctx.Fasthttp.Response.Header.Set("Content-Type", m)
}

// Attachment https://expressjs.com/en/4x/api.html#res.attachment
func (ctx *Context) Attachment(args ...interface{}) {
	if len(args) == 1 {
		filename, ok := args[0].(string)
		if !ok {
			panic("Invalid string")
		}
		ctx.Type(filepath.Ext(filename))
		ctx.Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		ctx.Set("Content-Disposition", "attachment")
	}
}

// Set https://expressjs.com/en/4x/api.html#res.set
func (ctx *Context) Set(key string, value string) {
	ctx.Fasthttp.Response.Header.Set(key, value)
}

// Get https://expressjs.com/en/4x/api.html#res.get
func (ctx *Context) Get(key string) string {
	return string(ctx.Fasthttp.Response.Header.Peek(key))
}

// Redirect https://expressjs.com/en/4x/api.html#res.redirect
func (ctx *Context) Redirect(args ...interface{}) *Context {
	if len(args) == 1 {
		str, ok := args[0].(string)
		if ok {
			ctx.Fasthttp.Redirect(str, 302)
		} else {
			panic("Invalid string url")
		}
	} else if len(args) == 2 {
		str, sOk := args[1].(string)
		code, cOk := args[0].(int)
		if sOk && cOk {
			ctx.Fasthttp.Redirect(str, code)
		} else {
			panic("Invalid statuscode or string")
		}
	} else {
		panic("You cannot have more than 1 argument")
	}
	return ctx
}

// Status https://expressjs.com/en/4x/api.html#res.status
func (ctx *Context) Status(code int) *Context {
	ctx.Fasthttp.SetStatusCode(code)
	return ctx
}

// Send https://expressjs.com/en/4x/api.html#res.send
func (ctx *Context) Send(args ...interface{}) {
	if len(args) > 2 {
		panic("To many arguments")
	}
	if len(args) == 1 {
		str, ok := args[0].(string)
		if ok {
			ctx.Fasthttp.SetBodyString(str)
			return
		}
		byt, ok := args[0].([]byte)
		if ok {
			ctx.Fasthttp.SetBody(byt)
			return
		}
		panic("Invalid string or []byte")
	} else if len(args) == 2 {
		reader, rOk := args[0].(io.Reader)
		bodysize, bOk := args[0].(int)
		if rOk && bOk {
			ctx.Fasthttp.SetBodyStream(reader, bodysize)
		} else {
			panic("Invalid io.Reader or bodysize(int)")
		}
	} else {
		panic("You cannot have more than 2 arguments")
	}
}

// SendFile https://expressjs.com/en/4x/api.html#res.sendFile
func (ctx *Context) SendFile(path string) {
	ctx.Type(filepath.Ext(path))
	// Shit doesnt work correctly,
	ctx.Fasthttp.SendFile(path)
}

// Write https://nodejs.org/docs/v0.4.7/api/all.html#response.write
func (ctx *Context) Write(args ...interface{}) {
	if len(args) > 2 {
		panic("To many arguments")
	}
	if len(args) == 1 {
		str, ok := args[0].(string)
		if ok {
			ctx.Fasthttp.WriteString(str)
			return
		}
		byt, ok := args[0].([]byte)
		if ok {
			ctx.Fasthttp.Write(byt)
			return
		}
		panic("Invalid string or []byte")
	} else {
		panic("You cannot have more than 1 argument")
	}
}
