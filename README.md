# Express for Go [![GoDoc](https://godoc.org/github.com/fenny/express?status.svg)](http://godoc.org/github.com/fenny/express) [![fuzzit](https://app.fuzzit.dev/badge?org_id=express&branch=master)](https://fuzzit.dev) [![Go Report](https://goreportcard.com/badge/github.com/fenny/express)](https://goreportcard.com/report/github.com/fenny/express)
Express is a router framework build on top of the fastest HTTP package [FastHTTP](https://github.com/valyala/fasthttp) for Go.  
This library is inspired by one of the most populair and well known web frameworks [ExpressJS](https://github.com/expressjs/express).

**This project is in alpha stage, do not use in production**

Express Go is **9** times faster than Express JS

<p align="left">
  <img height="100" src="https://i.imgur.com/Jh2BZ42.png">
</p>


[Install](#install)
[Getting started](#getting-started)  
[Benchmarks (Alpha)](#benchmarks-alpha)
[Route paths](#route-paths)  
[Route parameters](#route-parameters)  
[Route Handlers](#route-handlers)  
[API reference](#api-reference)  
[Credits](#credits--thanks)

## Install
```
go get -u github.com/fenny/express
```

## Getting started
This app starts a server and listens on port 8080 for connections. The app sets a 204 status code for all favicon.ico requests. /api/:user uses a middleware to set a dns-prefetch header before it responds with “Hello User!”. For every other path, it will respond with a 404 Not Found.

```go
package main

import "github.com/fenny/express"

func main() {
  app := express.Router()
  app.Get("/favicon.ico", func(c *express.Context) {
    c.Status(204)
  })
  app.Get("/api/*", func(c *express.Context) {
    // c.Params("*")
    c.Set("x-dns-prefetch-control", "on")
    c.Next()
  })
  app.Get("/api/:user/:imhere?", func(c *express.Context) {
    c.Write("Hello " + c.Params("user"))
    c.Write("Are you there? " + c.Params("imhere"))
  })
  app.All("*", func(c *express.Context) {
    c.Status(404).Send("404 Not Found.")
  })
  app.Listen(8080)
}
```


## Benchmarks (Alpha)
**Concurrency**  
In 30 ms processing time, the test result for 100, 1000, 5000 clients is:
![Benchmark (Round 3)](https://i.imgur.com/TtlMK42.png)

**Processing time**  
0 ms, 10 ms, 100 ms, 500 ms processing time in handlers.
![Benchmark (Round 3)](https://i.imgur.com/Auf621y.png)
![Benchmark (Round 3)](https://i.imgur.com/6pDQ2C1.png)

We still need to optimize this library, somewhere in the code it's still CPU bound.

## Route paths
Here are some examples of route paths based on strings.  
This route path will match requests to the root route, **/**.
```go
app.Get("/", func(c *express.Context) {
  c.Send("/")
})
```
This route path will match requests to **/about**.
```go
app.Get("/about", func(c *express.Context) {
  c.Send("/about")
})
```
This route path will match requests to **/random.text**.
```go
app.Get("/random.text", func(c *express.Context) {
  c.Send("/random.text")
})
```
Here are some examples of route paths based on string patterns.  
This route path will match **acd** and **abcd**.
```go
app.Get("/ab?cd", func(c *express.Context) {
  c.Send("/ab?cd")
})
```
 This route path will match **abcd**, **abbcd**, **abbbcd**, and so on.
```go
app.Get("/ab+cd", func(c *express.Context) {
  c.Send("/ab+cd")
})
```
This route path will match **/abe** and **/abcde**.
```go
app.Get("/ab(cd)?e", func(c *express.Context) {
  c.Send("/ab(cd)?e")
})
```

## Route parameters
Route parameters are named URL segments that are used to capture the values specified at their position in the URL. The captured values can be retreived with **Params(key string)**, with the name of the route parameter specified in the path as their respective keys.

```
Route path:           /users/:userId/books/:bookId
Request URL:          http://localhost:8080/users/34/books/8989
c.Params("userId") // STRING => 34
c.Params("bookId") // STRING => 8989
```

To define routes with route parameters, simply specify the route parameters in the path of the route as shown below.
```go
app.Get("/users/:userId/books/:bookId", func(c *express.Context) {
  c.Write(c.Params("userId"))
  c.Write(c.Params("bookId"))
})
```
```
The name of route parameters must be made up of “word characters” ([A-Za-z0-9_]).
```

## Route Handlers
You can provide multiple callback functions that behave like middleware to handle a request. By using the **Next()** function you can continue to the next middleware. You can use this mechanism to impose pre-conditions on a route, don't call **Next()** it there’s no reason to proceed with the to the next middleware route.

## API reference
The **express.Context** struct represents the HTTP request and response and has properties for the request query string, parameters, body, HTTP headers, and so on. In this documentation, the struct is always referred to as '**c**'.

```go
app.Get("/user/:name", func(c *express.Context) {
  // Hostname() contains the hostname derived from the Host HTTP header.
  c.Hostname() // STRING => example.com

  // Ip() contains the remote IP address of the request.
  c.Ip() // STRING => 127.0.0.1

  // Method contains a string corresponding to the HTTP method of the request: GET, POST, PUT, and so on.
  c.Method()) // STRING => GET

  // Url() This contains only the URL that is present in the actual HTTP request.
  c.Url() // STRING => /search?q=something

  // Params(string key) method returns the named route “parameter”. For example, if you have the route /user/:name, then the “name” property is available as c.Params("name").
  c.Params("name") // STRING => john

  // Path() contains the path part of the request URL.
  c.Path() // STRING => /user/john

  // Protocol() contains the request protocol string: either http or (for TLS requests) https.
  c.Protocol() // STRING =>  https

  // Query(key string) method returns a query string parameter in the route. If there is no query string, it returns a empty string.
  c.Query("search") // STRING => books

  // Secure() is a Boolean property that is true if a TLS connection is established.
  c.Secure() // BOOL => true

  // Xhr() is A Boolean property that is true if the request’s X-Requested-With header field is “XMLHttpRequest”.
  c.Xhr() // BOOL => true

  // Get(key string) returns the specified HTTP request header field (case-insensitive match). The Referrer and Referer fields are interchangeable.
  c.Get("Content-Type") // STRING => text/plain

  // Is(extension string) returns the matching content type if the incoming request’s “Content-Type” HTTP header field matches the MIME type specified by the type parameter.
  c.Is("json") // STRING => application/json

  // Attachment(filename string) sets the HTTP response Content-Disposition header field to “attachment”. If a filename is given, then it sets the Content-Type based on the extension name via c.Type(), and sets the Content-Disposition “filename=” parameter.
  c.Attachment("json")

  // Redirect(status int, path string) Redirects to the URL derived from the specified path, with specified status, a positive integer that corresponds to an HTTP status code . If not specified, status defaults to “302 “Found”.
  c.Redirect("/foo/bar")
  c.Redirect("http://google.com")
  c.Redirect(301, "http://google.com")

  // Send(body) Sends the HTTP response, the body parameter can be a []byte, string or a Reader
  c.Send("Hello World!")
  c.Send([]byte("Hello World!"))
  c.Send(reader io.Reader, bodysize int)

  // SendFile(path string) transfers the file at the given path. Sets the Content-Type response HTTP header field based on the filename’s extension.
  c.SendFile("tutorial.pdf")

  // Status(code int) sets the response HTTP status code to statusCode and send its string representation as the response body. It is a chainable method.
  c.Status(500)
  c.Status(200).Send("Woohoo")

  // Type(extension string) Sets the Content-Type HTTP header to the MIME type as determined by the strings file extention
  c.Type("json") // => c.Set("Content-Type", "application/json")
  c.Type("html") // => c.Set("Content-Type", "text/html")
  c.Type("html") // => c.Set("png", "image/png")

  // FastHttp contains all functions provided by the FastHTTP engine. In case express does not cover your needs, you can always fallback using the FastHTTP struct
  c.FastHttp
})
```

## Credits & Thanks
[Valyala (FastHTTP)](https://github.com/valyala)  
[Julien Schmidt (httpRouter)](https://github.com/julienschmidt/httprouter)  
[Savsgio (Atreugo)](https://github.com/savsgio/atreugo)  
[Vincent Liuxiang (Lu)](https://github.com/savsgio/atreugo)  
[Steambap (Koa-Tree-Router)](https://github.com/steambap/koa-tree-router)  
