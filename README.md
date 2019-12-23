# express
An Express inspired web router for Go using Fasthttp engine

This project is in alpha stage, do not use in production

# Getting started
This app starts a server and listens on port 3000 for connections. The app responds with “Hello World!” for requests to the root URL (/) or route. For every other path, it will respond with a 404 Not Found.

```go
package main

import (
	"github.com/fenny/express"
)

func main() {
	app := express.Router()

	app.Get("/", func(c *express.Context) {
		c.Send("Hello World!")
	})

	app.Listen(8080)
}
```

### c.Hostname()
Contains the hostname derived from the Host HTTP header.
```go
// Host: "example.com:8080"
fmt.Println(c.Hostname())
// => `example.com`
```

### c.Ip()
Contains the remote IP address of the request.
```go
fmt.Println(c.Ip())
// => '127.0.0.1'
```

### c.Method()
Contains a string corresponding to the HTTP method of the request: GET, POST, PUT, and so on.
```go
// GET http://localhost:8080
fmt.Println(c.Method())
// => 'GET'
```


### c.Method()
Contains a string corresponding to the HTTP method of the request: GET, POST, PUT, and so on.
```go
// GET http://localhost:8080
fmt.Println(c.Method())
// => 'GET'
```
