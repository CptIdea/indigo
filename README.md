<img src="indigo.svg" alt="This is just a logo" title="What are you looking for?"/>

Indigo is non-idiomatic, but focusing on simplicity and performance web-server

# Documentation

Documentation is available [here](https://floordiv.gitbook.io/indigo/). However, it isn't complete yet.

# Hello, world!

```golang
package main

import (
  "log"
  
  "github.com/indigo-web/indigo"
  "github.com/indigo-web/indigo/http"
  "github.com/indigo-web/indigo/router/inbuilt"
)

const addr = ":8080"

func MyHandler(request *http.Request) *http.Response {
  return request.Respond().String("Hello, world!")
}

func main() {
  r := inbuilt.New()
  r.Resource("/").
    Get(MyHandler).
    Post(MyHandler)

  app := indigo.NewApp(addr)
  if err := app.Serve(r); err != nil {
    log.Fatal(err)
  }
}
```

More examples in [examples/](https://github.com/indigo-web/indigo/tree/master/examples) folder.
