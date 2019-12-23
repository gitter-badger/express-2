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
// GET example.com
fmt.Println(c.Method())
// => 'GET'
```


### c.Url()
Request URL string. This contains only the URL that is present in the actual HTTP request. 
```go
// GET /search?q=something
fmt.Println(c.Method())
// => '/search?q=something'
```

### c.Params(key string)
This method returns the named route “parameters”. For example, if you have the route /user/:name, then the “name” property is available as c.Params("name").
```go
// GET /user/joe
fmt.Println(c.Params("name))
// => 'ryan'
```

### c.Path()
Contains the path part of the request URL.
```go
// example.com/users?sort=desc
fmt.Println(c.Params("name))
// => '/users'
```

### c.Protocol()
Contains the request protocol string: either http or (for TLS requests) https.
```go
// GET https://example.com
fmt.Println(c.Protocol())
// => 'https'
```

### c.Protocol()
This method returns the named route “parameters”. For example, if you have the route /user/:name, then the “name” property is available as c.Params("name").
```go
// GET https://example.com
fmt.Println(c.Protocol())
// => 'https'
```
