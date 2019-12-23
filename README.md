# Express running on FastHTTP
An Express inspired web router for Go using Fasthttp engine

This project is in alpha stage, do not use in production

### Getting started
This app starts a server and listens on port 8080 for connections. The app responds with “Hello World!” for requests to the root URL (/) or route. For every other path, it will respond with a 404 Not Found.

```go
package main

import (
	"github.com/fenny/express"
)

func main() {
	app := express.Router()

	app.Get("/api/:user", func(c *express.Context) {
		c.Send("Hello " + c.Params("user))
	})
  app.All("*", func(c *express.Context) {
    c.Status(404).Send("404 Not Found.")
  })
	app.Listen(8080)
}
```

### API Documentation
The ***express.Context** struct represents the HTTP request and response and has properties for the request query string, parameters, body, HTTP headers, and so on. In this documentation, the struct is always referred to as '**c**'.

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
})
```
